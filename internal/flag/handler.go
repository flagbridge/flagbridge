package flag

import (
	"encoding/json"
	"net/http"

	"github.com/flagbridge/flagbridge/internal/auth"
	"github.com/flagbridge/flagbridge/internal/cache"
	"github.com/flagbridge/flagbridge/internal/environment"
	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/flagbridge/flagbridge/internal/sse"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc        *Service
	projectSvc *project.Service
	envSvc     *environment.Service
	cache      cache.Provider
	hub        *sse.Hub
}

func NewHandler(svc *Service, projectSvc *project.Service, envSvc *environment.Service, c cache.Provider, hub *sse.Hub) *Handler {
	return &Handler{svc: svc, projectSvc: projectSvc, envSvc: envSvc, cache: c, hub: hub}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	claims := auth.GetClaims(r.Context())
	f, err := h.svc.Create(r.Context(), p.ID, req, claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"data": f})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	flags, err := h.svc.ListByProject(r.Context(), p.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": flags})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	key := chi.URLParam(r, "key")
	f, err := h.svc.GetByKey(r.Context(), p.ID, key)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": f})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	key := chi.URLParam(r, "key")
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	f, err := h.svc.Update(r.Context(), p.ID, key, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": f})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	key := chi.URLParam(r, "key")
	if err := h.svc.Delete(r.Context(), p.ID, key); err != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"status": "deleted"}})
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	key := chi.URLParam(r, "key")
	envSlug := chi.URLParam(r, "env")

	f, err := h.svc.GetByKey(r.Context(), p.ID, key)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag not found")
		return
	}

	env, err := h.envSvc.GetBySlug(r.Context(), p.ID, envSlug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "environment not found")
		return
	}

	state, err := h.svc.GetState(r.Context(), f.ID, env.ID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag state not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": state})
}

func (h *Handler) SetState(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	key := chi.URLParam(r, "key")
	envSlug := chi.URLParam(r, "env")

	f, err := h.svc.GetByKey(r.Context(), p.ID, key)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag not found")
		return
	}

	env, err := h.envSvc.GetBySlug(r.Context(), p.ID, envSlug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "environment not found")
		return
	}

	var req SetStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	claims := auth.GetClaims(r.Context())
	state, err := h.svc.SetState(r.Context(), f.ID, env.ID, claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	// Invalidate evaluation cache for this flag
	h.cache.Invalidate(r.Context(), "eval:"+p.Slug+":"+envSlug+":"+key)

	// Broadcast SSE event
	h.hub.Broadcast(envSlug, sse.Event{
		Type: "flag.updated",
		Data: map[string]string{
			"flag_key":    key,
			"environment": envSlug,
		},
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": state})
}

func (h *Handler) resolveProject(w http.ResponseWriter, r *http.Request) *project.Project {
	slug := chi.URLParam(r, "slug")
	p, err := h.projectSvc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return nil
	}
	return p
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
