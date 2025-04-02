-- +migrate Down

-- 1. Remove all foreign key constraints
ALTER TABLE player_builds
    DROP CONSTRAINT IF EXISTS fk_player_builds_report,
    DROP CONSTRAINT IF EXISTS fk_player_builds_dungeon;

ALTER TABLE group_compositions
    DROP CONSTRAINT IF EXISTS fk_group_compositions_report;

-- 2. Remove unique constraints
ALTER TABLE player_builds
    DROP CONSTRAINT IF EXISTS unique_player_builds_report_fight_actor;

-- 3. Remove all indexes
DROP INDEX IF EXISTS idx_player_builds_encounter_id;
DROP INDEX IF EXISTS idx_player_builds_report_code_fight_id;
DROP INDEX IF EXISTS idx_player_builds_player_info;
DROP INDEX IF EXISTS idx_player_builds_keystone_level;
DROP INDEX IF EXISTS idx_player_builds_actor_id;
DROP INDEX IF EXISTS idx_warcraft_logs_reports_code;

-- 4. Restore the original state of the primary key on warcraft_logs_reports
ALTER TABLE warcraft_logs_reports
    DROP CONSTRAINT IF EXISTS warcraft_logs_reports_pkey;

-- Create a new primary key on the code column only (assuming this was the original state)
ALTER TABLE warcraft_logs_reports
    ADD PRIMARY KEY (code);

-- 5. Make the NOT NULL constraints nullable again
ALTER TABLE player_builds
    ALTER COLUMN report_code DROP NOT NULL,
    ALTER COLUMN player_name DROP NOT NULL,
    ALTER COLUMN class DROP NOT NULL,
    ALTER COLUMN spec DROP NOT NULL,
    ALTER COLUMN fight_id DROP NOT NULL;

-- Note: We don't remove the unique constraint on dungeons.encounter_id
-- as it might be needed by other parts of the application