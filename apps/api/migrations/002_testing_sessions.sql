-- Testing Sessions for FlagBridge Testing API
-- Isolated sessions with scoped flag overrides for E2E testing.

CREATE TABLE IF NOT EXISTS testing_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    label TEXT DEFAULT '',
    overrides JSONB NOT NULL DEFAULT '{}',
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_testing_sessions_project ON testing_sessions (project_id, expires_at);
CREATE INDEX idx_testing_sessions_expires ON testing_sessions (expires_at);
