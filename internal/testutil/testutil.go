//go:build integration

package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/flagbridge/flagbridge/internal/apikey"
	"github.com/flagbridge/flagbridge/internal/audit"
	"github.com/flagbridge/flagbridge/internal/auth"
	"github.com/flagbridge/flagbridge/internal/cache"
	"github.com/flagbridge/flagbridge/internal/environment"
	"github.com/flagbridge/flagbridge/internal/evaluation"
	"github.com/flagbridge/flagbridge/internal/flag"
	"github.com/flagbridge/flagbridge/internal/middleware"
	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/flagbridge/flagbridge/internal/sse"
	fbtesting "github.com/flagbridge/flagbridge/internal/testing"
	"github.com/flagbridge/flagbridge/internal/targeting"
	"github.com/flagbridge/flagbridge/internal/webhook"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const TestJWTSecret = "integration-test-secret-do-not-use-in-prod"

// SeedUserID is the UUID of the seeded admin user from migrations/001_initial.sql.
const SeedUserID = "00000000-0000-0000-0000-000000000001"

// SeedUserEmail is the email of the seeded admin user.
const SeedUserEmail = "admin@flagbridge.io"

// SeedUserPassword is the plain-text password for the seeded admin user.
const SeedUserPassword = "flagbridge-admin-2026"

// SetupTestDB connects to a PostgreSQL database using DATABASE_URL.
// If DATABASE_URL is not set the test is skipped automatically.
// The caller is responsible for calling TeardownTestDB when done.
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}

	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		t.Fatalf("testutil: parse DATABASE_URL: %v", err)
	}
	cfg.MaxConns = 5
	cfg.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("testutil: create pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("testutil: ping database: %v", err)
	}

	return pool
}

// TeardownTestDB closes the pool and clears test-owned rows so tests are isolated.
// It truncates all tables in dependency order, preserving the seed admin user.
func TeardownTestDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE
			webhook_delivery_logs,
			webhooks,
			targeting_rules,
			flag_states,
			flags,
			api_keys,
			testing_sessions,
			audit_log,
			environments,
			projects
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Errorf("testutil: truncate tables: %v", err)
	}

	pool.Close()
}

// CreateTestServer builds a full Chi router wired up exactly as main.go does, backed
// by the provided pool. It returns both the test server and a cancel function that
// stops the server.
func CreateTestServer(t *testing.T, pool *pgxpool.Pool) *httptest.Server {
	t.Helper()

	memCache, err := cache.NewMemoryCache()
	if err != nil {
		t.Fatalf("testutil: init cache: %v", err)
	}

	hub := sse.NewHub()

	// Repositories
	projectRepo := project.NewRepository(pool)
	envRepo := environment.NewRepository(pool)
	flagRepo := flag.NewRepository(pool)
	apikeyRepo := apikey.NewRepository(pool)
	auditRepo := audit.NewRepository(pool)
	testingRepo := fbtesting.NewRepository(pool)
	webhookRepo := webhook.NewRepository(pool)

	// Services
	projectSvc := project.NewService(projectRepo)
	envSvc := environment.NewService(envRepo)
	flagSvc := flag.NewService(flagRepo)
	apikeySvc := apikey.NewService(apikeyRepo)
	auditSvc := audit.NewService(auditRepo)
	testingSvc := fbtesting.NewService(testingRepo)
	webhookDispatcher := webhook.NewDispatcher(webhookRepo)
	webhookSvc := webhook.NewService(webhookRepo, webhookDispatcher)

	// Handlers
	projectHandler := project.NewHandler(projectSvc)
	envHandler := environment.NewHandler(envSvc, projectSvc)
	targetingRepo := targeting.NewRepository(pool)
	targetingSvc := targeting.NewService(targetingRepo)
	flagHandler := flag.NewHandler(flagSvc, projectSvc, envSvc, targetingSvc, auditSvc, memCache, hub)
	evalHandler := evaluation.NewHandler(pool, memCache)
	apikeyHandler := apikey.NewHandler(apikeySvc)
	auditHandler := audit.NewHandler(auditSvc)
	testingHandler := fbtesting.NewHandler(testingSvc)
	webhookHandler := webhook.NewHandler(webhookSvc)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})
		r.Post("/auth/login", loginHandler(pool, TestJWTSecret))

		// Admin routes (JWT auth)
		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware(TestJWTSecret))

			r.Post("/projects", projectHandler.Create)
			r.Get("/projects", projectHandler.List)
			r.Get("/projects/{slug}", projectHandler.GetBySlug)
			r.Patch("/projects/{slug}", projectHandler.Update)
			r.Delete("/projects/{slug}", projectHandler.Delete)

			r.Post("/projects/{slug}/environments", envHandler.Create)
			r.Get("/projects/{slug}/environments", envHandler.List)

			r.Post("/projects/{slug}/flags", flagHandler.Create)
			r.Get("/projects/{slug}/flags", flagHandler.List)
			r.Get("/projects/{slug}/flags/{key}", flagHandler.Get)
			r.Patch("/projects/{slug}/flags/{key}", flagHandler.Update)
			r.Delete("/projects/{slug}/flags/{key}", flagHandler.Delete)

			r.Put("/projects/{slug}/flags/{key}/states/{env}", flagHandler.SetState)
			r.Get("/projects/{slug}/flags/{key}/states/{env}", flagHandler.GetState)

			r.Post("/api-keys", apikeyHandler.Create)
			r.Get("/api-keys", apikeyHandler.List)
			r.Delete("/api-keys/{id}", apikeyHandler.Delete)

			r.Post("/projects/{slug}/webhooks", webhookHandler.Create)
			r.Get("/projects/{slug}/webhooks", webhookHandler.List)
			r.Get("/webhooks/{webhookID}", webhookHandler.Get)
			r.Patch("/webhooks/{webhookID}", webhookHandler.Update)
			r.Delete("/webhooks/{webhookID}", webhookHandler.Delete)
			r.Get("/webhooks/{webhookID}/logs", webhookHandler.DeliveryLogs)
			r.Post("/webhooks/{webhookID}/test", webhookHandler.Test)

			r.Get("/audit-log", auditHandler.List)
		})

		// SDK routes (API key auth — eval scope)
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(pool, "eval"))
			r.Post("/evaluate", evalHandler.Evaluate)
			r.Post("/evaluate/batch", evalHandler.BatchEvaluate)
		})

		// Testing API routes (API key auth — test scope)
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(pool, "test"))
			r.Post("/testing/sessions", testingHandler.CreateSession)
			r.Get("/testing/sessions", testingHandler.ListSessions)
			r.Get("/testing/sessions/{sessionID}", testingHandler.GetSession)
			r.Delete("/testing/sessions/{sessionID}", testingHandler.DeleteSession)
			r.Put("/testing/sessions/{sessionID}/overrides", testingHandler.SetOverride)
			r.Put("/testing/sessions/{sessionID}/overrides/batch", testingHandler.SetOverridesBatch)
			r.Delete("/testing/sessions/{sessionID}/overrides/{flagKey}", testingHandler.DeleteOverride)
		})
	})

	return httptest.NewServer(r)
}

