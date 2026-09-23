package api

import (
	"net/http"
	"sync"
	"time"

	"traceknot/internal/httputil"
)

const (
	maxContextEntries        = 128
	maxContextPromptBytes    = 128 * 1024
	maxContextSessionIDBytes = 1024
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
	now := time.Now()
	oldestID := ""
	var oldestExpiry time.Time
	for id, entry := range registry.items {
		if !entry.expiresAt.After(now) {
			delete(registry.items, id)
			continue
		}
		if oldestID == "" || entry.expiresAt.Before(oldestExpiry) {
			oldestID, oldestExpiry = id, entry.expiresAt
		}
	}
	if _, exists := registry.items[sessionID]; !exists && len(registry.items) >= maxContextEntries {
		delete(registry.items, oldestID)
	}
	registry.items[sessionID] = promptContextEntry{
		prompt:    prompt,
		expiresAt: now.Add(registry.ttl),
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
	if len(body.SessionID) > maxContextSessionIDBytes || len(body.Prompt) > maxContextPromptBytes {
		httputil.WriteError(writer, http.StatusRequestEntityTooLarge, "context_too_large", "prompt context is too large")
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
