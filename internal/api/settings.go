package api

import (
	"net/http"

	"traceknot/internal/httputil"
	"traceknot/internal/install/agentenv"
	"traceknot/internal/settings"
)

type Settings struct{}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/settings", s.handleGet)
	mux.HandleFunc("PATCH /api/v1/settings", s.handlePatch)
	mux.HandleFunc("POST /api/v1/settings/restart-sessions", s.handleRestartSessions)
	return mux
}

func (s *Settings) handleGet(writer http.ResponseWriter, request *http.Request) {
	state, err := settings.Current(request.Context())
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "settings_read_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, state)
}

func (s *Settings) handlePatch(writer http.ResponseWriter, request *http.Request) {
	var patch settings.Patch
	if err := httputil.DecodeJSON(writer, request, &patch); err != nil {
		return
	}
	if err := settings.Apply(request.Context(), patch); err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "settings_apply_failed", err.Error())
		return
	}
	state, err := settings.Current(request.Context())
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "settings_read_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, state)
}

func (s *Settings) handleRestartSessions(writer http.ResponseWriter, request *http.Request) {
	closed, err := agentenv.CloseInteractiveSessions(request.Context())
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "restart_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"closed": closed})
}
