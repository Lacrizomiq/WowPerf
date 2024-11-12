CREATE TABLE IF NOT EXISTS rankings_update_states (
    id INTEGER PRIMARY KEY DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    last_update_time TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP - INTERVAL '25 hours'),
    CONSTRAINT ensure_single_row CHECK (id = 1)
);

CREATE TABLE IF NOT EXISTS player_rankings (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    dungeon_id INTEGER,
    name VARCHAR(255),
    class VARCHAR(255),
    spec VARCHAR(255),
    role VARCHAR(255),
    amount DOUBLE PRECISION,
    hard_mode_level INTEGER,
    duration BIGINT,
    start_time BIGINT,
    report_code VARCHAR(255),
    report_fight_id INTEGER,
    report_start_time BIGINT,
    guild_id INTEGER,
    guild_name VARCHAR(255),
    guild_faction INTEGER,
    server_id INTEGER,
    server_name VARCHAR(255),
    server_region VARCHAR(50),
    bracket_data INTEGER,
    faction INTEGER,
    affixes INTEGER[],
    medal VARCHAR(50),
    score DOUBLE PRECISION,
    leaderboard INTEGER DEFAULT 0
);

-- Indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_player_rankings_dungeon_id ON player_rankings(dungeon_id);
CREATE INDEX IF NOT EXISTS idx_player_rankings_name ON player_rankings(name);
CREATE INDEX IF NOT EXISTS idx_player_rankings_class ON player_rankings(class);
CREATE INDEX IF NOT EXISTS idx_player_rankings_spec ON player_rankings(spec);
CREATE INDEX IF NOT EXISTS idx_player_rankings_role ON player_rankings(role);
CREATE INDEX IF NOT EXISTS idx_player_rankings_guild_id ON player_rankings(guild_id);
CREATE INDEX IF NOT EXISTS idx_player_rankings_server_id ON player_rankings(server_id);
CREATE INDEX IF NOT EXISTS idx_player_rankings_deleted_at ON player_rankings(deleted_at);
CREATE INDEX IF NOT EXISTS idx_player_rankings_score ON player_rankings(score);
CREATE INDEX IF NOT EXISTS idx_rankings_update_states_deleted_at ON rankings_update_states(deleted_at);