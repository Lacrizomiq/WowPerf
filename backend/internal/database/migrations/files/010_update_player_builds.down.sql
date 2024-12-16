-- 010_update_player_builds.down.sql

-- Restore the name of the column
ALTER TABLE player_builds 
    RENAME COLUMN talent_code TO talent_import;

-- Restore the talent_tree_id column
ALTER TABLE player_builds 
    ADD COLUMN IF NOT EXISTS talent_tree_id INTEGER;

-- Delete the item_level column
ALTER TABLE player_builds 
    DROP COLUMN IF EXISTS item_level;

-- Delete the keystone_level and affixes columns
ALTER TABLE player_builds 
    DROP COLUMN IF EXISTS keystone_level,
    DROP COLUMN IF EXISTS affixes;

-- Delete the index
DROP INDEX IF EXISTS idx_player_builds_item_level;
DROP INDEX IF EXISTS idx_player_builds_keystone_level;
DROP INDEX IF EXISTS idx_player_builds_affixes;

DROP INDEX IF EXISTS idx_player_builds_gear;
DROP INDEX IF EXISTS idx_player_builds_stats;
DROP INDEX IF EXISTS idx_player_builds_combatant_info;
DROP INDEX IF EXISTS idx_player_builds_created_at;
DROP INDEX IF EXISTS idx_player_builds_class_spec_keystone;