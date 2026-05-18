CREATE TABLE IF NOT EXISTS episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id UUID REFERENCES contents(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    season_number INT NOT NULL,
    episode_number INT NOT NULL,
    duration_seconds INT NOT NULL,
    skip_intro_start INT DEFAULT 0,
    skip_intro_end INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
