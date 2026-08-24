package server

import (
	"context"
	"log/slog"
	"time"

	"traceknot/internal/ingest"
	"traceknot/internal/normalize/copilot"
	"traceknot/internal/normalize/shared"
	"traceknot/internal/rebuildstatus"
	"traceknot/internal/store"
)

type normalizerVersion struct {
	version int
	regroup ingest.RegroupFunc
}

var normalizerVersions = map[string]normalizerVersion{
	"claude":  {version: 1},
	"codex":   {version: 1},
	"copilot": {version: 3, regroup: copilot.RegroupNativeID},
}

func rebuildStaleNormalizers(
	storeHandle *store.Store,
	receiver *ingest.Receiver,
	normalizers map[string]shared.Normalizer,
	logger *slog.Logger,
	status *rebuildstatus.Status,
) {
	seen := map[string]bool{}
	for _, normalizer := range normalizers {
		provider := normalizer.Provider()
		if seen[provider] {
			continue
		}
		seen[provider] = true

		target, tracked := normalizerVersions[provider]
		if !tracked {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		current, err := storeHandle.NormalizerVersion(ctx, provider)
		if err != nil {
			logger.Warn("normalizer version lookup failed", "provider", provider, "error", err)
			cancel()
			continue
		}
		if current >= target.version {
			cancel()
			continue
		}

		status.Start(provider)
		logger.Info("rebuilding sessions for normalizer version bump",
			"provider", provider, "from", current, "to", target.version)
		if err := receiver.RebuildProvider(ctx, normalizer, target.regroup); err != nil {
			logger.Error("normalizer rebuild failed", "provider", provider, "error", err)
			status.Done(provider)
			cancel()
			continue
		}
		if err := storeHandle.SetNormalizerVersion(ctx, provider, target.version); err != nil {
			logger.Error("failed to persist normalizer version", "provider", provider, "error", err)
		} else {
			logger.Info("normalizer rebuild complete", "provider", provider, "version", target.version)
		}
		status.Done(provider)
		cancel()
	}
}
