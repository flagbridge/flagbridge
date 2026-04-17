package aiconfig

import (
	"encoding/json"
	"net/http"

	"github.com/flagbridge/flagbridge/internal/audit"
	"github.com/flagbridge/flagbridge/internal/auth"
	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc        *Service
	projectSvc *project.Service
	auditSvc   *audit.Service
}

func NewHandler(svc *Service, projectSvc *project.Service, auditSvc *audit.Service) *Handler {
	return &Handler{svc: svc, projectSvc: projectSvc, auditSvc: auditSvc}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.resolveProject(w, r)
	if err != nil {
		return
	}

	provider, err := h.svc.Get(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "AI provider config not found for this project")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": provider})
}

func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.resolveProject(w, r)
	if err != nil {
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	provider, err := h.svc.Create(r.Context(), projectID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	claims := auth.GetClaims(r.Context())
	slug := chi.URLParam(r, "slug")
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID: projectID, UserID: claims.UserID, Action: "configured",
		EntityType: "ai_provider", EntityID: provider.ID,
		Changes:   map[string]any{"provider": req.Provider, "model": provider.ModelID, "project": slug},
		IPAddress: r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": provider})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.resolveProject(w, r)
	if err != nil {
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	provider, err := h.svc.Update(r.Context(), projectID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}

	claims := auth.GetClaims(r.Context())
	slug := chi.URLParam(r, "slug")
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID: projectID, UserID: claims.UserID, Action: "updated",
		EntityType: "ai_provider", EntityID: provider.ID,
		Changes:   map[string]any{"provider": provider.ProviderName, "model": provider.ModelID, "project": slug},
		IPAddress: r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": provider})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.resolveProject(w, r)
	if err != nil {
		return
	}

	if err := h.svc.Delete(r.Context(), projectID); err != nil {
		writeError(w, http.StatusNotFound, "not_found", "AI provider config not found")
		return
	}

	claims := auth.GetClaims(r.Context())
	slug := chi.URLParam(r, "slug")
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID: projectID, UserID: claims.UserID, Action: "deleted",
		EntityType: "ai_provider", EntityID: projectID,
		Changes:   map[string]any{"project": slug},
		IPAddress: r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"status": "deleted"}})
}

func (h *Handler) resolveProject(w http.ResponseWriter, r *http.Request) (string, error) {
	slug := chi.URLParam(r, "slug")
	p, err := h.projectSvc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return "", err
	}
	return p.ID, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}
