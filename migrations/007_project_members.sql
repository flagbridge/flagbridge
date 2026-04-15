-- +goose Up
-- Project members: role-based access control per project.
-- Roles: admin, engineer, product, viewer.
-- The "product" role is the FlagBridge differentiator — PMs as first-class users.

CREATE TABLE IF NOT EXISTS project_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'viewer',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, user_id)
);

CREATE INDEX idx_project_members_project ON project_members (project_id);
CREATE INDEX idx_project_members_user ON project_members (user_id);

-- Auto-add existing project creators as admin members.
INSERT INTO project_members (project_id, user_id, role)
SELECT id, created_by, 'admin'
FROM projects
ON CONFLICT (project_id, user_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS project_members;
