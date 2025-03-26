-- 011_cleanup_player_builds.up.sql

-- Drop combatant_info index first
DROP INDEX IF EXISTS idx_player_builds_combatant_info;

-- Then drop combatant_info column
ALTER TABLE player_builds
    DROP COLUMN IF EXISTS combatant_info;