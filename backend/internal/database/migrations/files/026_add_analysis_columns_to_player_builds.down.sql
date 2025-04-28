-- down migration

-- remove columns from the player builds table
ALTER TABLE player_builds
DROP COLUMN equipment_processed_at,
DROP COLUMN talent_processed_at,
DROP COLUMN stat_processed_at,
DROP COLUMN equipment_status,
DROP COLUMN talent_status,
DROP COLUMN stat_status;

-- Drop indexes 
DROP INDEX IF EXISTS idx_player_builds_equipment_status;
DROP INDEX IF EXISTS idx_player_builds_talent_status;
DROP INDEX IF EXISTS idx_player_builds_stat_status;
DROP INDEX IF EXISTS idx_player_builds_equipment_processed_at;
DROP INDEX IF EXISTS idx_player_builds_talent_processed_at;
DROP INDEX IF EXISTS idx_player_builds_stat_processed_at;