-- Migration: 024_fix_missing_views.down.sql
-- Purpose: Drop views created in the up migration
-- Date: April 04, 2025

-- Since these views were meant to exist from previous migrations,
-- dropping them should generally not be necessary.
-- However, for completeness, we provide a way to revert if needed:

DROP VIEW IF EXISTS spec_dungeon_max_key_levels;
DROP VIEW IF EXISTS class_global_score_averages;
DROP VIEW IF EXISTS dungeon_avg_key_levels;
DROP VIEW IF EXISTS top_5_players_per_role;