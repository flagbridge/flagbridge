package member

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

func NewHandler(svc *Service, ps *project.Service, as *audit.Service) *Handler {
	return &Handler{svc: svc, projectSvc: ps, auditSvc: as}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	members, err := h.svc.List(r.Context(), p.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": members})
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	m, err := h.svc.Add(r.Context(), p.ID, req.Email, req.Role)
	if err != nil {
		writeError(w, http.StatusBadRequest, "add_failed", err.Error())
		return
	}

	claims := auth.GetClaims(r.Context())
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID:  p.ID,
		UserID:     claims.UserID,
		Action:     "member_added",
		EntityType: "project_member",
		EntityID:   m.ID,
		Changes:    map[string]any{"email": m.Email, "role": m.Role},
		IPAddress:  r.RemoteAddr,
	})

	writeJSON(w, http.StatusCreated, map[string]any{"data": m})
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	memberID := chi.URLParam(r, "memberID")

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	m, err := h.svc.UpdateRole(r.Context(), memberID, req.Role)
	if err != nil {
		writeError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}

	claims := auth.GetClaims(r.Context())
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID:  p.ID,
		UserID:     claims.UserID,
		Action:     "member_role_changed",
		EntityType: "project_member",
		EntityID:   m.ID,
		Changes:    map[string]any{"email": m.Email, "new_role": m.Role},
		IPAddress:  r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": m})
}

func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	p := h.resolveProject(w, r)
	if p == nil {
		return
	}

	memberID := chi.URLParam(r, "memberID")

	if err := h.svc.Remove(r.Context(), memberID); err != nil {
		writeError(w, http.StatusNotFound, "not_found", "member not found")
		return
	}

	claims := auth.GetClaims(r.Context())
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID:  p.ID,
		UserID:     claims.UserID,
		Action:     "member_removed",
		EntityType: "project_member",
		EntityID:   memberID,
		Changes:    nil,
		IPAddress:  r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"deleted": "true"}})
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}
