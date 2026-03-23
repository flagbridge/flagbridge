CREATE TABLE IF NOT EXISTS flag_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id UUID NOT NULL REFERENCES flags(id) ON DELETE CASCADE,
    environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT false,
    default_value JSONB,
    rules JSONB DEFAULT '[]',
    updated_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (flag_id, environment_id)
);

CREATE INDEX idx_flag_states_lookup ON flag_states (flag_id, environment_id);
