-- +goose Up
-- Product Cards: hypothesis, success metrics, and lifecycle per flag.
-- Promoted from Pro to CE as core differentiator (v2-ia pivot).

CREATE TABLE IF NOT EXISTS product_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id UUID NOT NULL UNIQUE REFERENCES flags(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    hypothesis TEXT NOT NULL DEFAULT '',
    success_metrics TEXT NOT NULL DEFAULT '',
    go_no_go TEXT NOT NULL DEFAULT '',
    owner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'planning',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_product_cards_flag ON product_cards (flag_id);
CREATE INDEX idx_product_cards_project ON product_cards (project_id);

-- +goose Down
DROP TABLE IF EXISTS product_cards;
