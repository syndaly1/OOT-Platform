CREATE TYPE profile_type AS ENUM ('adult', 'teen', 'kid');

CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type profile_type DEFAULT 'adult',
    language VARCHAR(10) DEFAULT 'en',
    avatar VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);