package api

import (
	"net/http"
	"time"

	"traceknot/internal/httputil"
	"traceknot/internal/store"
)

type skipOutcomeRequest struct {
	SessionID string `json:"session_id"`
}

func (picker *Picker) handleRecordSkip(writer http.ResponseWriter, request *http.Request) {
	var body skipOutcomeRequest
	if err := httputil.DecodeJSON(writer, request, &body); err != nil {
		return
	}
	if body.SessionID == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_session", "session_id is required")
		return
	}
	sessionID, err := picker.resolveSessionID(request.Context(), body.SessionID)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "outcome_failed", err.Error())
		return
	}
	if err := picker.store.RecordSkip(request.Context(), sessionID, time.Now().UnixMilli()); err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "outcome_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"ok": true})
}

func (picker *Picker) handleOutcome(writer http.ResponseWriter, request *http.Request) {
	sessionID, err := picker.resolveSessionID(request.Context(), request.URL.Query().Get("session_id"))
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "outcome_failed", err.Error())
		return
	}
	if sessionID == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_session", "session_id is required")
		return
	}

	outcome, err := picker.store.OutcomeForSession(request.Context(), sessionID)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "outcome_failed", err.Error())
		return
	}
	switch {
	case outcome == nil:
		httputil.WriteJSON(writer, http.StatusOK, map[string]any{"status": "pending"})
	case outcome.Status == store.ClaimStatusClaimed:
		httputil.WriteJSON(writer, http.StatusOK, map[string]any{
			"status":        "claimed",
			"work_item_key": outcome.WorkItemKey,
		})
	case outcome.Status == store.ClaimStatusSkipped:
		httputil.WriteJSON(writer, http.StatusOK, map[string]any{"status": "skipped"})
	default:
		httputil.WriteJSON(writer, http.StatusOK, map[string]any{"status": "pending"})
	}
}
