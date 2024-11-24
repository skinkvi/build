-- +migrate Up

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_type VARCHAR(50) NOT NULL CHECK (project_type IN ('design', 'fix', 'identify')),
    request_data TEXT,
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +migrate Down

DROP TABLE IF EXISTS projects CASCADE; 