package api

import (
	"net/http"
	"strings"

	"traceknot/internal/httputil"
	"traceknot/internal/store"
)

type SessionDetail struct {
	store *store.Store
}

func NewSessionDetail(storeHandle *store.Store) *SessionDetail {
	return &SessionDetail{store: storeHandle}
}

func (detail *SessionDetail) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/sessions/{id}/tree", detail.handleTree)
	mux.HandleFunc("GET /api/v1/nodes/{id}", detail.handleNode)
	return mux
}

func (detail *SessionDetail) handleTree(writer http.ResponseWriter, request *http.Request) {
	sessionID := request.PathValue("id")
	meta, err := detail.store.LoadSessionMeta(request.Context(), sessionID)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "tree_failed", err.Error())
		return
	}
	if meta == nil {
		httputil.WriteError(writer, http.StatusNotFound, "session_not_found", "session not found")
		return
	}
	rows, err := detail.store.LoadSessionTree(request.Context(), sessionID)
	if err != nil {
		httputil.WriteError(writer, http.StatusInternalServerError, "tree_failed", err.Error())
		return
	}
	nodes := make([]treeRowJSON, 0, len(rows))
	for _, row := range rows {
		nodes = append(nodes, toTreeRowJSON(row))
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{
		"session": metaJSON(meta),
		"nodes":   nodes,
	})
}

func (detail *SessionDetail) handleNode(writer http.ResponseWriter, request *http.Request) {
	nodeID := request.PathValue("id")
	var entry *store.NodeDetail
	if strings.HasPrefix(nodeID, "synthetic:root:") {
		sessionID := strings.TrimPrefix(nodeID, "synthetic:root:")
		rows, err := detail.store.LoadSessionTree(request.Context(), sessionID)
		if err != nil {
			httputil.WriteError(writer, http.StatusInternalServerError, "node_failed", err.Error())
			return
		}
		entry = syntheticRootDetail(rows, sessionID)
	} else {
		var err error
		entry, err = detail.store.LoadNodeDetail(request.Context(), nodeID)
		if err != nil {
			httputil.WriteError(writer, http.StatusInternalServerError, "node_failed", err.Error())
			return
		}
	}
	if entry == nil {
		httputil.WriteError(writer, http.StatusNotFound, "node_not_found", "node not found")
		return
	}
	httputil.WriteJSON(writer, http.StatusOK, map[string]any{"node": nodeDetailJSON(entry)})
}
