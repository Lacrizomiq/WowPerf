/*
 * MIGRATION: Add Analytics Functions
 * 
 * This migration adds a set of SQL functions that provide analytical insights
 * for character builds, talent configurations, enchants, gems, statistics, and more.
 * These functions transform raw data into actionable insights for players.
 *
 * Each function accepts class and spec parameters (except get_spec_comparison which only needs class)
 * and returns structured data that can be used by the application to present analytics.
 */

-- Function 1: Get Popular Items By Slot
-- Returns the most popular items for each slot for a specific class and spec
-- Limits to top 4 items per slot, ranked by usage count
CREATE OR REPLACE FUNCTION get_popular_items_by_slot(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    encounter_id INTEGER,
    item_slot INTEGER,
    item_id INTEGER,
    item_name TEXT,
    item_icon TEXT,
    item_quality INTEGER,
    item_level NUMERIC,
    usage_count INTEGER,
    usage_percentage NUMERIC,
    avg_keystone_level NUMERIC,
    rank BIGINT
) AS $$
BEGIN
    RETURN QUERY
    WITH ranked_items AS (
        SELECT 
            encounter_id,
            item_slot, 
            item_id, 
            item_name, 
            item_icon,
            item_quality,
            item_level,
            usage_count,
            usage_percentage,
            avg_keystone_level,
            ROW_NUMBER() OVER (PARTITION BY encounter_id, item_slot ORDER BY usage_count DESC) as rank
        FROM build_statistics
        WHERE class = p_class
        AND spec = p_spec
    )
    SELECT * FROM ranked_items
    WHERE rank <= 4
    ORDER BY encounter_id, item_slot, rank;
END;
$$ LANGUAGE plpgsql;

