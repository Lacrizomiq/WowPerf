-- Add encounter_id to the existing dungeons table
ALTER TABLE dungeons 
ADD COLUMN IF NOT EXISTS encounter_id INTEGER UNIQUE;

-- Update encounter_id for existing dungeons
UPDATE dungeons SET encounter_id = 12660 WHERE slug = 'arakara-city-of-echoes';
UPDATE dungeons SET encounter_id = 12669 WHERE slug = 'city-of-threads';
UPDATE dungeons SET encounter_id = 60670 WHERE slug = 'grim-batol';
UPDATE dungeons SET encounter_id = 62290 WHERE slug = 'mists-of-tirna-scithe';
UPDATE dungeons SET encounter_id = 61822 WHERE slug = 'siege-of-boralus';
UPDATE dungeons SET encounter_id = 12662 WHERE slug = 'the-dawnbreaker';
UPDATE dungeons SET encounter_id = 62286 WHERE slug = 'the-necrotic-wake';
UPDATE dungeons SET encounter_id = 12652 WHERE slug = 'the-stonevault';

-- Class/spec specific rankings
CREATE TABLE IF NOT EXISTS class_rankings (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    player_name VARCHAR(255) NOT NULL,
    class VARCHAR(255) NOT NULL,
    spec VARCHAR(255) NOT NULL,
    dungeon_id INTEGER REFERENCES dungeons(id),
    encounter_id INTEGER REFERENCES dungeons(encounter_id),
    amount NUMERIC,
    hard_mode_level INTEGER,
    duration BIGINT,
    start_time BIGINT,
    rank_position INTEGER,
    score NUMERIC,
    medal VARCHAR(50),
    leaderboard INTEGER,
    
    server_id INTEGER,
    server_name VARCHAR(255),
    server_region VARCHAR(50),
    
    guild_id INTEGER,
    guild_name VARCHAR(255),
    guild_faction INTEGER,
    
    report_code VARCHAR(255),
    report_fight_id INTEGER,
    report_start_time BIGINT,
    
    faction INTEGER,
    affixes INTEGER[]
);

