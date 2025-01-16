-- +migrate Up
-- Rename dungeon_id to encounter_id
ALTER TABLE player_builds RENAME COLUMN dungeon_id TO encounter_id;

-- Add NOT NULL constraints
ALTER TABLE player_builds 
    ALTER COLUMN report_code SET NOT NULL,
    ALTER COLUMN player_name SET NOT NULL,
    ALTER COLUMN class SET NOT NULL,
    ALTER COLUMN spec SET NOT NULL,
    ALTER COLUMN fight_id SET NOT NULL;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_player_builds_encounter_id ON player_builds(encounter_id);
CREATE INDEX IF NOT EXISTS idx_player_builds_report_code_fight_id ON player_builds(report_code, fight_id);
CREATE INDEX IF NOT EXISTS idx_player_builds_player_info ON player_builds(player_name, class, spec);
CREATE INDEX IF NOT EXISTS idx_player_builds_keystone_level ON player_builds(keystone_level);

-- Add foreign key to warcraft_logs_reports
ALTER TABLE player_builds
    ADD CONSTRAINT fk_player_builds_report 
    FOREIGN KEY (report_code, fight_id) 
    REFERENCES warcraft_logs_reports(code, fight_id)
    ON DELETE CASCADE;

-- Add foreign key to dungeons
ALTER TABLE player_builds
    ADD CONSTRAINT fk_player_builds_dungeon
    FOREIGN KEY (encounter_id)
    REFERENCES dungeons(encounter_id)
    ON DELETE CASCADE;