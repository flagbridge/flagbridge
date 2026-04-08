package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const (
	ClaimsContextKey  contextKey = "claims"
	APIKeyContextKey  contextKey = "apikey"
)

type APIKeyInfo struct {
	ID            string
	ProjectID     string
	EnvironmentID string
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
				db.Exec(context.Background(), `UPDATE api_keys SET last_used_at = now() WHERE id = $1`, info.ID)
			}()

			ctx := context.WithValue(r.Context(), APIKeyContextKey, &info)
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
	w.Write([]byte(`{"error":{"code":"unauthorized","message":"` + msg + `"}}`))
}
