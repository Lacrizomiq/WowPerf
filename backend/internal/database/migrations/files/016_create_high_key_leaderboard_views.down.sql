-- Drop all views created in the up migration, in reverse order
DROP VIEW IF EXISTS top_5_players_per_role;
DROP VIEW IF EXISTS top_10_players_per_spec;
DROP VIEW IF EXISTS dungeon_avg_key_levels;
DROP VIEW IF EXISTS class_global_score_averages;
DROP VIEW IF EXISTS spec_dungeon_max_key_levels;
DROP VIEW IF EXISTS spec_global_score_averages;