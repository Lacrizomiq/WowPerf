-- 001_initial_tables.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    battle_net_id INTEGER UNIQUE,
    battle_tag VARCHAR(255) UNIQUE,
    encrypted_token BYTEA,
    battle_net_expires_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_username_change_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    battle_net_refresh_token TEXT,
    battle_net_token_type VARCHAR(50),
    battle_net_scopes TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_email_key UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS dungeons (
    id SERIAL PRIMARY KEY,
    challenge_mode_id INTEGER UNIQUE,
    slug VARCHAR(255),
    name VARCHAR(255),
    short_name VARCHAR(255),
    media_url VARCHAR(255),
    icon VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS seasons (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255),
    name VARCHAR(255),
    short_name VARCHAR(255),
    seasonal_affix_id INTEGER,
    starts_us TIMESTAMP WITH TIME ZONE,
    starts_eu TIMESTAMP WITH TIME ZONE,
    starts_tw TIMESTAMP WITH TIME ZONE,
    starts_kr TIMESTAMP WITH TIME ZONE,
    starts_cn TIMESTAMP WITH TIME ZONE,
    ends_us TIMESTAMP WITH TIME ZONE,
    ends_eu TIMESTAMP WITH TIME ZONE,
    ends_tw TIMESTAMP WITH TIME ZONE,
    ends_kr TIMESTAMP WITH TIME ZONE,
    ends_cn TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS affixes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    icon VARCHAR(255),
    description TEXT,
    wowhead_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS key_stone_upgrades (
    id SERIAL PRIMARY KEY,
    challenge_mode_id INTEGER,
    qualifying_duration BIGINT,
    upgrade_level INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS raids (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) UNIQUE,
    name VARCHAR(255),
    short_name VARCHAR(255),
    expansion VARCHAR(255),
    media_url VARCHAR(255),
    icon VARCHAR(255),
    starts JSONB DEFAULT '{}'::jsonb,
    ends JSONB DEFAULT '{}'::jsonb,
    encounters JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS season_dungeons (
    season_id INTEGER REFERENCES seasons(id),
    dungeon_id INTEGER REFERENCES dungeons(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (season_id, dungeon_id)
);

-- Create indexes for soft delete
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_dungeons_deleted_at ON dungeons(deleted_at);
CREATE INDEX IF NOT EXISTS idx_seasons_deleted_at ON seasons(deleted_at);
CREATE INDEX IF NOT EXISTS idx_affixes_deleted_at ON affixes(deleted_at);
CREATE INDEX IF NOT EXISTS idx_key_stone_upgrades_deleted_at ON key_stone_upgrades(deleted_at);
CREATE INDEX IF NOT EXISTS idx_raids_deleted_at ON raids(deleted_at);
CREATE INDEX IF NOT EXISTS idx_raids_starts ON raids USING gin (starts);
CREATE INDEX IF NOT EXISTS idx_raids_ends ON raids USING gin (ends);
CREATE INDEX IF NOT EXISTS idx_raids_encounters ON raids USING gin (encounters);
CREATE INDEX IF NOT EXISTS idx_season_dungeons_deleted_at ON season_dungeons(deleted_at);