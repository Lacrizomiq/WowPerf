-- Migration: 19_create_stat_statistics.down.sql
DROP TABLE IF EXISTS stat_statistics;

DROP INDEX IF EXISTS idx_stat_statistics_class_spec;
DROP INDEX IF EXISTS idx_stat_statistics_encounter_id;
DROP INDEX IF EXISTS idx_stat_statistics_category;
DROP INDEX IF EXISTS idx_stat_statistics_deleted_at;

