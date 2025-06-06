-- Migration pour la création des tables mythicplus de la feature Statistiques Mythic+ Raider.io 
-- 27/05/2025


-- Création de la table team_compositions en premier (référencée par les autres)
-- Cette table contient les compositions de donjons mythiques
CREATE TABLE mythicplus_team_compositions (
    id BIGSERIAL PRIMARY KEY,
    composition_hash VARCHAR(64) UNIQUE NOT NULL,
    
    -- Tank
    tank_class VARCHAR(50) NOT NULL,
    tank_spec VARCHAR(50) NOT NULL,
    
    -- Healer
    healer_class VARCHAR(50) NOT NULL,
    healer_spec VARCHAR(50) NOT NULL,
    
    -- DPS (ordonnés alphabétiquement)
    dps1_class VARCHAR(50) NOT NULL,
    dps1_spec VARCHAR(50) NOT NULL,
    dps2_class VARCHAR(50) NOT NULL,
    dps2_spec VARCHAR(50) NOT NULL,
    dps3_class VARCHAR(50) NOT NULL,
    dps3_spec VARCHAR(50) NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Création de la table mythicplus_runs
-- Cette table contient les runs mythiques
CREATE TABLE mythicplus_runs (
    id BIGSERIAL PRIMARY KEY,
    keystone_run_id BIGINT UNIQUE NOT NULL,
    season VARCHAR(50) NOT NULL,
    region VARCHAR(10) NOT NULL,
    dungeon_slug VARCHAR(100) NOT NULL,
    dungeon_name VARCHAR(200),
    mythic_level INTEGER NOT NULL,
    score DECIMAL(10,2),
    status VARCHAR(20),
    clear_time_ms BIGINT,
    keystone_time_ms BIGINT,
    completed_at TIMESTAMP WITH TIME ZONE,
    num_chests INTEGER,
    time_remaining_ms BIGINT,
    team_composition_id BIGINT REFERENCES mythicplus_team_compositions(id) ON DELETE RESTRICT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Création de la table run_roster (vue dénormalisée pour analyses)
-- Cette table contient les rosters des runs mythiques
CREATE TABLE mythicplus_run_roster (
    id BIGSERIAL PRIMARY KEY,
    team_composition_id BIGINT NOT NULL REFERENCES mythicplus_team_compositions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL, -- tank, healer, dps
    class_name VARCHAR(50) NOT NULL,
    spec_name VARCHAR(50) NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Index de performance pour mythicplus_runs
CREATE INDEX idx_mythicplus_runs_season_region ON mythicplus_runs(season, region);
CREATE INDEX idx_mythicplus_runs_dungeon_level ON mythicplus_runs(dungeon_slug, mythic_level);
CREATE INDEX idx_mythicplus_runs_score ON mythicplus_runs(score DESC);
CREATE INDEX idx_mythicplus_runs_completed_at ON mythicplus_runs(completed_at DESC);
CREATE INDEX idx_mythicplus_runs_status ON mythicplus_runs(status);
CREATE INDEX idx_mythicplus_runs_team_composition_id ON mythicplus_runs(team_composition_id);
CREATE INDEX idx_mythicplus_runs_deleted_at ON mythicplus_runs(deleted_at);

-- Index de performance pour team_compositions
CREATE INDEX idx_mythicplus_team_compositions_deleted_at ON mythicplus_team_compositions(deleted_at);

-- Index de performance pour run_roster (analyses rapides)
CREATE INDEX idx_mythicplus_run_roster_team_composition_id ON mythicplus_run_roster(team_composition_id);
CREATE INDEX idx_mythicplus_run_roster_role ON mythicplus_run_roster(role);
CREATE INDEX idx_mythicplus_run_roster_class_spec ON mythicplus_run_roster(class_name, spec_name);
CREATE INDEX idx_mythicplus_run_roster_deleted_at ON mythicplus_run_roster(deleted_at);

-- Index composite pour analyses avancées
CREATE INDEX idx_mythicplus_runs_analytics ON mythicplus_runs(season, region, dungeon_slug, score DESC, completed_at DESC);