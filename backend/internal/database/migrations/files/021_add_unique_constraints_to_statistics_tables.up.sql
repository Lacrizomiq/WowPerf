-- Migration 021: Add unique constraints to statistics tables
-- First drop existing constraints if they exist
ALTER TABLE build_statistics DROP CONSTRAINT IF EXISTS build_statistics_unique;
ALTER TABLE talent_statistics DROP CONSTRAINT IF EXISTS talent_statistics_unique;
ALTER TABLE stat_statistics DROP CONSTRAINT IF EXISTS stat_statistics_unique;

-- Then add the constraints
ALTER TABLE build_statistics ADD CONSTRAINT build_statistics_unique
UNIQUE (class, spec, encounter_id, item_slot, item_id);

ALTER TABLE talent_statistics ADD CONSTRAINT talent_statistics_unique
UNIQUE (class, spec, encounter_id, talent_import);

ALTER TABLE stat_statistics ADD CONSTRAINT stat_statistics_unique
UNIQUE (class, spec, encounter_id, stat_name);