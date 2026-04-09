package environment

import (
	"encoding/json"
	"net/http"

	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc        *Service
	projectSvc *project.Service
}

func NewHandler(svc *Service, projectSvc *project.Service) *Handler {
	return &Handler{svc: svc, projectSvc: projectSvc}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	p, err := h.projectSvc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	env, err := h.svc.Create(r.Context(), p.ID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"data": env})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	p, err := h.projectSvc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return
	}

	envs, err := h.svc.ListByProject(r.Context(), p.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": envs})
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
