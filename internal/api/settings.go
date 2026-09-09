package api

import (
	"net/http"

	"traceknot/internal/config"
	"traceknot/internal/httputil"
)

type Settings struct{}

func NewSettings() *Settings {
	return &Settings{}
}

func (settings *Settings) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/settings", settings.handleGet)
	return mux
}

func (settings *Settings) handleGet(writer http.ResponseWriter, request *http.Request) {
	cfg := config.Load()
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"require_work_item": cfg.RequireWorkItem})
}
