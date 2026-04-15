package productcard

import (
	"encoding/json"
	"net/http"

	"github.com/flagbridge/flagbridge/internal/audit"
	"github.com/flagbridge/flagbridge/internal/auth"
	"github.com/flagbridge/flagbridge/internal/flag"
	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc        *Service
	projectSvc *project.Service
	flagSvc    *flag.Service
	auditSvc   *audit.Service
}

func NewHandler(svc *Service, ps *project.Service, fs *flag.Service, as *audit.Service) *Handler {
	return &Handler{svc: svc, projectSvc: ps, flagSvc: fs, auditSvc: as}
}

// Get returns the product card for a flag, or null if none exists.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p, f := h.resolve(w, r)
	if p == nil {
		return
	}

	card, err := h.svc.Get(r.Context(), f.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	if card == nil {
		writeJSON(w, http.StatusOK, map[string]any{"data": nil})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": card})
}

// Upsert creates or updates the product card for a flag.
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	p, f := h.resolve(w, r)
	if p == nil {
		return
	}

	var req UpsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	card, err := h.svc.Upsert(r.Context(), f.ID, p.ID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "upsert_failed", err.Error())
		return
	}

	claims := auth.GetClaims(r.Context())
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID:  p.ID,
		UserID:     claims.UserID,
		Action:     "updated",
		EntityType: "product_card",
		EntityID:   card.ID,
		Changes: map[string]any{
			"flag_key":   f.Key,
			"hypothesis": card.Hypothesis,
			"status":     card.Status,
		},
		IPAddress: r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": card})
}

// Delete removes the product card from a flag.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	p, f := h.resolve(w, r)
	if p == nil {
		return
	}

	if err := h.svc.Delete(r.Context(), f.ID); err != nil {
		writeError(w, http.StatusNotFound, "not_found", "product card not found")
		return
	}

	claims := auth.GetClaims(r.Context())
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID:  p.ID,
		UserID:     claims.UserID,
		Action:     "deleted",
		EntityType: "product_card",
		EntityID:   f.ID,
		Changes:    map[string]any{"flag_key": f.Key},
		IPAddress:  r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]string{"deleted": "true"}})
}

func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) (*project.Project, *flag.Flag) {
	slug := chi.URLParam(r, "slug")
	key := chi.URLParam(r, "key")

	p, err := h.projectSvc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return nil, nil
	}

	f, err := h.flagSvc.GetByKey(r.Context(), p.ID, key)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag not found")
		return nil, nil
	}

	return p, f
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
