-- +migrate Down
-- Remove foreign keys
ALTER TABLE player_builds DROP CONSTRAINT IF EXISTS fk_player_builds_report;
ALTER TABLE player_builds DROP CONSTRAINT IF EXISTS fk_player_builds_dungeon;

-- Remove indexes
DROP INDEX IF EXISTS idx_player_builds_encounter_id;
DROP INDEX IF EXISTS idx_player_builds_report_code_fight_id;
DROP INDEX IF EXISTS idx_player_builds_player_info;
DROP INDEX IF EXISTS idx_player_builds_keystone_level;

-- Remove NOT NULL constraints
ALTER TABLE player_builds 
    ALTER COLUMN report_code DROP NOT NULL,
    ALTER COLUMN player_name DROP NOT NULL,
    ALTER COLUMN class DROP NOT NULL,
    ALTER COLUMN spec DROP NOT NULL,
    ALTER COLUMN fight_id DROP NOT NULL;

-- Rename encounter_id to dungeon_id
ALTER TABLE player_builds RENAME COLUMN encounter_id TO dungeon_id;