CREATE TABLE IF NOT EXISTS watch_history (
    id SERIAL PRIMARY KEY,

    profile_id TEXT NOT NULL,
    content_id TEXT NOT NULL,
    episode_id TEXT,

    total_watched_sec INT DEFAULT 0,

    completed BOOLEAN DEFAULT FALSE,

    watched_at TIMESTAMP DEFAULT NOW()
);