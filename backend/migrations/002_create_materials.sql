-- +migrate Up

CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    article VARCHAR(255) NOT NULL,
    image_url TEXT,
    description TEXT,
    purchase_location TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +migrate Down

