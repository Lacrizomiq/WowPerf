/*
 * MIGRATION: Revert Fix Remaining Issues
 *
 * This migration removes the fixed functions, allowing the previous versions
 * to be recreated if needed.
 */

-- Drop functions with new types
DROP FUNCTION IF EXISTS get_optimal_build(TEXT, TEXT);
DROP FUNCTION IF EXISTS get_stat_priorities(TEXT, TEXT);

-- Recreate previous functions from migration 029
-- Function 6: Get Stat Priorities (Previous version)
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
            ss.stat_name::TEXT,
            ss.stat_category::TEXT,
            AVG(ss.avg_value) as avg_value,
            MIN(ss.min_value) as min_value,
            MAX(ss.max_value) as max_value,
            SUM(ss.sample_size)::BIGINT as total_samples,
            AVG(ss.avg_keystone_level) as avg_keystone_level
        FROM stat_statistics ss
        WHERE ss.class = p_class
        AND ss.spec = p_spec
        GROUP BY ss.stat_name, ss.stat_category
    )
    SELECT 
        sa.stat_name,
        sa.stat_category,
        sa.avg_value,
        sa.min_value,
        sa.max_value,
        sa.total_samples,
        sa.avg_keystone_level,
        ROW_NUMBER() OVER (PARTITION BY sa.stat_category ORDER BY sa.avg_value DESC)::BIGINT as priority_rank
    FROM stat_aggregates sa
    ORDER BY sa.stat_category, ROW_NUMBER() OVER (PARTITION BY sa.stat_category ORDER BY sa.avg_value DESC);
END;
$$ LANGUAGE plpgsql;

-- Function 7: Get Optimal Build (Previous version)
CREATE OR REPLACE FUNCTION get_optimal_build(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    top_talent_import TEXT,
    stat_priority TEXT,
    top_items JSON
) AS $$
BEGIN
    RETURN QUERY
    WITH top_talents AS (
        SELECT ts.talent_import::TEXT
        FROM talent_statistics ts
        WHERE ts.class = p_class
        AND ts.spec = p_spec
        GROUP BY ts.talent_import
        ORDER BY SUM(ts.usage_count) DESC
        LIMIT 1
    ),
    top_stats AS (
        SELECT STRING_AGG(s.stat_name, ' > ' ORDER BY s.avg_value DESC)::TEXT as stat_priority
        FROM (
            SELECT 
                ss.stat_name::TEXT, AVG(ss.avg_value) as avg_value
            FROM stat_statistics ss
            WHERE ss.class = p_class
            AND ss.spec = p_spec
            AND ss.stat_category = 'secondary'
            GROUP BY ss.stat_name
        ) s
    ),
    top_items AS (
        SELECT 
            jsonb_object_agg(b.item_slot::text, 
                jsonb_build_object(
                    'name', b.item_name,
                    'icon', b.item_icon,
                    'quality', b.item_quality,
                    'usage_count', b.usage_count
                )
            ) as items_json
        FROM (
            SELECT DISTINCT ON (bs.item_slot)
                bs.item_slot, bs.item_name, bs.item_icon, bs.item_quality, bs.usage_count
            FROM build_statistics bs
            WHERE bs.class = p_class
            AND bs.spec = p_spec
            ORDER BY bs.item_slot, bs.usage_count DESC
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