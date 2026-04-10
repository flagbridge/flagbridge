package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flagbridge/flagbridge/internal/apikey"
	"github.com/flagbridge/flagbridge/internal/audit"
	"github.com/flagbridge/flagbridge/internal/auth"
	"github.com/flagbridge/flagbridge/internal/cache"
	"github.com/flagbridge/flagbridge/internal/config"
	"github.com/flagbridge/flagbridge/internal/database"
	"github.com/flagbridge/flagbridge/internal/environment"
	"github.com/flagbridge/flagbridge/internal/evaluation"
	"github.com/flagbridge/flagbridge/internal/flag"
	fbtesting "github.com/flagbridge/flagbridge/internal/testing"
	"github.com/flagbridge/flagbridge/internal/webhook"
	"github.com/flagbridge/flagbridge/internal/middleware"
	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/flagbridge/flagbridge/internal/sse"
	"github.com/flagbridge/flagbridge/internal/targeting"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Structured JSON logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(1)
	}
	if cfg.JWTSecret == "" {
		slog.Error("JWT_SECRET is required")
		os.Exit(1)
	}

	ctx := context.Background()

	// Database
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Cache
	memCache, err := cache.NewMemoryCache()
	if err != nil {
		slog.Error("failed to initialize cache", "error", err)
		os.Exit(1)
	}

	// SSE Hub
	hub := sse.NewHub()

	// Repositories
	projectRepo := project.NewRepository(db)
	envRepo := environment.NewRepository(db)
	flagRepo := flag.NewRepository(db)
	targetingRepo := targeting.NewRepository(db)
	apikeyRepo := apikey.NewRepository(db)
	auditRepo := audit.NewRepository(db)
	testingRepo := fbtesting.NewRepository(db)
	webhookRepo := webhook.NewRepository(db)

	// Services
	projectSvc := project.NewService(projectRepo)
	envSvc := environment.NewService(envRepo)
	flagSvc := flag.NewService(flagRepo)
	targetingSvc := targeting.NewService(targetingRepo)
	apikeySvc := apikey.NewService(apikeyRepo)
	auditSvc := audit.NewService(auditRepo)
	testingSvc := fbtesting.NewService(testingRepo)
	webhookDispatcher := webhook.NewDispatcher(webhookRepo)
	webhookSvc := webhook.NewService(webhookRepo, webhookDispatcher)

	// Handlers
	projectHandler := project.NewHandler(projectSvc)
	envHandler := environment.NewHandler(envSvc, projectSvc)
	flagHandler := flag.NewHandler(flagSvc, projectSvc, envSvc, targetingSvc, auditSvc, memCache, hub)
	evalHandler := evaluation.NewHandler(db, memCache)
	apikeyHandler := apikey.NewHandler(apikeySvc)
	auditHandler := audit.NewHandler(auditSvc)
	testingHandler := fbtesting.NewHandler(testingSvc)
	webhookHandler := webhook.NewHandler(webhookSvc)


	// Build targeting handler inline since it needs cross-package deps
	targetingHandler := newTargetingHandler(targetingSvc, projectSvc, flagSvc, envSvc, auditSvc, memCache, hub)

	// Router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Recovery)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// All /v1 routes
	r.Route("/v1", func(r chi.Router) {
		// Public
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})
		r.Post("/auth/login", loginHandler(db, cfg.JWTSecret))

		// Admin routes (JWT auth)
		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware(cfg.JWTSecret))

			// Projects
			r.Post("/projects", projectHandler.Create)
			r.Get("/projects", projectHandler.List)
			r.Get("/projects/{slug}", projectHandler.GetBySlug)
			r.Patch("/projects/{slug}", projectHandler.Update)
			r.Delete("/projects/{slug}", projectHandler.Delete)

			// Environments
			r.Post("/projects/{slug}/environments", envHandler.Create)
			r.Get("/projects/{slug}/environments", envHandler.List)

			// Flags
			r.Post("/projects/{slug}/flags", flagHandler.Create)
			r.Get("/projects/{slug}/flags", flagHandler.List)
			r.Get("/projects/{slug}/flags/{key}", flagHandler.Get)
			r.Patch("/projects/{slug}/flags/{key}", flagHandler.Update)
			r.Delete("/projects/{slug}/flags/{key}", flagHandler.Delete)

			// Flag states
			r.Put("/projects/{slug}/flags/{key}/states/{env}", flagHandler.SetState)
			r.Get("/projects/{slug}/flags/{key}/states/{env}", flagHandler.GetState)

			// Targeting rules
			r.Put("/projects/{slug}/flags/{key}/rules/{env}", targetingHandler.SetRules)
			r.Get("/projects/{slug}/flags/{key}/rules/{env}", targetingHandler.GetRules)

			// API keys
			r.Post("/api-keys", apikeyHandler.Create)
			r.Get("/api-keys", apikeyHandler.List)
			r.Delete("/api-keys/{id}", apikeyHandler.Delete)

			// Webhooks
			r.Post("/projects/{slug}/webhooks", webhookHandler.Create)
			r.Get("/projects/{slug}/webhooks", webhookHandler.List)
			r.Get("/projects/{slug}/webhooks/{webhook_id}", webhookHandler.Get)
			r.Patch("/projects/{slug}/webhooks/{webhook_id}", webhookHandler.Update)
			r.Delete("/projects/{slug}/webhooks/{webhook_id}", webhookHandler.Delete)
			r.Get("/projects/{slug}/webhooks/{webhook_id}/logs", webhookHandler.DeliveryLogs)
			r.Post("/projects/{slug}/webhooks/{webhook_id}/test", webhookHandler.Test)

			// Audit log
			r.Get("/audit-log", auditHandler.List)
		})

		// SDK routes (API key auth — eval scope)
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(db, "eval"))

			r.Post("/evaluate", evalHandler.Evaluate)
			r.Post("/evaluate/batch", evalHandler.BatchEvaluate)
		})

		// Testing API routes (API key auth — test scope)
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(db, "test"))

			r.Post("/testing/sessions", testingHandler.CreateSession)
			r.Get("/testing/sessions", testingHandler.ListSessions)
			r.Get("/testing/sessions/{sessionID}", testingHandler.GetSession)
			r.Delete("/testing/sessions/{sessionID}", testingHandler.DeleteSession)
			r.Put("/testing/sessions/{sessionID}/overrides", testingHandler.SetOverride)
			r.Put("/testing/sessions/{sessionID}/overrides/batch", testingHandler.SetOverridesBatch)
			r.Delete("/testing/sessions/{sessionID}/overrides/{flagKey}", testingHandler.DeleteOverride)
		})

		// SSE routes (API key auth — eval scope)
		r.Group(func(r chi.Router) {
			r.Use(auth.APIKeyMiddleware(db, "eval"))
			r.Get("/sse/{environment}", hub.ServeHTTP)
		})
	})

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("starting flagbridge api server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}

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
			slog.Error("failed to generate token", "error", err)
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

// targetingHandlerWrapper bridges the targeting package with project/flag/env resolution.
type targetingHandlerWrapper struct {
	svc        *targeting.Service
	projectSvc *project.Service
	flagSvc    *flag.Service
	envSvc     *environment.Service
	auditSvc   *audit.Service
	cache      cache.Provider
	hub        *sse.Hub
}

func newTargetingHandler(svc *targeting.Service, ps *project.Service, fs *flag.Service, es *environment.Service, as *audit.Service, c cache.Provider, h *sse.Hub) *targetingHandlerWrapper {
	return &targetingHandlerWrapper{svc: svc, projectSvc: ps, flagSvc: fs, envSvc: es, auditSvc: as, cache: c, hub: h}
}

func (h *targetingHandlerWrapper) GetRules(w http.ResponseWriter, r *http.Request) {
	flagID, envID, _, _, _, err := h.resolve(w, r)
	if err != nil {
		return
	}

	rules, err := h.svc.GetRules(r.Context(), flagID, envID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "fetch_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": rules})
}

