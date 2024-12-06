-- 008_alter_warcraft_logs_reports.up.sql
ALTER TABLE warcraft_logs_reports 
  -- base data
  ADD COLUMN IF NOT EXISTS total_time BIGINT,
  ADD COLUMN IF NOT EXISTS item_level NUMERIC,
  
  -- composition and performances
  ADD COLUMN IF NOT EXISTS composition JSONB,              -- composition array
  ADD COLUMN IF NOT EXISTS damage_done JSONB,              -- damageDone array
  ADD COLUMN IF NOT EXISTS healing_done JSONB,             -- healingDone array
  ADD COLUMN IF NOT EXISTS damage_taken JSONB,             -- damageTaken array
  ADD COLUMN IF NOT EXISTS death_events JSONB,             -- deathEvents array
  
  -- player details
  ADD COLUMN IF NOT EXISTS player_details_dps JSONB,       -- playerDetails.dps array
  ADD COLUMN IF NOT EXISTS player_details_healers JSONB,   -- playerDetails.healers array
  ADD COLUMN IF NOT EXISTS player_details_tanks JSONB,     -- playerDetails.tanks array
  
  -- combat data
  ADD COLUMN IF NOT EXISTS log_version INTEGER,
  ADD COLUMN IF NOT EXISTS game_version INTEGER,
  
  -- mythic key data
  ADD COLUMN IF NOT EXISTS friendly_players INTEGER[],     -- fights.friendlyPlayers array
  ADD COLUMN IF NOT EXISTS talent_codes JSONB,             -- Map of talentImportCode per actorID
  
  -- raw_data modification
  ALTER COLUMN raw_data DROP NOT NULL;

-- Indexes to optimize searches
CREATE INDEX IF NOT EXISTS idx_reports_composition ON warcraft_logs_reports USING gin (composition jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_damage_done ON warcraft_logs_reports USING gin (damage_done jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_healing_done ON warcraft_logs_reports USING gin (healing_done jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_damage_taken ON warcraft_logs_reports USING gin (damage_taken jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_death_events ON warcraft_logs_reports USING gin (death_events jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_player_details_dps ON warcraft_logs_reports USING gin (player_details_dps jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_player_details_healers ON warcraft_logs_reports USING gin (player_details_healers jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_player_details_tanks ON warcraft_logs_reports USING gin (player_details_tanks jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_reports_talent_codes ON warcraft_logs_reports USING gin (talent_codes jsonb_path_ops);