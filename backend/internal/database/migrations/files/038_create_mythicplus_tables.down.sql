-- Suppression des index dans l'ordre inverse
DROP INDEX IF EXISTS idx_mythicplus_runs_analytics;
DROP INDEX IF EXISTS idx_mythicplus_run_roster_deleted_at;
DROP INDEX IF EXISTS idx_mythicplus_run_roster_class_spec;
DROP INDEX IF EXISTS idx_mythicplus_run_roster_role;
DROP INDEX IF EXISTS idx_mythicplus_run_roster_team_composition_id;
DROP INDEX IF EXISTS idx_mythicplus_team_compositions_deleted_at;
DROP INDEX IF EXISTS idx_mythicplus_runs_deleted_at;
DROP INDEX IF EXISTS idx_mythicplus_runs_team_composition_id;
DROP INDEX IF EXISTS idx_mythicplus_runs_status;
DROP INDEX IF EXISTS idx_mythicplus_runs_completed_at;
DROP INDEX IF EXISTS idx_mythicplus_runs_score;
DROP INDEX IF EXISTS idx_mythicplus_runs_dungeon_level;
DROP INDEX IF EXISTS idx_mythicplus_runs_season_region;

-- Suppression des tables dans l'ordre inverse pour respecter les contraintes FK
DROP TABLE IF EXISTS mythicplus_run_roster;
DROP TABLE IF EXISTS mythicplus_runs;
DROP TABLE IF EXISTS mythicplus_team_compositions;