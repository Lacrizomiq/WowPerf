/*
 * MIGRATION: Remove Analytics Functions
 * 
 * This migration drops all the analytical functions that were added
 * in the corresponding up migration. Each function is dropped with its
 * full signature to ensure proper removal.
 *
 * The functions are dropped in reverse order of their creation to handle
 * any potential dependencies, although in this case they are independent.
 */

-- Drop the Bonus Function: Class Spec Summary
DROP FUNCTION IF EXISTS get_class_spec_summary(TEXT, TEXT);

-- Drop Function 8: Spec Comparison
DROP FUNCTION IF EXISTS get_spec_comparison(TEXT);

-- Drop Function 7: Optimal Build
DROP FUNCTION IF EXISTS get_optimal_build(TEXT, TEXT);

-- Drop Function 6: Stat Priorities
DROP FUNCTION IF EXISTS get_stat_priorities(TEXT, TEXT);

-- Drop Function 5: Talent Builds By Dungeon
DROP FUNCTION IF EXISTS get_talent_builds_by_dungeon(TEXT, TEXT);

-- Drop Function 4: Top Talent Builds
DROP FUNCTION IF EXISTS get_top_talent_builds(TEXT, TEXT);

-- Drop Function 3: Gem Usage
DROP FUNCTION IF EXISTS get_gem_usage(TEXT, TEXT);

-- Drop Function 2: Enchant Usage
DROP FUNCTION IF EXISTS get_enchant_usage(TEXT, TEXT);

-- Drop Function 1: Popular Items By Slot
DROP FUNCTION IF EXISTS get_popular_items_by_slot(TEXT, TEXT);