-- Function 2: Get Enchant Usage
-- Analyzes enchantment popularity for a given class and spec
-- Provides usage statistics and ranks enchants by popularity
CREATE OR REPLACE FUNCTION get_enchant_usage(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    item_slot INTEGER,
    permanent_enchant_id INTEGER,
    permanent_enchant_name TEXT,
    usage_count BIGINT,
    avg_keystone_level NUMERIC,
    avg_item_level NUMERIC,
    max_keystone_level BIGINT,
    rank BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        item_slot,
        permanent_enchant_id,
        permanent_enchant_name,
        COUNT(*) as usage_count,
        AVG(avg_keystone_level) as avg_keystone_level,
        AVG(avg_item_level) as avg_item_level,
        MAX(max_keystone_level) as max_keystone_level,
        ROW_NUMBER() OVER (PARTITION BY item_slot ORDER BY COUNT(*) DESC) as rank
    FROM build_statistics
    WHERE class = p_class
    AND spec = p_spec
    AND has_permanent_enchant = true
    GROUP BY item_slot, permanent_enchant_id, permanent_enchant_name
    ORDER BY item_slot, rank;
END;
$$ LANGUAGE plpgsql;

-- Function 3: Get Gem Usage
-- Analyzes gem configurations and popularity for a given class and spec
-- Returns arrays of gem IDs, icons and levels for each slot configuration
CREATE OR REPLACE FUNCTION get_gem_usage(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    item_slot INTEGER,
    gems_count INTEGER,
    gem_ids_array INTEGER[],
    gem_icons_array TEXT[],
    gem_levels_array NUMERIC[],
    usage_count BIGINT,
    avg_keystone_level NUMERIC,
    avg_item_level NUMERIC,
    rank BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        item_slot,
        gems_count,
        gem_ids as gem_ids_array,
        gem_icons as gem_icons_array,
        gem_levels as gem_levels_array,
        COUNT(*) as usage_count,
        AVG(avg_keystone_level) as avg_keystone_level,
        AVG(avg_item_level) as avg_item_level,
        ROW_NUMBER() OVER (PARTITION BY item_slot ORDER BY COUNT(*) DESC) as rank
    FROM build_statistics
    WHERE class = p_class
    AND spec = p_spec
    AND has_gems = true
    GROUP BY item_slot, gems_count, gem_ids, gem_icons, gem_levels
    ORDER BY item_slot, rank;
END;
$$ LANGUAGE plpgsql;

-- Function 4: Get Top Talent Builds
-- Returns the 3 most popular talent configurations for a given class and spec
-- Includes usage statistics and average keystone level
CREATE OR REPLACE FUNCTION get_top_talent_builds(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    talent_import TEXT,
    total_usage BIGINT,
    avg_usage_percentage NUMERIC,
    avg_keystone_level NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        talent_import,
        SUM(usage_count) as total_usage,
        AVG(usage_percentage) as avg_usage_percentage,
        AVG(avg_keystone_level) as avg_keystone_level
    FROM talent_statistics
    WHERE class = p_class
    AND spec = p_spec
    GROUP BY talent_import
    ORDER BY total_usage DESC
    LIMIT 3;
END;
$$ LANGUAGE plpgsql;

-- Function 5: Get Talent Builds By Dungeon
-- Analyzes talent configurations per dungeon for a specific class and spec
-- Joins with dungeons table to include dungeon names
CREATE OR REPLACE FUNCTION get_talent_builds_by_dungeon(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    class TEXT,
    spec TEXT,
    encounter_id INTEGER,
    dungeon_name TEXT,
    talent_import TEXT,
    total_usage BIGINT,
    avg_usage_percentage NUMERIC,
    avg_keystone_level NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ts.class,
        ts.spec,
        ts.encounter_id,
        d.name as dungeon_name,
        ts.talent_import,
        SUM(ts.usage_count) as total_usage,
        AVG(ts.usage_percentage) as avg_usage_percentage,
        AVG(ts.avg_keystone_level) as avg_keystone_level
    FROM talent_statistics ts
    JOIN dungeons d ON ts.encounter_id = d.encounter_id
    WHERE ts.class = p_class
    AND ts.spec = p_spec
    GROUP BY ts.encounter_id, d.name, ts.talent_import, ts.class, ts.spec
    ORDER BY ts.class, ts.spec, ts.encounter_id, total_usage DESC;
END;
$$ LANGUAGE plpgsql;

-- Function 6: Get Stat Priorities
-- Analyzes statistical priorities for a given class and spec
-- Ranks stats within each category (primary, secondary, etc.)
CREATE OR REPLACE FUNCTION get_stat_priorities(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    stat_name TEXT,
    stat_category TEXT,
    avg_value NUMERIC,
    min_value NUMERIC,
    max_value NUMERIC,
    total_samples BIGINT,
    avg_keystone_level NUMERIC,
    priority_rank BIGINT
) AS $$
BEGIN
    RETURN QUERY
    WITH stat_aggregates AS (
        SELECT 
            stat_name,
            stat_category,
            AVG(avg_value) as avg_value,
            MIN(min_value) as min_value,
            MAX(max_value) as max_value,
            SUM(sample_size) as total_samples,
            AVG(avg_keystone_level) as avg_keystone_level
        FROM stat_statistics
        WHERE class = p_class
        AND spec = p_spec
        GROUP BY stat_name, stat_category
    )
    SELECT 
        stat_name,
        stat_category,
        avg_value,
        min_value,
        max_value,
        total_samples,
        avg_keystone_level,
        ROW_NUMBER() OVER (PARTITION BY stat_category ORDER BY avg_value DESC) as priority_rank
    FROM stat_aggregates
    ORDER BY stat_category, priority_rank;
END;
$$ LANGUAGE plpgsql;

-- Function 7: Get Optimal Build
-- Provides a comprehensive overview of the optimal build for a class/spec
-- Combines top talent, stat priority, and best-in-slot items into a single result
CREATE OR REPLACE FUNCTION get_optimal_build(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    top_talent_import TEXT,
    stat_priority TEXT,
    top_items JSON
) AS $$
BEGIN
    RETURN QUERY
    WITH top_talents AS (
        SELECT talent_import
        FROM talent_statistics
        WHERE class = p_class
        AND spec = p_spec
        GROUP BY talent_import
        ORDER BY SUM(usage_count) DESC
        LIMIT 1
    ),
    top_stats AS (
        SELECT STRING_AGG(stat_name, ' > ' ORDER BY avg_value DESC) as stat_priority
        FROM (
            SELECT 
                stat_name, AVG(avg_value) as avg_value
            FROM stat_statistics
            WHERE class = p_class
            AND spec = p_spec
            AND stat_category = 'secondary'
            GROUP BY stat_name
        ) s
    ),
    top_items AS (
        SELECT 
            jsonb_object_agg(item_slot::text, 
                jsonb_build_object(
                    'name', item_name,
                    'icon', item_icon,
                    'quality', item_quality,
                    'usage_count', usage_count
                )
            ) as items_json
        FROM (
            SELECT DISTINCT ON (item_slot)
                item_slot, item_name, item_icon, item_quality, usage_count
            FROM build_statistics
            WHERE class = p_class
            AND spec = p_spec
            ORDER BY item_slot, usage_count DESC
        ) b
    )
    SELECT 
        tt.talent_import,
        ts.stat_priority,
        ti.items_json::json
    FROM top_talents tt
    CROSS JOIN top_stats ts
    CROSS JOIN top_items ti;
END;
$$ LANGUAGE plpgsql;

-- Function 8: Get Spec Comparison
-- Compares different specializations within a class on various metrics
-- Provides stat priorities and performance metrics for each spec
CREATE OR REPLACE FUNCTION get_spec_comparison(p_class TEXT) 
RETURNS TABLE (
    spec TEXT,
    avg_keystone_level NUMERIC,
    max_keystone_level BIGINT,
    avg_item_level NUMERIC,
    dungeons_count BIGINT,
    stat_priority TEXT
) AS $$
BEGIN
    RETURN QUERY
    WITH spec_stats AS (
        SELECT 
            spec,
            AVG(avg_keystone_level) as avg_keystone_level,
            MAX(max_keystone_level) as max_keystone_level,
            AVG(avg_item_level) as avg_item_level,
            COUNT(DISTINCT encounter_id) as dungeons_count
        FROM talent_statistics
        WHERE class = p_class
        GROUP BY spec
    ),
    spec_stat_priorities AS (
        SELECT 
            spec,
            STRING_AGG(stat_name, ' > ' ORDER BY avg_stat_val DESC) as stat_priority
        FROM (
            SELECT 
                spec, 
                stat_name,
                AVG(avg_value) as avg_stat_val
            FROM stat_statistics
            WHERE class = p_class
            AND stat_category = 'secondary'
            GROUP BY spec, stat_name
        ) as stat_values
        GROUP BY spec
    )
    SELECT 
        ss.spec,
        ss.avg_keystone_level,
        ss.max_keystone_level,
        ss.avg_item_level,
        ss.dungeons_count,
        ssp.stat_priority
    FROM spec_stats ss
    LEFT JOIN spec_stat_priorities ssp ON ss.spec = ssp.spec
    ORDER BY ss.avg_keystone_level DESC;
END;
$$ LANGUAGE plpgsql;

-- Bonus Function: Get Class Spec Summary
-- Provides a concise summary of a class/spec performance and configuration
-- Combines key metrics into a single, easily accessible result
CREATE OR REPLACE FUNCTION get_class_spec_summary(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    avg_keystone_level NUMERIC,
    max_keystone_level BIGINT,
    avg_item_level NUMERIC,
    top_talent_import TEXT,
    stat_priority TEXT,
    dungeons_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    WITH talent_stats AS (
        SELECT 
            AVG(avg_keystone_level) as avg_keystone_level,
            MAX(max_keystone_level) as max_keystone_level,
            AVG(avg_item_level) as avg_item_level,
            COUNT(DISTINCT encounter_id) as dungeons_count
        FROM talent_statistics
        WHERE class = p_class
        AND spec = p_spec
    ),
    top_talent AS (
        SELECT talent_import
        FROM talent_statistics
        WHERE class = p_class
        AND spec = p_spec
        GROUP BY talent_import
        ORDER BY SUM(usage_count) DESC
        LIMIT 1
    ),
    stat_priority AS (
        SELECT STRING_AGG(stat_name, ' > ' ORDER BY avg_value DESC) as stat_priority
        FROM (
            SELECT 
                stat_name, AVG(avg_value) as avg_value
            FROM stat_statistics
            WHERE class = p_class
            AND spec = p_spec
            AND stat_category = 'secondary'
            GROUP BY stat_name
        ) s
    )
    SELECT 
        ts.avg_keystone_level,
        ts.max_keystone_level,
        ts.avg_item_level,
        tt.talent_import,
        sp.stat_priority,
        ts.dungeons_count
    FROM talent_stats ts
    CROSS JOIN top_talent tt
    CROSS JOIN stat_priority sp;
END;
$$ LANGUAGE plpgsql;