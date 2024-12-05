-- Drop indexes
DROP INDEX IF EXISTS idx_group_compositions_encounter;
DROP INDEX IF EXISTS idx_group_compositions_specs;
DROP INDEX IF EXISTS idx_group_compositions_deleted_at;

DROP INDEX IF EXISTS idx_talent_statistics_encounter;
DROP INDEX IF EXISTS idx_talent_statistics_class_spec;
DROP INDEX IF EXISTS idx_talent_statistics_deleted_at;

DROP INDEX IF EXISTS idx_build_statistics_encounter;
DROP INDEX IF EXISTS idx_build_statistics_item;
DROP INDEX IF EXISTS idx_build_statistics_class_spec;
DROP INDEX IF EXISTS idx_build_statistics_deleted_at;

DROP INDEX IF EXISTS idx_player_builds_encounter;
DROP INDEX IF EXISTS idx_player_builds_actor;
DROP INDEX IF EXISTS idx_player_builds_report;
DROP INDEX IF EXISTS idx_player_builds_class_spec;
DROP INDEX IF EXISTS idx_player_builds_deleted_at;

DROP INDEX IF EXISTS idx_reports_keystone;
DROP INDEX IF EXISTS idx_reports_encounter;
DROP INDEX IF EXISTS idx_warcraft_logs_reports_deleted_at;

DROP INDEX IF EXISTS idx_class_rankings_server;
DROP INDEX IF EXISTS idx_class_rankings_medal;
DROP INDEX IF EXISTS idx_class_rankings_hard_mode_level;
DROP INDEX IF EXISTS idx_class_rankings_score;
DROP INDEX IF EXISTS idx_class_rankings_encounter;
DROP INDEX IF EXISTS idx_class_rankings_dungeon;
DROP INDEX IF EXISTS idx_class_rankings_class_spec;
DROP INDEX IF EXISTS idx_class_rankings_deleted_at;

DROP INDEX IF EXISTS idx_dungeons_encounter_id;

-- Drop tables in reverse order of their creation
DROP TABLE IF EXISTS group_compositions;
DROP TABLE IF EXISTS talent_statistics;
DROP TABLE IF EXISTS build_statistics;
DROP TABLE IF EXISTS player_builds;
DROP TABLE IF EXISTS warcraft_logs_reports;
DROP TABLE IF EXISTS class_rankings;

-- Drop encounter_id column from dungeons table
ALTER TABLE dungeons DROP COLUMN IF EXISTS encounter_id;
