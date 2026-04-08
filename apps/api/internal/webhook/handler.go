package webhook

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/flagbridge/flagbridge/apps/api/internal/auth"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectSlug := chi.URLParam(r, "slug")

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	claims := auth.GetClaims(r.Context())
	// projectSlug is used as projectID here; caller resolves slug→ID before routing
	wh, err := h.svc.Create(r.Context(), projectSlug, req, claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	// Return secret only on creation
	writeJSON(w, http.StatusCreated, map[string]any{
		"data": map[string]any{
			"id":         wh.ID,
			"project_id": wh.ProjectID,
			"url":        wh.URL,
			"secret":     wh.Secret,
			"events":     wh.Events,
			"active":     wh.Active,
			"created_at": wh.CreatedAt,
		},
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projectSlug := chi.URLParam(r, "slug")

	webhooks, err := h.svc.ListByProject(r.Context(), projectSlug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": webhooks})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhookID")
	wh, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "webhook not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": wh})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhookID")

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	wh, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": wh})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhookID")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeliveryLogs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhookID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	logs, err := h.svc.ListDeliveryLogs(r.Context(), id, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": logs})
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "webhookID")
	if err := h.svc.TestWebhook(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"status": "test_sent"}})
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
