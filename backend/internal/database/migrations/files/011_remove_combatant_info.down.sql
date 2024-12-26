-- 011_cleanup_player_builds.down.sql

-- Restore combatant_info column
ALTER TABLE player_builds
    ADD COLUMN IF NOT EXISTS combatant_info JSONB;

-- Recreate the index
CREATE INDEX IF NOT EXISTS idx_player_builds_combatant_info ON player_builds USING gin(combatant_info jsonb_path_ops);