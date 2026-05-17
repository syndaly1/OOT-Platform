CREATE TABLE IF NOT EXISTS watch_progress (
    id SERIAL PRIMARY KEY,

    profile_id TEXT NOT NULL,
    content_id TEXT NOT NULL,
    episode_id TEXT,

    current_position_sec INT DEFAULT 0,
    duration_sec INT DEFAULT 0,

    watched_percent REAL DEFAULT 0,

    completed BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);