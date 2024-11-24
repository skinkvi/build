-- +migrate Up

CREATE TABLE IF NOT EXISTS ai_responses (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    response_data TEXT,
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +migrate Down

DROP TABLE IF EXISTS ai_responses CASCADE; 