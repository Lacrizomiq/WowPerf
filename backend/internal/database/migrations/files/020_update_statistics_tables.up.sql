-- Migration: 20_update_statistics_tables.up.sql

-- 1. Update build_statistics table
-- Drop columns
ALTER TABLE build_statistics
    DROP COLUMN IF EXISTS dungeon_id,
    DROP COLUMN IF EXISTS period_start,
    DROP COLUMN IF EXISTS period_end,
    DROP COLUMN IF EXISTS avg_score;

-- Add new columns for item details
ALTER TABLE build_statistics
    ADD COLUMN item_name VARCHAR(255),
    ADD COLUMN item_icon VARCHAR(255),
    ADD COLUMN item_quality INTEGER DEFAULT 0,
    ADD COLUMN item_level NUMERIC DEFAULT 0;

-- Add set bonus information
ALTER TABLE build_statistics
    ADD COLUMN has_set_bonus BOOLEAN DEFAULT FALSE,
    ADD COLUMN set_id INTEGER DEFAULT 0;

-- Add bonus IDs array
ALTER TABLE build_statistics
    ADD COLUMN bonus_ids INTEGER[];

-- Add gem information
ALTER TABLE build_statistics
    ADD COLUMN has_gems BOOLEAN DEFAULT FALSE,
    ADD COLUMN gems_count INTEGER DEFAULT 0,
    ADD COLUMN gem_ids INTEGER[],
    ADD COLUMN gem_icons TEXT[],
    ADD COLUMN gem_levels NUMERIC[];

-- Add enchant information
ALTER TABLE build_statistics
    ADD COLUMN has_permanent_enchant BOOLEAN DEFAULT FALSE,
    ADD COLUMN permanent_enchant_id INTEGER DEFAULT 0,
    ADD COLUMN permanent_enchant_name VARCHAR(255),
    ADD COLUMN has_temporary_enchant BOOLEAN DEFAULT FALSE,
    ADD COLUMN temporary_enchant_id INTEGER DEFAULT 0,
    ADD COLUMN temporary_enchant_name VARCHAR(255);

-- Add item level statistics
ALTER TABLE build_statistics
    ADD COLUMN avg_item_level NUMERIC DEFAULT 0,
    ADD COLUMN min_item_level NUMERIC DEFAULT 0,
    ADD COLUMN max_item_level NUMERIC DEFAULT 0;

-- Add keystone level statistics
ALTER TABLE build_statistics
    ADD COLUMN avg_keystone_level NUMERIC DEFAULT 0,
    ADD COLUMN min_keystone_level INTEGER DEFAULT 0,
    ADD COLUMN max_keystone_level INTEGER DEFAULT 0;

-- 2. Update talent_statistics table
-- Drop columns
ALTER TABLE talent_statistics
    DROP COLUMN IF EXISTS dungeon_id,
    DROP COLUMN IF EXISTS period_start,
    DROP COLUMN IF EXISTS period_end,
    DROP COLUMN IF EXISTS avg_score;

-- Add keystone level statistics
ALTER TABLE talent_statistics
    ADD COLUMN avg_keystone_level NUMERIC DEFAULT 0,
    ADD COLUMN min_keystone_level INTEGER DEFAULT 0,
    ADD COLUMN max_keystone_level INTEGER DEFAULT 0;

-- Add item level statistics
ALTER TABLE talent_statistics
    ADD COLUMN avg_item_level NUMERIC DEFAULT 0,
    ADD COLUMN min_item_level NUMERIC DEFAULT 0,
    ADD COLUMN max_item_level NUMERIC DEFAULT 0;

-- 3. Update stat_statistics table
-- Add sample size
ALTER TABLE stat_statistics
    ADD COLUMN sample_size INTEGER DEFAULT 0;

-- Add keystone level statistics
ALTER TABLE stat_statistics
    ADD COLUMN avg_keystone_level NUMERIC DEFAULT 0,
    ADD COLUMN min_keystone_level INTEGER DEFAULT 0,
    ADD COLUMN max_keystone_level INTEGER DEFAULT 0;

-- Add item level statistics
ALTER TABLE stat_statistics
    ADD COLUMN avg_item_level NUMERIC DEFAULT 0,
    ADD COLUMN min_item_level NUMERIC DEFAULT 0,
    ADD COLUMN max_item_level NUMERIC DEFAULT 0;

-- Update indexes for build_statistics
DROP INDEX IF EXISTS idx_build_statistics_dungeon;
CREATE INDEX IF NOT EXISTS idx_build_statistics_item_name ON build_statistics(item_name);
CREATE INDEX IF NOT EXISTS idx_build_statistics_item_quality ON build_statistics(item_quality);
CREATE INDEX IF NOT EXISTS idx_build_statistics_set_id ON build_statistics(set_id);
CREATE INDEX IF NOT EXISTS idx_build_statistics_has_gems ON build_statistics(has_gems);
CREATE INDEX IF NOT EXISTS idx_build_statistics_has_permanent_enchant ON build_statistics(has_permanent_enchant);
CREATE INDEX IF NOT EXISTS idx_build_statistics_keystone_level ON build_statistics(avg_keystone_level);

-- Update indexes for stat_statistics
CREATE INDEX IF NOT EXISTS idx_stat_statistics_stat_name ON stat_statistics(stat_name);
CREATE INDEX IF NOT EXISTS idx_stat_statistics_avg_value ON stat_statistics(avg_value);
CREATE INDEX IF NOT EXISTS idx_stat_statistics_keystone_level ON stat_statistics(avg_keystone_level);