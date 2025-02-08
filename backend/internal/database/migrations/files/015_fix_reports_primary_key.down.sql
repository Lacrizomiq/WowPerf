-- +migrate Down

-- Remove foreign key constraints
ALTER TABLE player_builds
    DROP CONSTRAINT IF EXISTS fk_player_builds_report;

ALTER TABLE group_compositions
    DROP CONSTRAINT IF EXISTS fk_group_compositions_report;

-- Remove composite primary key and indexes
ALTER TABLE warcraft_logs_reports
    DROP CONSTRAINT IF EXISTS warcraft_logs_reports_pkey;

DROP INDEX IF EXISTS idx_warcraft_logs_reports_code;

-- Restore original primary key on code
ALTER TABLE warcraft_logs_reports
    ADD PRIMARY KEY (code);

-- Restore original foreign key constraints
ALTER TABLE player_builds
    ADD CONSTRAINT fk_player_builds_report
    FOREIGN KEY (report_code)
    REFERENCES warcraft_logs_reports(code)
    ON DELETE CASCADE;

ALTER TABLE group_compositions
    ADD CONSTRAINT group_compositions_report_code_fkey
    FOREIGN KEY (report_code)
    REFERENCES warcraft_logs_reports(code)
    ON DELETE CASCADE;