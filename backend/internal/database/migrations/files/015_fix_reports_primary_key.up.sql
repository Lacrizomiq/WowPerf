-- +migrate Up

-- First, remove existing foreign key constraints
ALTER TABLE player_builds
    DROP CONSTRAINT IF EXISTS fk_player_builds_report;

ALTER TABLE group_compositions
    DROP CONSTRAINT IF EXISTS group_compositions_report_code_fkey;

-- Drop any existing primary key and unique constraints
ALTER TABLE warcraft_logs_reports
    DROP CONSTRAINT IF EXISTS warcraft_logs_reports_pkey CASCADE;

DROP INDEX IF EXISTS warcraft_logs_reports_code_key;
DROP INDEX IF EXISTS warcraft_logs_reports_code_fight_id_key;

-- Create new composite primary key
ALTER TABLE warcraft_logs_reports
    ADD PRIMARY KEY (code, fight_id);

-- Add index on code alone for performance (only if it doesn't exist)
DROP INDEX IF EXISTS idx_warcraft_logs_reports_code;
CREATE INDEX idx_warcraft_logs_reports_code 
    ON warcraft_logs_reports(code);

-- Recreate foreign key constraints with new composite key
ALTER TABLE player_builds
    ADD CONSTRAINT fk_player_builds_report 
    FOREIGN KEY (report_code, fight_id) 
    REFERENCES warcraft_logs_reports(code, fight_id)
    ON DELETE CASCADE;

ALTER TABLE group_compositions
    ADD CONSTRAINT fk_group_compositions_report
    FOREIGN KEY (report_code, fight_id)
    REFERENCES warcraft_logs_reports(code, fight_id)
    ON DELETE CASCADE;