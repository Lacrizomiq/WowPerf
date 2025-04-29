/*
 * DOWN MIGRATION 031
 * Restore previous versions of functions from migration 030
 */

-- Drop the fixed functions
DROP FUNCTION IF EXISTS get_gem_usage(TEXT, TEXT);
DROP FUNCTION IF EXISTS get_optimal_build(TEXT, TEXT);

-- Restore Get Gem Usage function from migration 029
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
        bs.item_slot,
        bs.gems_count,
        CASE WHEN bs.gem_ids IS NULL THEN '{}'::INTEGER[] ELSE bs.gem_ids END as gem_ids_array,
        CASE WHEN bs.gem_icons IS NULL THEN '{}'::TEXT[] ELSE bs.gem_icons END as gem_icons_array,
        CASE WHEN bs.gem_levels IS NULL THEN '{}'::NUMERIC[] ELSE bs.gem_levels END as gem_levels_array,
        COUNT(*)::BIGINT as usage_count,
        AVG(bs.avg_keystone_level) as avg_keystone_level,
        AVG(bs.avg_item_level) as avg_item_level,
        ROW_NUMBER() OVER (PARTITION BY bs.item_slot ORDER BY COUNT(*) DESC)::BIGINT as rank
    FROM build_statistics bs
    WHERE bs.class = p_class
    AND bs.spec = p_spec
    AND bs.has_gems = true
    GROUP BY bs.item_slot, bs.gems_count, bs.gem_ids, bs.gem_icons, bs.gem_levels
    ORDER BY bs.item_slot, ROW_NUMBER() OVER (PARTITION BY bs.item_slot ORDER BY COUNT(*) DESC);
END;
$$ LANGUAGE plpgsql;

-- Restore Get Optimal Build function from migration 030
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