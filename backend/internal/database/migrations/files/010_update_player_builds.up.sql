-- 010_update_player_builds.up.sql

-- Delete of the talent_tree_id column (déjà fait apparemment)
ALTER TABLE player_builds 
    DROP COLUMN IF EXISTS talent_tree_id;

-- La colonne talent_code existe déjà, pas besoin de renommage

-- Add the item_level column
ALTER TABLE player_builds 
    ADD COLUMN IF NOT EXISTS item_level NUMERIC;

-- Add the keystone_level and affixes columns
ALTER TABLE player_builds 
    ADD COLUMN IF NOT EXISTS keystone_level INTEGER,
    ADD COLUMN IF NOT EXISTS affixes INTEGER[];

-- Drop the dungeon_id column
ALTER TABLE player_builds
    DROP COLUMN IF EXISTS dungeon_id;

-- Add additional indexes to improve query performance
CREATE INDEX IF NOT EXISTS idx_player_builds_item_level ON player_builds(item_level);
CREATE INDEX IF NOT EXISTS idx_player_builds_keystone_level ON player_builds(keystone_level);
CREATE INDEX IF NOT EXISTS idx_player_builds_affixes ON player_builds USING gin(affixes);

-- Index for JSON data
CREATE INDEX IF NOT EXISTS idx_player_builds_gear ON player_builds USING gin(gear jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_player_builds_stats ON player_builds USING gin(stats jsonb_path_ops);

-- Index for temporal data
CREATE INDEX IF NOT EXISTS idx_player_builds_created_at ON player_builds(created_at);

-- Index for frequent searches
CREATE INDEX IF NOT EXISTS idx_player_builds_class_spec_keystone ON player_builds(class, spec, keystone_level);

-- Drop the dungeon index
DROP INDEX IF EXISTS idx_player_builds_dungeon;