func (h *targetingHandlerWrapper) SetRules(w http.ResponseWriter, r *http.Request) {
	flagID, envID, envSlug, projectSlug, projectID, err := h.resolve(w, r)
	if err != nil {
		return
	}

	var req targeting.SetRulesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	rules, err := h.svc.SetRules(r.Context(), flagID, envID, req.Rules)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	flagKey := chi.URLParam(r, "key")

	claims := auth.GetClaims(r.Context())
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID: projectID, UserID: claims.UserID, Action: "updated",
		EntityType: "targeting_rules", EntityID: flagID,
		Changes: map[string]any{"flag": flagKey, "environment": envSlug, "rule_count": len(rules)},
		IPAddress: r.RemoteAddr,
	})

	h.cache.Invalidate(r.Context(), "eval:"+projectSlug+":"+envSlug+":"+flagKey)

	h.hub.Broadcast(envSlug, sse.Event{
		Type: "flag.updated",
		Data: map[string]string{
			"flag_key":    flagKey,
			"environment": envSlug,
		},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": rules})
}

func (h *targetingHandlerWrapper) resolve(w http.ResponseWriter, r *http.Request) (flagID, envID, envSlug, projectSlug, projectID string, err error) {
	projectSlug = chi.URLParam(r, "slug")
	flagKey := chi.URLParam(r, "key")
	envSlug = chi.URLParam(r, "env")

	p, pErr := h.projectSvc.GetBySlug(r.Context(), projectSlug)
	if pErr != nil {
		writeError(w, http.StatusNotFound, "not_found", "project not found")
		return "", "", "", "", "", fmt.Errorf("project not found")
	}

	f, fErr := h.flagSvc.GetByKey(r.Context(), p.ID, flagKey)
	if fErr != nil {
		writeError(w, http.StatusNotFound, "not_found", "flag not found")
		return "", "", "", "", "", fmt.Errorf("flag not found")
	}

	e, eErr := h.envSvc.GetBySlug(r.Context(), p.ID, envSlug)
	if eErr != nil {
		writeError(w, http.StatusNotFound, "not_found", "environment not found")
		return "", "", "", "", "", fmt.Errorf("environment not found")
	}

	return f.ID, e.ID, envSlug, projectSlug, p.ID, nil
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}
