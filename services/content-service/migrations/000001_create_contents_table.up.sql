CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    poster_url TEXT,
    type TEXT NOT NULL,
    genre TEXT NOT NULL,
    release_year INT NOT NULL,
    rating DOUBLE PRECISION DEFAULT 0,
    is_featured BOOLEAN DEFAULT false,
    is_trending BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