// MustLogin calls POST /v1/auth/login and returns the JWT token.
// It fatally fails the test if login does not succeed.
func MustLogin(t *testing.T, srv *httptest.Server) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{
		"email":    SeedUserEmail,
		"password": SeedUserPassword,
	})

	resp, err := srv.Client().Post(srv.URL+"/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("testutil: login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("testutil: login returned %d", resp.StatusCode)
	}

	var out struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("testutil: decode login response: %v", err)
	}
	if out.Data.Token == "" {
		t.Fatal("testutil: login returned empty token")
	}
	return out.Data.Token
}

// MustCreateProject creates a project via the API and returns the created project data.
// It fatally fails the test on any error.
func MustCreateProject(t *testing.T, srv *httptest.Server, token, name, slug string) map[string]any {
	t.Helper()

	body, _ := json.Marshal(map[string]string{
		"name": name,
		"slug": slug,
	})

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("testutil: create project request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("testutil: create project returned %d", resp.StatusCode)
	}

	var out struct {
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("testutil: decode create project response: %v", err)
	}
	return out.Data
}

// MustCreateFlag creates a flag inside a project via the API and returns the created flag data.
func MustCreateFlag(t *testing.T, srv *httptest.Server, token, projectSlug, key, name string) map[string]any {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"key":           key,
		"name":          name,
		"type":          "boolean",
		"default_value": false,
	})

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/projects/%s/flags", srv.URL, projectSlug), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("testutil: create flag request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("testutil: create flag returned %d (project=%s, key=%s)", resp.StatusCode, projectSlug, key)
	}

	var out struct {
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("testutil: decode create flag response: %v", err)
	}
	return out.Data
}

// AuthedRequest builds an *http.Request pre-filled with the Bearer token.
func AuthedRequest(t *testing.T, method, url, token string, body any) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("testutil: encode request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		t.Fatalf("testutil: build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

// DecodeBody decodes the response body into the provided target.
func DecodeBody(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("testutil: decode response body: %v", err)
	}
}

// loginHandler is a local copy of the one in main.go so the test server does not
// import the main package (which would create a dependency cycle).
func loginHandler(db *pgxpool.Pool, jwtSecret string) http.HandlerFunc {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
			return
		}
		if req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "missing_fields", "email and password are required")
			return
		}

		var userID, email, name, role, passwordHash string
		err := db.QueryRow(r.Context(), `
			SELECT id, email, name, role, password FROM users WHERE email = $1
		`, req.Email).Scan(&userID, &email, &name, &role, &passwordHash)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
			return
		}

		if !auth.CheckPassword(passwordHash, req.Password) {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
			return
		}

		token, err := auth.GenerateToken(jwtSecret, userID, email, role)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "token_error", "failed to generate token")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"token": token,
				"user": map[string]string{
					"id":    userID,
					"email": email,
					"name":  name,
					"role":  role,
				},
			},
		})
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}
