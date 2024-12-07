-- Add missing columns to match the Report model structure
ALTER TABLE warcraft_logs_reports 
    -- Add missing mythic+ data columns
    ADD COLUMN IF NOT EXISTS keystone_level INTEGER,       -- M+ key level
    ADD COLUMN IF NOT EXISTS keystone_time BIGINT,        -- Time taken to complete the key
    ADD COLUMN IF NOT EXISTS affixes INTEGER[],           -- Array of affix IDs for the key

    -- Add missing performance data columns
    ADD COLUMN IF NOT EXISTS damage_taken JSONB,          -- Damage taken data for the run

    -- Rename any columns to match exactly with the model
    -- None needed as columns already match model names

    -- Set correct column types if different
    ALTER COLUMN fight_id SET DATA TYPE INTEGER,
    ALTER COLUMN encounter_id SET DATA TYPE INTEGER,
    ALTER COLUMN total_time SET DATA TYPE BIGINT,
    ALTER COLUMN item_level SET DATA TYPE NUMERIC,

    -- Set correct column constraints
    ALTER COLUMN fight_id SET NOT NULL,
    ALTER COLUMN code SET NOT NULL;

-- Remove any old indexes that might conflict
DROP INDEX IF EXISTS idx_reports_keystone;

-- Create or update all necessary indexes to match the model
CREATE INDEX IF NOT EXISTS idx_reports_encounter_id ON warcraft_logs_reports(encounter_id);
CREATE INDEX IF NOT EXISTS idx_reports_keystone_level ON warcraft_logs_reports(keystone_level);
CREATE INDEX IF NOT EXISTS idx_reports_keystone_time ON warcraft_logs_reports(keystone_time);

-- Performance data indexes
CREATE INDEX IF NOT EXISTS idx_reports_damage_taken ON warcraft_logs_reports USING gin (damage_taken jsonb_path_ops);

-- Add proper foreign key constraints if needed
ALTER TABLE warcraft_logs_reports
    DROP CONSTRAINT IF EXISTS fk_warcraft_logs_reports_encounter_id,
    ADD CONSTRAINT fk_warcraft_logs_reports_encounter_id 
        FOREIGN KEY (encounter_id) 
        REFERENCES dungeons(encounter_id);