-- 010_update_player_builds.up.sql

-- Delete of the talent_tree_id column
ALTER TABLE player_builds 
    DROP COLUMN IF EXISTS talent_tree_id;

-- Rename of the talent_import column to talent_code
ALTER TABLE player_builds 
    RENAME COLUMN talent_import TO talent_code;

-- Add the item_level column
ALTER TABLE player_builds 
    ADD COLUMN IF NOT EXISTS item_level NUMERIC;

-- Add the keystone_level and affixes columns
ALTER TABLE player_builds 
    ADD COLUMN IF NOT EXISTS keystone_level INTEGER,
    ADD COLUMN IF NOT EXISTS affixes INTEGER[];

-- Add additional indexes to improve query performance
CREATE INDEX IF NOT EXISTS idx_player_builds_item_level ON player_builds(item_level);
CREATE INDEX IF NOT EXISTS idx_player_builds_keystone_level ON player_builds(keystone_level);
CREATE INDEX IF NOT EXISTS idx_player_builds_affixes ON player_builds USING gin(affixes);

-- Index for JSON data
CREATE INDEX IF NOT EXISTS idx_player_builds_gear ON player_builds USING gin(gear jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_player_builds_stats ON player_builds USING gin(stats jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_player_builds_combatant_info ON player_builds USING gin(combatant_info jsonb_path_ops);

-- Index for temporal data
CREATE INDEX IF NOT EXISTS idx_player_builds_created_at ON player_builds(created_at);

-- Index for frequent searches
CREATE INDEX IF NOT EXISTS idx_player_builds_class_spec_keystone ON player_builds(class, spec, keystone_level);
