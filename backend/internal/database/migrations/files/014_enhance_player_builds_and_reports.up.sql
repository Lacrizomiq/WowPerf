-- +migrate Up

-- 1. First drop all dependent foreign keys
-- Ensure we drop all possible constraints that might exist
ALTER TABLE player_builds
    DROP CONSTRAINT IF EXISTS fk_player_builds_report,
    DROP CONSTRAINT IF EXISTS player_builds_report_code_fkey;

ALTER TABLE group_compositions
    DROP CONSTRAINT IF EXISTS group_compositions_report_code_fkey,
    DROP CONSTRAINT IF EXISTS fk_group_compositions_report;

-- Now it's safe to modify the primary key
ALTER TABLE warcraft_logs_reports 
    DROP CONSTRAINT IF EXISTS warcraft_logs_reports_pkey;

ALTER TABLE warcraft_logs_reports
    ADD PRIMARY KEY (code, fight_id);

CREATE INDEX IF NOT EXISTS idx_warcraft_logs_reports_code
    ON warcraft_logs_reports(code);

-- 2. Enhance player_builds table structure
ALTER TABLE player_builds 
    ALTER COLUMN report_code SET NOT NULL,
    ALTER COLUMN player_name SET NOT NULL,
    ALTER COLUMN class SET NOT NULL,
    ALTER COLUMN spec SET NOT NULL,
    ALTER COLUMN fight_id SET NOT NULL;

-- 3. Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_player_builds_encounter_id ON player_builds(encounter_id);
CREATE INDEX IF NOT EXISTS idx_player_builds_report_code_fight_id ON player_builds(report_code, fight_id);
CREATE INDEX IF NOT EXISTS idx_player_builds_player_info ON player_builds(player_name, class, spec);
CREATE INDEX IF NOT EXISTS idx_player_builds_keystone_level ON player_builds(keystone_level);
CREATE INDEX IF NOT EXISTS idx_player_builds_actor_id ON player_builds(actor_id);

-- 4. Add unique constraint for UPSERT operations
ALTER TABLE player_builds
    ADD CONSTRAINT unique_player_builds_report_fight_actor 
    UNIQUE (report_code, fight_id, actor_id);

-- 5. Finally, add all foreign key constraints
-- Check if dungeons table has the necessary unique constraint
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'uk_dungeons_encounter_id' 
        AND conrelid = 'dungeons'::regclass
    ) THEN
        ALTER TABLE dungeons 
        ADD CONSTRAINT uk_dungeons_encounter_id UNIQUE (encounter_id);
    END IF;
END $$;

-- Now add the foreign keys
ALTER TABLE player_builds
    ADD CONSTRAINT fk_player_builds_report 
    FOREIGN KEY (report_code, fight_id) 
    REFERENCES warcraft_logs_reports(code, fight_id)
    ON DELETE CASCADE;

ALTER TABLE player_builds
    ADD CONSTRAINT fk_player_builds_dungeon
    FOREIGN KEY (encounter_id)
    REFERENCES dungeons(encounter_id)
    ON DELETE CASCADE;

ALTER TABLE group_compositions
    ADD CONSTRAINT fk_group_compositions_report
    FOREIGN KEY (report_code, fight_id)
    REFERENCES warcraft_logs_reports(code, fight_id)
    ON DELETE CASCADE;