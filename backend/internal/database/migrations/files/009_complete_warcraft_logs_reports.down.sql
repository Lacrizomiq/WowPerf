ALTER TABLE warcraft_logs_reports
    DROP COLUMN IF EXISTS keystone_level,
    DROP COLUMN IF EXISTS keystone_time,
    DROP COLUMN IF EXISTS damage_taken,
    DROP COLUMN IF EXISTS affixes;

DROP INDEX IF EXISTS idx_reports_encounter_id;
DROP INDEX IF EXISTS idx_reports_keystone_level;
DROP INDEX IF EXISTS idx_reports_keystone_time;
DROP INDEX IF EXISTS idx_reports_damage_taken;

ALTER TABLE warcraft_logs_reports
    DROP CONSTRAINT IF EXISTS fk_warcraft_logs_reports_encounter_id;