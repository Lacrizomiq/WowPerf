CREATE TABLE IF NOT EXISTS dungeon_stats (
    id SERIAL PRIMARY KEY,
    season VARCHAR(255),
    region VARCHAR(255),
    dungeon_slug VARCHAR(255),
    role_stats JSONB DEFAULT '{}'::jsonb,
    spec_stats JSONB DEFAULT '{}'::jsonb,
    level_stats JSONB DEFAULT '{}'::jsonb,
    team_comp JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uni_dungeon_stats_season_region_dungeon UNIQUE (season, region, dungeon_slug)
);

CREATE INDEX IF NOT EXISTS idx_dungeon_stats_deleted_at ON dungeon_stats(deleted_at);
CREATE INDEX IF NOT EXISTS idx_dungeon_stats_role_stats ON dungeon_stats USING gin (role_stats);
CREATE INDEX IF NOT EXISTS idx_dungeon_stats_spec_stats ON dungeon_stats USING gin (spec_stats);
CREATE INDEX IF NOT EXISTS idx_dungeon_stats_level_stats ON dungeon_stats USING gin (level_stats);
CREATE INDEX IF NOT EXISTS idx_dungeon_stats_team_comp ON dungeon_stats USING gin (team_comp);

CREATE TABLE IF NOT EXISTS update_states (
    id SERIAL PRIMARY KEY,
    last_update_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP - INTERVAL '25 hours',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_update_states_deleted_at ON update_states(deleted_at);