-- WarcraftLogs reports storage
CREATE TABLE IF NOT EXISTS warcraft_logs_reports (
    code VARCHAR(255) PRIMARY KEY,
    fight_id INTEGER NOT NULL,
    encounter_id INTEGER REFERENCES dungeons(encounter_id),
    raw_data JSONB NOT NULL,
    player_details JSONB,
    keystoneLevel INTEGER,
    keystoneTime BIGINT,
    affixes INTEGER[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Player builds
CREATE TABLE IF NOT EXISTS player_builds (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    player_name VARCHAR(255) NOT NULL,
    class VARCHAR(255) NOT NULL,
    spec VARCHAR(255) NOT NULL,
    report_code VARCHAR(255) REFERENCES warcraft_logs_reports(code),
    fight_id INTEGER NOT NULL,
    talent_import TEXT,
    talent_tree JSONB,
    talent_tree_id INTEGER,
    actor_id INTEGER,
    gear JSONB,
    stats JSONB,
    combatant_info JSONB,
    dungeon_id INTEGER REFERENCES dungeons(id),
    encounter_id INTEGER REFERENCES dungeons(encounter_id)
);

-- Build statistics
CREATE TABLE IF NOT EXISTS build_statistics (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    class VARCHAR(255) NOT NULL,
    spec VARCHAR(255) NOT NULL,
    dungeon_id INTEGER REFERENCES dungeons(id),
    encounter_id INTEGER REFERENCES dungeons(encounter_id),
    item_slot INTEGER,
    item_id INTEGER,
    usage_count INTEGER DEFAULT 0,
    usage_percentage NUMERIC DEFAULT 0,
    avg_score NUMERIC DEFAULT 0,
    period_start TIMESTAMP WITH TIME ZONE,
    period_end TIMESTAMP WITH TIME ZONE
);

-- Talent statistics
CREATE TABLE IF NOT EXISTS talent_statistics (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    class VARCHAR(255) NOT NULL,
    spec VARCHAR(255) NOT NULL,
    dungeon_id INTEGER REFERENCES dungeons(id),
    encounter_id INTEGER REFERENCES dungeons(encounter_id),
    talent_import TEXT,
    usage_count INTEGER DEFAULT 0,
    usage_percentage NUMERIC DEFAULT 0,
    avg_score NUMERIC DEFAULT 0,
    period_start TIMESTAMP WITH TIME ZONE,
    period_end TIMESTAMP WITH TIME ZONE
);

-- Group compositions
CREATE TABLE IF NOT EXISTS group_compositions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    report_code VARCHAR(255) REFERENCES warcraft_logs_reports(code),
    fight_id INTEGER,
    tank_specs TEXT[],
    healer_specs TEXT[],
    dps_specs TEXT[],
    success BOOLEAN,
    dungeon_id INTEGER REFERENCES dungeons(id),
    encounter_id INTEGER REFERENCES dungeons(encounter_id),
    keystone_level INTEGER
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_dungeons_encounter_id ON dungeons(encounter_id);

CREATE INDEX IF NOT EXISTS idx_class_rankings_deleted_at ON class_rankings(deleted_at);
CREATE INDEX IF NOT EXISTS idx_class_rankings_class_spec ON class_rankings(class, spec);
CREATE INDEX IF NOT EXISTS idx_class_rankings_dungeon ON class_rankings(dungeon_id);
CREATE INDEX IF NOT EXISTS idx_class_rankings_encounter ON class_rankings(encounter_id);
CREATE INDEX IF NOT EXISTS idx_class_rankings_score ON class_rankings(score);
CREATE INDEX IF NOT EXISTS idx_class_rankings_hard_mode_level ON class_rankings(hard_mode_level);
CREATE INDEX IF NOT EXISTS idx_class_rankings_medal ON class_rankings(medal);
CREATE INDEX IF NOT EXISTS idx_class_rankings_server ON class_rankings(server_name, server_region);

CREATE INDEX IF NOT EXISTS idx_warcraft_logs_reports_deleted_at ON warcraft_logs_reports(deleted_at);
CREATE INDEX IF NOT EXISTS idx_reports_encounter ON warcraft_logs_reports(encounter_id);
CREATE INDEX IF NOT EXISTS idx_reports_keystone ON warcraft_logs_reports(keystoneLevel);

CREATE INDEX IF NOT EXISTS idx_player_builds_deleted_at ON player_builds(deleted_at);
CREATE INDEX IF NOT EXISTS idx_player_builds_class_spec ON player_builds(class, spec);
CREATE INDEX IF NOT EXISTS idx_player_builds_report ON player_builds(report_code, fight_id);
CREATE INDEX IF NOT EXISTS idx_player_builds_actor ON player_builds(actor_id);
CREATE INDEX IF NOT EXISTS idx_player_builds_encounter ON player_builds(encounter_id);

CREATE INDEX IF NOT EXISTS idx_build_statistics_deleted_at ON build_statistics(deleted_at);
CREATE INDEX IF NOT EXISTS idx_build_statistics_class_spec ON build_statistics(class, spec);
CREATE INDEX IF NOT EXISTS idx_build_statistics_item ON build_statistics(item_id);
CREATE INDEX IF NOT EXISTS idx_build_statistics_encounter ON build_statistics(encounter_id);

CREATE INDEX IF NOT EXISTS idx_talent_statistics_deleted_at ON talent_statistics(deleted_at);
CREATE INDEX IF NOT EXISTS idx_talent_statistics_class_spec ON talent_statistics(class, spec);
CREATE INDEX IF NOT EXISTS idx_talent_statistics_encounter ON talent_statistics(encounter_id);

CREATE INDEX IF NOT EXISTS idx_group_compositions_deleted_at ON group_compositions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_group_compositions_specs ON group_compositions USING gin(tank_specs, healer_specs, dps_specs);
CREATE INDEX IF NOT EXISTS idx_group_compositions_encounter ON group_compositions(encounter_id);