-- +goose Up
-- AI provider configuration: encrypted API keys for LLM providers.
-- Supports Anthropic, OpenAI, Ollama (self-hosted).
-- One provider config per project (can be updated, not multiple).

CREATE TABLE IF NOT EXISTS ai_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    provider TEXT NOT NULL DEFAULT 'anthropic',
    model_id TEXT NOT NULL DEFAULT 'claude-sonnet-4-20250514',
    encrypted_api_key TEXT,
    base_url TEXT,
    max_tokens INTEGER NOT NULL DEFAULT 4096,
    temperature NUMERIC(3,2) NOT NULL DEFAULT 0.7,
    monthly_usage INTEGER NOT NULL DEFAULT 0,
    monthly_limit INTEGER NOT NULL DEFAULT 100,
    last_reset_at TIMESTAMPTZ NOT NULL DEFAULT date_trunc('month', now()),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id)
);

CREATE INDEX idx_ai_providers_project ON ai_providers (project_id);

-- +goose Down
DROP TABLE IF EXISTS ai_providers;
