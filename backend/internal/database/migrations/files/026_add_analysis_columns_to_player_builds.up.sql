-- up migration

-- This migration is used to add columns to the player builds table to track the status of the analysis
-- It is used to decouple the analysis workflows from the player builds table

-- add columns to track the status of the analysis on the player builds table
ALTER TABLE player_builds
ADD COLUMN equipment_processed_at TIMESTAMP NULL,
ADD COLUMN talent_processed_at TIMESTAMP NULL,
ADD COLUMN stat_processed_at TIMESTAMP NULL,
ADD COLUMN equipment_status VARCHAR(255) NULL,
ADD COLUMN talent_status VARCHAR(255) NULL,
ADD COLUMN stat_status VARCHAR(255) NULL;

-- indexes for performance 
CREATE INDEX idx_player_builds_equipment_status ON player_builds(equipment_status);
CREATE INDEX idx_player_builds_talent_status ON player_builds(talent_status);
CREATE INDEX idx_player_builds_stat_status ON player_builds(stat_status);
CREATE INDEX idx_player_builds_equipment_processed_at ON player_builds(equipment_processed_at);
CREATE INDEX idx_player_builds_talent_processed_at ON player_builds(talent_processed_at);
CREATE INDEX idx_player_builds_stat_processed_at ON player_builds(stat_processed_at);