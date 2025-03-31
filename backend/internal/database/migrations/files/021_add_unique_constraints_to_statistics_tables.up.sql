-- Migration 021: Add unique constraints to statistics tables
ALTER TABLE build_statistics ADD CONSTRAINT IF NOT EXISTS build_statistics_unique
UNIQUE (class, spec, encounter_id, item_slot, item_id);

ALTER TABLE talent_statistics ADD CONSTRAINT IF NOT EXISTS talent_statistics_unique
UNIQUE (class, spec, encounter_id, talent_import);

ALTER TABLE stat_statistics ADD CONSTRAINT IF NOT EXISTS stat_statistics_unique
UNIQUE (class, spec, encounter_id, stat_name);