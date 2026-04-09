-- FlagBridge Schema Alignment
-- Aligns the manually-created Supabase schema with migration 001+002+003

-- ============================================================
-- PROJECTS: add created_by
-- ============================================================
ALTER TABLE projects ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id);
UPDATE projects SET created_by = (SELECT id FROM users LIMIT 1) WHERE created_by IS NULL;
ALTER TABLE projects ALTER COLUMN created_by SET NOT NULL;

-- ============================================================
-- FLAGS: add name, created_by; rename key if needed
-- ============================================================
ALTER TABLE flags ADD COLUMN IF NOT EXISTS name TEXT;
UPDATE flags SET name = key WHERE name IS NULL;
ALTER TABLE flags ALTER COLUMN name SET NOT NULL;

ALTER TABLE flags ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id);
UPDATE flags SET created_by = (SELECT id FROM users LIMIT 1) WHERE created_by IS NULL;
ALTER TABLE flags ALTER COLUMN created_by SET NOT NULL;

-- ============================================================
-- API_KEYS: add environment_id, created_by, expires_at; rename last_used
-- ============================================================
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS environment_id UUID REFERENCES environments(id) ON DELETE SET NULL;

ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id);
UPDATE api_keys SET created_by = (SELECT id FROM users LIMIT 1) WHERE created_by IS NULL;
ALTER TABLE api_keys ALTER COLUMN created_by SET NOT NULL;

ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;

-- Rename last_used → last_used_at if the old column exists
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'api_keys' AND column_name = 'last_used') THEN
    ALTER TABLE api_keys RENAME COLUMN last_used TO last_used_at;
  ELSIF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'api_keys' AND column_name = 'last_used_at') THEN
    ALTER TABLE api_keys ADD COLUMN last_used_at TIMESTAMPTZ;
  END IF;
END
$$;

-- ============================================================
-- AUDIT LOG (from migration 001)
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    changes JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_project ON audit_log (project_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_entity ON audit_log (entity_type, entity_id);

-- ============================================================
-- TESTING SESSIONS (from migration 002)
-- ============================================================
CREATE TABLE IF NOT EXISTS testing_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    name TEXT,
    created_by UUID REFERENCES users(id),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS testing_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES testing_sessions(id) ON DELETE CASCADE,
    flag_key TEXT NOT NULL,
    value JSONB NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (session_id, flag_key)
);

CREATE INDEX IF NOT EXISTS idx_testing_sessions_project ON testing_sessions (project_id);
CREATE INDEX IF NOT EXISTS idx_testing_overrides_session ON testing_overrides (session_id);

-- ============================================================
-- WEBHOOKS (from migration 003)
-- ============================================================
CREATE TABLE IF NOT EXISTS webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    events TEXT[] NOT NULL DEFAULT '{}',
    active BOOLEAN NOT NULL DEFAULT true,
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS webhook_delivery_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    response_status INT,
    response_body TEXT,
    attempt INT NOT NULL DEFAULT 1,
    success BOOLEAN NOT NULL DEFAULT false,
    delivered_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_webhooks_project ON webhooks (project_id);
CREATE INDEX IF NOT EXISTS idx_webhook_logs_webhook ON webhook_delivery_logs (webhook_id, delivered_at DESC);
