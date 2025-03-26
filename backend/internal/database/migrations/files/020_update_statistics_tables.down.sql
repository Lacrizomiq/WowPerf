-- Migration: 20_update_statistics_tables.down.sql

-- 1. Revert build_statistics table changes
-- Remove added columns - item details
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS item_name,
    DROP COLUMN IF EXISTS item_icon,
    DROP COLUMN IF EXISTS item_quality,
    DROP COLUMN IF EXISTS item_level;

-- Remove set bonus information
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS has_set_bonus,
    DROP COLUMN IF EXISTS set_id;

-- Remove bonus IDs array
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS bonus_ids;

-- Remove gem information
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS has_gems,
    DROP COLUMN IF EXISTS gems_count,
    DROP COLUMN IF EXISTS gem_ids,
    DROP COLUMN IF EXISTS gem_icons,
    DROP COLUMN IF EXISTS gem_levels;

-- Remove enchant information
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS has_permanent_enchant,
    DROP COLUMN IF EXISTS permanent_enchant_id,
    DROP COLUMN IF EXISTS permanent_enchant_name,
    DROP COLUMN IF EXISTS has_temporary_enchant,
    DROP COLUMN IF EXISTS temporary_enchant_id,
    DROP COLUMN IF EXISTS temporary_enchant_name;

-- Remove item level statistics
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS avg_item_level,
    DROP COLUMN IF EXISTS min_item_level,
    DROP COLUMN IF EXISTS max_item_level;

-- Remove keystone level statistics
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS avg_keystone_level,
    DROP COLUMN IF EXISTS min_keystone_level,
    DROP COLUMN IF EXISTS max_keystone_level;

-- Add back removed columns
ALTER TABLE build_statistics
    ADD COLUMN dungeon_id INTEGER,
    ADD COLUMN period_start TIMESTAMP WITH TIME ZONE,
    ADD COLUMN period_end TIMESTAMP WITH TIME ZONE,
    ADD COLUMN avg_score NUMERIC DEFAULT 0;

-- 2. Revert talent_statistics table changes
-- Remove keystone level statistics
ALTER TABLE talent_statistics
    DROP COLUMN IF EXISTS avg_keystone_level,
    DROP COLUMN IF EXISTS min_keystone_level,
    DROP COLUMN IF EXISTS max_keystone_level;

-- Remove item level statistics
ALTER TABLE talent_statistics
    DROP COLUMN IF EXISTS avg_item_level,
    DROP COLUMN IF EXISTS min_item_level,
    DROP COLUMN IF EXISTS max_item_level;

-- Add back removed columns
ALTER TABLE talent_statistics
    ADD COLUMN dungeon_id INTEGER,
    ADD COLUMN period_start TIMESTAMP WITH TIME ZONE,
    ADD COLUMN period_end TIMESTAMP WITH TIME ZONE,
    ADD COLUMN avg_score NUMERIC DEFAULT 0;

-- 3. Revert stat_statistics table changes
-- Remove sample size
ALTER TABLE stat_statistics
    DROP COLUMN IF EXISTS sample_size;

-- Remove keystone level statistics
ALTER TABLE stat_statistics
    DROP COLUMN IF EXISTS avg_keystone_level,
    DROP COLUMN IF EXISTS min_keystone_level,
    DROP COLUMN IF EXISTS max_keystone_level;

-- Remove item level statistics
ALTER TABLE stat_statistics
    DROP COLUMN IF EXISTS avg_item_level,
    DROP COLUMN IF EXISTS min_item_level,
    DROP COLUMN IF EXISTS max_item_level;

-- Remove added indexes
DROP INDEX IF EXISTS idx_build_statistics_item_name;
DROP INDEX IF EXISTS idx_build_statistics_item_quality;
DROP INDEX IF EXISTS idx_build_statistics_set_id;
DROP INDEX IF EXISTS idx_build_statistics_has_gems;
DROP INDEX IF EXISTS idx_build_statistics_has_permanent_enchant;
DROP INDEX IF EXISTS idx_build_statistics_keystone_level;
DROP INDEX IF EXISTS idx_stat_statistics_stat_name;
DROP INDEX IF EXISTS idx_stat_statistics_avg_value;
DROP INDEX IF EXISTS idx_stat_statistics_keystone_level;

-- Recreate original indexes
CREATE INDEX IF NOT EXISTS idx_build_statistics_dungeon ON build_statistics(dungeon_id);