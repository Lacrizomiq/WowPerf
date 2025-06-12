/*
 * MIGRATION: Fix Remaining Type Issues
 * 
 * This migration addresses:
 * 1. Type mismatch in get_stat_priorities (double precision vs numeric)
 * 2. NULL handling in get_optimal_build for top_items
 */

-- Drop existing functions before recreating them with different return types
DROP FUNCTION IF EXISTS get_stat_priorities(TEXT, TEXT);
DROP FUNCTION IF EXISTS get_optimal_build(TEXT, TEXT);

-- Function 6: Get Stat Priorities (Fixed Double Precision)
CREATE OR REPLACE FUNCTION get_stat_priorities(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    stat_name TEXT,
    stat_category TEXT,
    avg_value DOUBLE PRECISION, -- Changed from NUMERIC to DOUBLE PRECISION
    min_value DOUBLE PRECISION, -- Changed from NUMERIC to DOUBLE PRECISION
    max_value DOUBLE PRECISION, -- Changed from NUMERIC to DOUBLE PRECISION
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

-- Function 7: Get Optimal Build (Fixed NULL handling)
CREATE OR REPLACE FUNCTION get_optimal_build(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    top_talent_import TEXT,
    stat_priority TEXT,
    top_items JSON
) AS $$
BEGIN
    RETURN QUERY
    WITH top_talents AS (
        SELECT COALESCE(ts.talent_import::TEXT, '') AS talent_import
        FROM talent_statistics ts
        WHERE ts.class = p_class
        AND ts.spec = p_spec
        GROUP BY ts.talent_import
        ORDER BY SUM(ts.usage_count) DESC
        LIMIT 1
    ),
    top_stats AS (
        SELECT COALESCE(STRING_AGG(s.stat_name, ' > ' ORDER BY s.avg_value DESC), '')::TEXT as stat_priority
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
            COALESCE(
                jsonb_object_agg(
                    COALESCE(b.item_slot::text, '0'), 
                    jsonb_build_object(
                        'name', COALESCE(b.item_name, 'Unknown'),
                        'icon', COALESCE(b.item_icon, ''),
                        'quality', COALESCE(b.item_quality, 0),
                        'usage_count', COALESCE(b.usage_count, 0)
                    )
                ),
                '{}'::jsonb
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
        CASE 
            WHEN ti.items_json = '{}'::jsonb THEN NULL
            ELSE ti.items_json::json
        END as top_items
    FROM top_talents tt
    CROSS JOIN top_stats ts
    CROSS JOIN top_items ti;
END;
$$ LANGUAGE plpgsql;