package api

import (
	"net/http"
	"sync"
	"time"

	"traceknot/internal/httputil"
)

type contextRegistry struct {
	mu    sync.Mutex
	items map[string]promptContextEntry
	ttl   time.Duration
}

type promptContextEntry struct {
	prompt    string
	expiresAt time.Time
}

func newContextRegistry() *contextRegistry {
	return &contextRegistry{
		items: make(map[string]promptContextEntry),
		ttl:   10 * time.Minute,
	}
}

func (registry *contextRegistry) record(sessionID string, prompt string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.items[sessionID] = promptContextEntry{
		prompt:    prompt,
		expiresAt: time.Now().Add(registry.ttl),
	}
}

func (registry *contextRegistry) prompt(sessionID string) string {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	entry, found := registry.items[sessionID]
	if !found {
		return ""
	}
	if time.Now().After(entry.expiresAt) {
		delete(registry.items, sessionID)
		return ""
	}
	return entry.prompt
}

type storeContextRequest struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

func (picker *Picker) handleStoreContext(writer http.ResponseWriter, request *http.Request) {
	var body storeContextRequest
	if err := httputil.DecodeJSON(writer, request, &body); err != nil {
		return
	}
	if body.SessionID == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_session", "session_id is required")
		return
	}
	picker.contexts.record(body.SessionID, body.Prompt)
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"ok": true})
}

func (picker *Picker) handleGetContext(writer http.ResponseWriter, request *http.Request) {
	sessionID := request.URL.Query().Get("session_id")
	if sessionID == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_session", "session_id is required")
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{
		"prompt": picker.contexts.prompt(sessionID),
	})
}
