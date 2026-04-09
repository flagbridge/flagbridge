package evaluation

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/flagbridge/flagbridge/internal/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db    *pgxpool.Pool
	cache cache.Provider
}

func NewHandler(db *pgxpool.Pool, c cache.Provider) *Handler {
	return &Handler{db: db, cache: c}
}

type EvaluateRequest struct {
	FlagKey     string      `json:"flag_key"`
	Project     string      `json:"project"`
	Environment string      `json:"environment"`
	Context     EvalContext  `json:"context"`
}

type BatchEvaluateRequest struct {
	Project     string      `json:"project"`
	Environment string      `json:"environment"`
	FlagKeys    []string    `json:"flag_keys"`
	Context     EvalContext  `json:"context"`
}

func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	if req.FlagKey == "" || req.Project == "" || req.Environment == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "flag_key, project, and environment are required")
		return
	}

	result, err := h.evaluateFlag(r.Context(), req.Project, req.Environment, req.FlagKey, req.Context)
	if err != nil {
		if err == errFlagNotFound {
			writeError(w, http.StatusNotFound, "flag_not_found", fmt.Sprintf("flag %q not found", req.FlagKey))
			return
		}
		slog.Error("evaluation failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "evaluation failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": result})
}

func (h *Handler) BatchEvaluate(w http.ResponseWriter, r *http.Request) {
	var req BatchEvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	if req.Project == "" || req.Environment == "" || len(req.FlagKeys) == 0 {
		writeError(w, http.StatusBadRequest, "missing_fields", "project, environment, and flag_keys are required")
		return
	}

	flags := make(map[string]*Result, len(req.FlagKeys))
	for _, key := range req.FlagKeys {
		result, err := h.evaluateFlag(r.Context(), req.Project, req.Environment, key, req.Context)
		if err != nil {
			if err == errFlagNotFound {
				continue
			}
			slog.Error("batch evaluation failed", "flag_key", key, "error", err)
			continue
		}
		flags[key] = result
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"flags": flags}})
}

var errFlagNotFound = fmt.Errorf("flag not found")

func (h *Handler) evaluateFlag(ctx context.Context, project, env, flagKey string, evalCtx EvalContext) (*Result, error) {
	// Check cache
	cacheKey := buildCacheKey(project, env, flagKey, evalCtx)
	if cached, ok := h.cache.Get(ctx, cacheKey); ok {
		var result Result
		if err := json.Unmarshal(cached, &result); err == nil {
			return &result, nil
		}
	}

	// Fetch flag data from DB
	data, err := h.fetchFlagData(ctx, project, env, flagKey)
	if err != nil {
		return nil, err
	}

	result := Evaluate(flagKey, *data, evalCtx)

	// Cache result
	if b, err := json.Marshal(result); err == nil {
		h.cache.Set(ctx, cacheKey, b, 10*time.Second)
	}

	return &result, nil
}

func (h *Handler) fetchFlagData(ctx context.Context, projectSlug, envSlug, flagKey string) (*FlagData, error) {
	var data FlagData
	var flagID, envID string

	err := h.db.QueryRow(ctx, `
		SELECT f.id, f.default_value, fs.enabled, fs.value, e.id
		FROM flags f
		JOIN projects p ON p.id = f.project_id
		JOIN environments e ON e.project_id = p.id AND e.slug = $2
		LEFT JOIN flag_states fs ON fs.flag_id = f.id AND fs.environment_id = e.id
		WHERE p.slug = $1 AND f.key = $3
	`, projectSlug, envSlug, flagKey).Scan(
		&flagID, &data.DefaultValue, &data.Enabled, &data.StateValue, &envID,
	)
	if err != nil {
		return nil, errFlagNotFound
	}

	rows, err := h.db.Query(ctx, `
		SELECT id, priority, conditions, value, enabled
		FROM targeting_rules
		WHERE flag_id = $1 AND environment_id = $2
		ORDER BY priority ASC
	`, flagID, envID)
	if err != nil {
		return nil, fmt.Errorf("fetching targeting rules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rule Rule
		var conditionsJSON []byte
		if err := rows.Scan(&rule.ID, &rule.Priority, &conditionsJSON, &rule.Value, &rule.Enabled); err != nil {
			return nil, fmt.Errorf("scanning targeting rule: %w", err)
		}
		if err := json.Unmarshal(conditionsJSON, &rule.Conditions); err != nil {
			slog.Warn("invalid conditions JSON", "rule_id", rule.ID, "error", err)
			continue
		}
		data.Rules = append(data.Rules, rule)
	}

	return &data, nil
}

func buildCacheKey(project, env, flagKey string, ctx EvalContext) string {
	ctxJSON, _ := json.Marshal(ctx)
	hash := sha256.Sum256(ctxJSON)
	return fmt.Sprintf("eval:%s:%s:%s:%x", project, env, flagKey, hash)
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
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
