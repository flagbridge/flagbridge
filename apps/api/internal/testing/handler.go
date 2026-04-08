package testing

import (
	"encoding/json"
	"net/http"

	"github.com/flagbridge/flagbridge/apps/api/internal/auth"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	// Get user from API key or JWT context
	userID := resolveUserID(r)

	session, err := h.svc.CreateSession(r.Context(), req, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"data": session})
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sessionID")
	session, err := h.svc.GetSession(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "session not found or expired")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": session})
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "sessionID")
	if err := h.svc.DeleteSession(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing_project", "project_id query param is required")
		return
	}

	sessions, err := h.svc.ListSessions(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": sessions})
}

func (h *Handler) SetOverride(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	var req SetOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	session, err := h.svc.SetOverride(r.Context(), sessionID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "override_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": session})
}

func (h *Handler) SetOverridesBatch(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	var req SetOverridesBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	session, err := h.svc.SetOverridesBatch(r.Context(), sessionID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "override_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": session})
}

func (h *Handler) DeleteOverride(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	flagKey := chi.URLParam(r, "flagKey")

	session, err := h.svc.DeleteOverride(r.Context(), sessionID, flagKey)
	if err != nil {
		writeError(w, http.StatusBadRequest, "delete_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": session})
}

func resolveUserID(r *http.Request) string {
	if claims := auth.GetClaims(r.Context()); claims != nil {
		return claims.UserID
	}
	if info := auth.GetAPIKeyInfo(r.Context()); info != nil {
		return "apikey:" + info.ID
	}
	return "unknown"
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}
