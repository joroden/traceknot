package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"traceknot/internal/httputil"
	"traceknot/internal/providers"
	"traceknot/internal/ptr"
	"traceknot/internal/stableid"
	"traceknot/internal/store"
	"traceknot/internal/store/claim"
)

type Picker struct {
	store    *store.Store
	registry *providers.Registry
	contexts *contextRegistry
}

func NewPicker(storeHandle *store.Store, registry *providers.Registry) *Picker {
	return &Picker{
		store:    storeHandle,
		registry: registry,
		contexts: newContextRegistry(),
	}
}

func (picker *Picker) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/picker/recent", picker.handleRecent)
	mux.HandleFunc("POST /api/v1/picker/offer", picker.handleOfferPicker)
	mux.HandleFunc("GET /api/v1/picker/outcome", picker.handleOutcome)
	mux.HandleFunc("POST /api/v1/picker/outcome", picker.handleRecordSkip)
	mux.HandleFunc("GET /api/v1/picker/context", picker.handleGetContext)
	mux.HandleFunc("POST /api/v1/picker/context", picker.handleStoreContext)
	mux.HandleFunc("GET /api/v1/providers", picker.handleProviders)
	mux.HandleFunc("GET /api/v1/providers/{id}/search", picker.handleSearch)
	mux.HandleFunc("POST /api/v1/claims", picker.handleCreateClaim)
	mux.HandleFunc("POST /api/v1/claims/{id}/attach", picker.handleAttachClaim)
	return mux
}

func (picker *Picker) handleRecent(writer http.ResponseWriter, request *http.Request) {
	providerFilter := request.URL.Query().Get("provider")
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}

	items, err := picker.store.ListRecentWorkItems(request.Context(), providerFilter, limit)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "recent_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"items": items})
}

func (picker *Picker) handleProviders(writer http.ResponseWriter, request *http.Request) {
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{
		"providers": picker.registry.Probes(request.Context()),
	})
}

func (picker *Picker) handleSearch(writer http.ResponseWriter, request *http.Request) {
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	items, err := picker.registry.Search(request.Context(), request.PathValue("id"), request.URL.Query().Get("q"), limit)
	if err != nil {
		httputil.WriteError(writer, http.StatusBadGateway, "search_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"items": items})
}

type createClaimRequest struct {
	SessionID *string `json:"session_id"`
	Source    string  `json:"source"`
	WorkItem  struct {
		Key      string `json:"key"`
		Title    string `json:"title"`
		Provider string `json:"provider"`
		Project  string `json:"project"`
	} `json:"work_item"`
}

type claimResponse struct {
	ClaimID       string `json:"claim_id"`
	WorkItemKey   string `json:"work_item_key"`
	WorkItemTitle string `json:"work_item_title"`
	Provider      string `json:"provider"`
}

type offerRequest struct {
	SessionID string `json:"session_id"`
}

func (picker *Picker) handleOfferPicker(writer http.ResponseWriter, request *http.Request) {
	var body offerRequest
	if err := httputil.DecodeJSON(writer, request, &body); err != nil {
		return
	}
	if body.SessionID == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_session", "session_id is required")
		return
	}
	sessionID, err := picker.resolveSessionID(request.Context(), body.SessionID)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "offer_failed", err.Error())
		return
	}
	status, created, err := picker.store.OfferPicker(request.Context(), sessionID, time.Now().UnixMilli())
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "offer_failed", err.Error())
		return
	}
	if created {
		status = "offered"
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"status": status})
}

func (picker *Picker) resolveSessionID(ctx context.Context, raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	resolved, err := picker.store.ResolveSessionID(ctx, raw)
	if err != nil {
		return "", err
	}
	if resolved != nil {
		return *resolved, nil
	}
	return raw, nil
}

func (picker *Picker) handleCreateClaim(writer http.ResponseWriter, request *http.Request) {
	var body createClaimRequest
	if err := httputil.DecodeJSON(writer, request, &body); err != nil {
		return
	}
	if body.WorkItem.Key == "" || body.WorkItem.Provider == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_fields", "work_item.key and work_item.provider are required")
		return
	}

	sessionID, err := picker.resolveSessionID(request.Context(), ptr.Deref(body.SessionID))
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "session_resolve_failed", err.Error())
		return
	}

	nowUnixMs := time.Now().UnixMilli()
	claimScope := sessionID
	if claimScope == "" {
		claimScope = strconv.FormatInt(nowUnixMs, 10)
	}
	claimID := stableid.From("claim", body.WorkItem.Provider, body.WorkItem.Key, claimScope)

	source := body.Source
	if source == "" {
		source = "hook"
	}

	stored := store.Claim{
		ClaimID:         claimID,
		SessionID:       ptr.Optional(sessionID),
		WorkItemKey:     body.WorkItem.Key,
		WorkItemTitle:   body.WorkItem.Title,
		Provider:        body.WorkItem.Provider,
		Project:         body.WorkItem.Project,
		Source:          source,
		ClaimedAtUnixMs: nowUnixMs,
		UpdatedAtUnixMs: nowUnixMs,
	}
	if err := picker.store.UpsertClaim(request.Context(), stored); err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "claim_failed", err.Error())
		return
	}

	recent := store.RecentWorkItem{
		Key:                    body.WorkItem.Key,
		Provider:               body.WorkItem.Provider,
		Title:                  body.WorkItem.Title,
		Project:                body.WorkItem.Project,
		LastAttributedAtUnixMs: nowUnixMs,
	}
	if err := picker.store.UpsertRecentWorkItem(request.Context(), recent); err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "recent_failed", err.Error())
		return
	}

	httputil.WriteJSON(writer, http.StatusCreated, claimResponse{
		ClaimID:       claimID,
		WorkItemKey:   stored.WorkItemKey,
		WorkItemTitle: stored.WorkItemTitle,
		Provider:      stored.Provider,
	})
}

type attachClaimRequest struct {
	SessionID string `json:"session_id"`
}

func (picker *Picker) handleAttachClaim(writer http.ResponseWriter, request *http.Request) {
	var body attachClaimRequest
	if err := httputil.DecodeJSON(writer, request, &body); err != nil {
		return
	}
	if body.SessionID == "" {
		httputil.WriteError(writer, http.StatusBadRequest, "missing_session", "session_id is required")
		return
	}

	sessionID, err := picker.resolveSessionID(request.Context(), body.SessionID)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "session_resolve_failed", err.Error())
		return
	}

	if err := picker.store.AttachClaimToSession(
		request.Context(),
		request.PathValue("id"),
		sessionID,
		time.Now().UnixMilli(),
	); err != nil {
		if errors.Is(err, claim.ErrClaimNotFound) {
			httputil.WriteError(writer, http.StatusNotFound, "claim_not_found", "no claim matches the given id")
			return
		}
		httputil.WriteError(writer, http.StatusInternalServerError, "attach_failed", err.Error())
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"ok": true})
}
