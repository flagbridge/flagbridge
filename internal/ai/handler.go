package ai

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flagbridge/flagbridge/internal/aiconfig"
	"github.com/flagbridge/flagbridge/internal/audit"
	"github.com/flagbridge/flagbridge/internal/auth"
	"github.com/flagbridge/flagbridge/internal/project"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	aiConfigSvc *aiconfig.Service
	projectSvc  *project.Service
	auditSvc    *audit.Service
	ctxBuilder  *ContextBuilder
}

func NewHandler(aiConfigSvc *aiconfig.Service, projectSvc *project.Service, auditSvc *audit.Service, ctxBuilder *ContextBuilder) *Handler {
	return &Handler{
		aiConfigSvc: aiConfigSvc,
		projectSvc:  projectSvc,
		auditSvc:    auditSvc,
		ctxBuilder:  ctxBuilder,
	}
}

// Chat handles POST /v1/projects/{slug}/ai/chat — streaming SSE response.
func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.resolveProject(w, r)
	if err != nil {
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}
	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "missing_message", "message is required")
		return
	}

	// Rate limit check
	if err := h.aiConfigSvc.CheckAndIncrementUsage(r.Context(), projectID); err != nil {
		writeError(w, http.StatusTooManyRequests, "rate_limited", err.Error())
		return
	}

	// Load AI config
	cfg, err := h.aiConfigSvc.Get(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_configured", "AI provider not configured for this project")
		return
	}

	// Decrypt API key (not needed for Ollama)
	var apiKey string
	if cfg.ProviderName != "ollama" {
		apiKey, err = h.aiConfigSvc.DecryptAPIKey(r.Context(), projectID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "no_api_key", "AI provider API key not configured")
			return
		}
	}

	// Build context
	projectContext, err := h.ctxBuilder.Build(r.Context(), projectID)
	if err != nil {
		slog.Error("failed to build AI context", "error", err, "project_id", projectID)
		writeError(w, http.StatusInternalServerError, "context_error", "failed to build project context")
		return
	}

	// Create provider
	var baseURL string
	if cfg.BaseURL != nil {
		baseURL = *cfg.BaseURL
	}
	provider, err := NewProvider(cfg.ProviderName, apiKey, baseURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "provider_error", err.Error())
		return
	}

	// Build request
	provReq := ProviderRequest{
		Model:       cfg.ModelID,
		System:      SystemPrompt(projectContext),
		Messages:    []Message{{Role: "user", Content: req.Message}},
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		Stream:      true,
	}

	// Stream response via SSE
	stream, err := provider.ChatStream(provReq)
	if err != nil {
		slog.Error("AI stream failed", "error", err, "provider", cfg.ProviderName)
		writeError(w, http.StatusBadGateway, "provider_error", "failed to connect to AI provider")
		return
	}

	// Audit log
	claims := auth.GetClaims(r.Context())
	slug := chi.URLParam(r, "slug")
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID: projectID, UserID: claims.UserID, Action: "ai_chat",
		EntityType: "ai_provider", EntityID: cfg.ID,
		Changes:   map[string]any{"provider": cfg.ProviderName, "model": cfg.ModelID, "project": slug},
		IPAddress: r.RemoteAddr,
	})

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "stream_error", "streaming not supported")
		return
	}

	for chunk := range stream {
		if chunk.Err != nil {
			writeSSE(w, flusher, StreamEvent{Type: "error", Error: chunk.Err.Error()})
			return
		}
		if chunk.Done {
			writeSSE(w, flusher, StreamEvent{Type: "done", Usage: chunk.Usage})
			return
		}
		if chunk.Content != "" {
			writeSSE(w, flusher, StreamEvent{Type: "content", Content: chunk.Content})
		}
	}
}

// AnalyzeFlag handles POST /v1/projects/{slug}/ai/analyze-flag — non-streaming.
func (h *Handler) AnalyzeFlag(w http.ResponseWriter, r *http.Request) {
	projectID, err := h.resolveProject(w, r)
	if err != nil {
		return
	}

	flagKey := chi.URLParam(r, "key")
	if flagKey == "" {
		writeError(w, http.StatusBadRequest, "missing_flag_key", "flag key is required")
		return
	}

	if err := h.aiConfigSvc.CheckAndIncrementUsage(r.Context(), projectID); err != nil {
		writeError(w, http.StatusTooManyRequests, "rate_limited", err.Error())
		return
	}

	cfg, err := h.aiConfigSvc.Get(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_configured", "AI provider not configured for this project")
		return
	}

	var apiKey string
	if cfg.ProviderName != "ollama" {
		apiKey, err = h.aiConfigSvc.DecryptAPIKey(r.Context(), projectID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "no_api_key", "AI provider API key not configured")
			return
		}
	}

	projectContext, err := h.ctxBuilder.Build(r.Context(), projectID)
	if err != nil {
		slog.Error("failed to build AI context", "error", err, "project_id", projectID)
		writeError(w, http.StatusInternalServerError, "context_error", "failed to build project context")
		return
	}

	var baseURL string
	if cfg.BaseURL != nil {
		baseURL = *cfg.BaseURL
	}
	provider, err := NewProvider(cfg.ProviderName, apiKey, baseURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "provider_error", err.Error())
		return
	}

	provReq := ProviderRequest{
		Model:       cfg.ModelID,
		System:      SystemPrompt(projectContext),
		Messages:    []Message{{Role: "user", Content: AnalyzeFlagPrompt(flagKey)}},
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
	}

	resp, err := provider.Chat(provReq)
	if err != nil {
		slog.Error("AI chat failed", "error", err, "provider", cfg.ProviderName)
		writeError(w, http.StatusBadGateway, "provider_error", "failed to get AI response")
		return
	}

	claims := auth.GetClaims(r.Context())
	slug := chi.URLParam(r, "slug")
	h.auditSvc.Log(r.Context(), audit.LogInput{
		ProjectID: projectID, UserID: claims.UserID, Action: "ai_analyze_flag",
		EntityType: "flag", EntityID: flagKey,
		Changes:   map[string]any{"provider": cfg.ProviderName, "model": cfg.ModelID, "project": slug},
		IPAddress: r.RemoteAddr,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"data": ChatResponse{
			Response: resp.Content,
			Model:    cfg.ModelID,
			Provider: cfg.ProviderName,
			Usage:    resp.Usage,
		},
	})
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

func writeSSE(w http.ResponseWriter, flusher http.Flusher, event StreamEvent) {
	data, _ := json.Marshal(event)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
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
