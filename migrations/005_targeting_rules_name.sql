-- +goose Up
-- Add missing 'name' column to targeting_rules
-- The table was created in Supabase without this column

ALTER TABLE targeting_rules ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT '';
