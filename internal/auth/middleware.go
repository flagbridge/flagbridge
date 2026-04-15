package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const (
	ClaimsContextKey       contextKey = "claims"
	APIKeyContextKey       contextKey = "apikey"
	ProjectRoleContextKey  contextKey = "project_role"
)

type APIKeyInfo struct {
	ID            string
	ProjectID     string
	EnvironmentID *string
	Scope         string
}

// JWTMiddleware validates JWT tokens from the Authorization header.
func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				unauthorizedJSON(w, "missing authorization token")
				return
			}

			claims, err := ValidateToken(secret, token)
			if err != nil {
				unauthorizedJSON(w, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// APIKeyMiddleware validates API keys from the Authorization header.
func APIKeyMiddleware(db *pgxpool.Pool, requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := extractBearerToken(r)
			if key == "" || !strings.HasPrefix(key, "fb_sk_") {
				unauthorizedJSON(w, "missing or invalid API key")
				return
			}

			hash := HashAPIKey(key)

			var info APIKeyInfo
			err := db.QueryRow(r.Context(), `
				SELECT id, project_id, environment_id, scope
				FROM api_keys
				WHERE key_hash = $1 AND (expires_at IS NULL OR expires_at > now())
			`, hash).Scan(&info.ID, &info.ProjectID, &info.EnvironmentID, &info.Scope)
			if err != nil {
				unauthorizedJSON(w, "invalid API key")
				return
			}

			if requiredScope != "" && info.Scope != requiredScope && info.Scope != "full" {
				unauthorizedJSON(w, "insufficient API key scope")
				return
			}

			// Update last_used_at
			go func() {
				_, _ = db.Exec(context.Background(), `UPDATE api_keys SET last_used_at = now() WHERE id = $1`, info.ID)
			}()

			ctx := context.WithValue(r.Context(), APIKeyContextKey, &info)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireProjectRole checks that the authenticated user has one of the allowed roles
// in the project identified by the {slug} URL parameter. Global admins (users.role="admin")
// bypass role checks. The resolved project role is stored in context for downstream use.
func RequireProjectRole(db *pgxpool.Pool, allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				unauthorizedJSON(w, "authentication required")
				return
			}

			// Global admins bypass project role checks
			if claims.Role == "admin" {
				ctx := context.WithValue(r.Context(), ProjectRoleContextKey, "admin")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			slug := chi.URLParam(r, "slug")
			if slug == "" {
				forbiddenJSON(w, "project context required")
				return
			}

			var projectRole string
			err := db.QueryRow(r.Context(), `
				SELECT pm.role FROM project_members pm
				JOIN projects p ON p.id = pm.project_id
				WHERE p.slug = $1 AND pm.user_id = $2
			`, slug, claims.UserID).Scan(&projectRole)
			if err == pgx.ErrNoRows {
				forbiddenJSON(w, "not a member of this project")
				return
			}
			if err != nil {
				slog.Error("failed to check project role", "error", err)
				forbiddenJSON(w, "failed to verify permissions")
				return
			}

			if !allowed[projectRole] {
				forbiddenJSON(w, "insufficient role: requires "+strings.Join(allowedRoles, " or "))
				return
			}

			ctx := context.WithValue(r.Context(), ProjectRoleContextKey, projectRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) *Claims {
	claims, _ := ctx.Value(ClaimsContextKey).(*Claims)
	return claims
}

func GetAPIKeyInfo(ctx context.Context) *APIKeyInfo {
	info, _ := ctx.Value(APIKeyContextKey).(*APIKeyInfo)
	return info
}

func GetProjectRole(ctx context.Context) string {
	role, _ := ctx.Value(ProjectRoleContextKey).(string)
	return role
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

func unauthorizedJSON(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"` + msg + `"}}`))
}

func forbiddenJSON(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":{"code":"forbidden","message":"` + msg + `"}}`))
}
