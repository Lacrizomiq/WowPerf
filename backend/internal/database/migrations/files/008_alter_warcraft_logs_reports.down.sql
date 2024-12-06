-- 008_alter_warcraft_logs_reports.down.sql
ALTER TABLE warcraft_logs_reports 
  DROP COLUMN IF EXISTS total_time,
  DROP COLUMN IF EXISTS item_level,
  DROP COLUMN IF EXISTS composition,
  DROP COLUMN IF EXISTS damage_done,
  DROP COLUMN IF EXISTS healing_done,
  DROP COLUMN IF EXISTS damage_taken,
  DROP COLUMN IF EXISTS death_events,
  DROP COLUMN IF EXISTS player_details_dps,
  DROP COLUMN IF EXISTS player_details_healers,
  DROP COLUMN IF EXISTS player_details_tanks,
  DROP COLUMN IF EXISTS log_version,
  DROP COLUMN IF EXISTS game_version,
  DROP COLUMN IF EXISTS friendly_players,
  DROP COLUMN IF EXISTS talent_codes,
  ALTER COLUMN raw_data SET NOT NULL;

DROP INDEX IF EXISTS idx_reports_composition;
DROP INDEX IF EXISTS idx_reports_damage_done;
DROP INDEX IF EXISTS idx_reports_healing_done;
DROP INDEX IF EXISTS idx_reports_damage_taken;
DROP INDEX IF EXISTS idx_reports_death_events;
DROP INDEX IF EXISTS idx_reports_player_details_dps;
DROP INDEX IF EXISTS idx_reports_player_details_healers;
DROP INDEX IF EXISTS idx_reports_player_details_tanks;
DROP INDEX IF EXISTS idx_reports_talent_codes;