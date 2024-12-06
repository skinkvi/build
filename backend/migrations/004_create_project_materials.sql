-- +migrate Up

CREATE TABLE IF NOT EXISTS project_materials (
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    material_id INTEGER NOT NULL REFERENCES materials(id),
    PRIMARY KEY (project_id, material_id)
);

-- +migrate Down

DROP TABLE IF EXISTS project_materials CASCADE; 