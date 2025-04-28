-- Migration 021: Remove unique constraints from statistics tables
ALTER TABLE build_statistics DROP CONSTRAINT IF EXISTS build_statistics_unique;
ALTER TABLE talent_statistics DROP CONSTRAINT IF EXISTS talent_statistics_unique;
ALTER TABLE stat_statistics DROP CONSTRAINT IF EXISTS stat_statistics_unique;