-- +migrate Up

CREATE INDEX IF NOT EXISTS idx_materials_name_trgm ON materials USING gin (name gin_trgm_ops);
