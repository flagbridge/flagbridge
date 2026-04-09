package audit

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		// Try resolving from route param if nested under /projects/{slug}
		slug := chi.URLParam(r, "slug")
		if slug != "" {
			projectID = slug // will be resolved by caller
		}
	}

	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing_project", "project_id query param is required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	entries, total, err := h.svc.List(r.Context(), ListParams{
		ProjectID:  projectID,
		UserID:     r.URL.Query().Get("user_id"),
		Action:     r.URL.Query().Get("action"),
		EntityType: r.URL.Query().Get("entity_type"),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data":  entries,
		"total": total,
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}
