package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"traceknot/internal/api"
	"traceknot/internal/httputil"
	"traceknot/internal/ingest"
	"traceknot/internal/ingest/claudeusage"
	"traceknot/internal/ingest/codexrollout"
	"traceknot/internal/normalize/claude"
	"traceknot/internal/normalize/codex"
	"traceknot/internal/normalize/copilot"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/pricing"
	"traceknot/internal/providers"
	"traceknot/internal/store"
	"traceknot/internal/tokenize"
	"traceknot/ui"
)

func serviceNormalizers(catalog *pricing.Catalog) map[string]shared.Normalizer {
	estimator := tokenize.NewEstimator(0)
	copilotNormalizer := copilot.NewNormalizer(catalog, estimator)
	return map[string]shared.Normalizer{
		"codex_cli_rs":   codex.NewNormalizer(catalog, estimator),
		"claude-code":    claude.NewNormalizer(catalog, estimator),
		"github-copilot": copilotNormalizer,
		"copilot-chat":   copilotNormalizer,
	}
}

func Run(dbPath string, listenAddr string, logger *slog.Logger) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return fmt.Errorf("create database directory: %w", err)
	}

	storeHandle, err := store.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer storeHandle.Close()

	reconcileClaimAttachments(storeHandle, logger)

	catalog, err := pricing.LoadEmbeddedCatalog()
	if err != nil {
		return fmt.Errorf("load pricing catalog: %w", err)
	}

	normalizers := serviceNormalizers(catalog)
	receiver := ingest.NewReceiver(storeHandle, normalizers, logger)

	watcherCancel := startRolloutWatcher(receiver, normalizers, logger)
	defer watcherCancel()

	usageWatcherCancel := startClaudeUsageWatcher(receiver, normalizers, logger)
	defer usageWatcherCancel()

	httpServer := &http.Server{
		Addr:              listenAddr,
		Handler:           requireSameOrigin(buildMux(storeHandle, receiver), listenAddr),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return serveUntilShutdown(httpServer, listenAddr, dbPath, logger)
}

func reconcileClaimAttachments(storeHandle *store.Store, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := storeHandle.ReconcileClaimAttachments(ctx); err != nil {
		logger.Warn("claim attachment reconciliation failed", "error", err)
	}
}

func startRolloutWatcher(
	receiver *ingest.Receiver,
	normalizers map[string]shared.Normalizer,
	logger *slog.Logger,
) context.CancelFunc {
	rolloutWatcher := codexrollout.New(receiver, normalizers["codex_cli_rs"], logger)
	ctx, cancel := context.WithCancel(context.Background())
	go rolloutWatcher.Run(ctx)
	return cancel
}

func startClaudeUsageWatcher(
	receiver *ingest.Receiver,
	normalizers map[string]shared.Normalizer,
	logger *slog.Logger,
) context.CancelFunc {
	usageWatcher := claudeusage.New(receiver, normalizers["claude-code"], logger)
	ctx, cancel := context.WithCancel(context.Background())
	go usageWatcher.Run(ctx)
	return cancel
}

func requireSameOrigin(next http.Handler, listenAddr string) http.Handler {
	allowed := allowedOrigins(listenAddr)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !isSafeMethod(request.Method) {
			if origin := request.Header.Get("Origin"); origin != "" && !allowed[origin] {
				httputil.WriteError(writer, http.StatusForbidden, "cross_origin_forbidden", "cross-origin requests are not allowed")
				return
			}
		}
		next.ServeHTTP(writer, request)
	})
}

func isSafeMethod(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

func allowedOrigins(listenAddr string) map[string]bool {
	_, port, err := net.SplitHostPort(listenAddr)
	if err != nil {
		port = ""
	}
	origins := map[string]bool{}
	for _, host := range []string{"127.0.0.1", "localhost", "[::1]"} {
		origin := "http://" + host
		if port != "" {
			origin += ":" + port
		}
		origins[origin] = true
	}
	return origins
}

func buildMux(storeHandle *store.Store, receiver *ingest.Receiver) http.Handler {
	registry := providers.NewRegistry(providers.NewGitHub(), providers.NewJira())
	picker := api.NewPicker(storeHandle, registry)
	dashboard := api.NewDashboard(storeHandle)
	sessionDetail := api.NewSessionDetail(storeHandle)
	sessionsList := api.NewSessions(storeHandle)
	workItems := api.NewWorkItems(storeHandle)

	mux := http.NewServeMux()
	mux.Handle("/v1/", receiver.Handler())
	mux.Handle("/api/v1/", picker.Handler())
	mux.Handle("/api/v1/dashboard", dashboard.Handler())
	mux.Handle("/api/v1/sessions", sessionsList.Handler())
	mux.Handle("/api/v1/work-items", workItems.Handler())
	mux.Handle("/api/v1/sessions/", sessionDetail.Handler())
	mux.Handle("/api/v1/nodes/", sessionDetail.Handler())
	mux.Handle("/", ui.Handler())
	mux.HandleFunc("GET /healthz", func(writer http.ResponseWriter, _ *http.Request) {
		httputil.WriteJSON(writer, http.StatusOK, map[string]any{"ok": true, "service": "traceknot"})
	})
	mux.HandleFunc("GET /readyz", func(writer http.ResponseWriter, _ *http.Request) {
		httputil.WriteJSON(writer, http.StatusOK, map[string]any{"ok": true, "ready": true})
	})
	return mux
}

func serveUntilShutdown(httpServer *http.Server, listenAddr string, dbPath string, logger *slog.Logger) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("traceknot listening", "addr", listenAddr, "db", dbPath)
		serveErr <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if err != http.ErrServerClosed {
			return err
		}
	case <-shutdown:
		logger.Info("shutting down")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpServer.Shutdown(ctx)
}
