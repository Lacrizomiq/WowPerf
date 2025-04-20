/*
 * MIGRATION 031
 * Fix Gem and Optimal Build Functions
 * 
 * This migration addresses issues with null arrays in get_gem_usage function
 * and JSON conversion in get_optimal_build function
 */

-- Drop existing functions before recreating them
DROP FUNCTION IF EXISTS get_gem_usage(TEXT, TEXT);
DROP FUNCTION IF EXISTS get_optimal_build(TEXT, TEXT);

-- Function: Get Gem Usage (Fixed Null Arrays)
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
        COALESCE(bs.gem_ids, ARRAY[]::INTEGER[]) as gem_ids_array,
        COALESCE(bs.gem_icons, ARRAY[]::TEXT[]) as gem_icons_array,
        COALESCE(bs.gem_levels, ARRAY[]::NUMERIC[]) as gem_levels_array,
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

-- Function: Get Optimal Build (Fixed JSON handling)
CREATE OR REPLACE FUNCTION get_optimal_build(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
    top_talent_import TEXT,
    stat_priority TEXT,
    top_items JSON
) AS $$
DECLARE
    items_json JSONB;
BEGIN
    -- Get the top talent import
    SELECT COALESCE(ts.talent_import, '') INTO top_talent_import
    FROM talent_statistics ts
    WHERE ts.class = p_class
    AND ts.spec = p_spec
    GROUP BY ts.talent_import
    ORDER BY SUM(ts.usage_count) DESC
    LIMIT 1;
    
    -- Get the stat priority string
    SELECT COALESCE(STRING_AGG(s.stat_name, ' > ' ORDER BY s.avg_value DESC), '') INTO stat_priority
    FROM (
        SELECT 
            ss.stat_name, AVG(ss.avg_value) as avg_value
        FROM stat_statistics ss
        WHERE ss.class = p_class
        AND ss.spec = p_spec
        AND ss.stat_category = 'secondary'
        GROUP BY ss.stat_name
    ) s;
    
    -- Get the top items JSON
    SELECT 
        COALESCE(
            jsonb_object_agg(
                b.item_slot::text, 
                jsonb_build_object(
                    'name', COALESCE(b.item_name, 'Unknown'),
                    'icon', COALESCE(b.item_icon, ''),
                    'quality', COALESCE(b.item_quality, 0),
                    'usage_count', COALESCE(b.usage_count, 0)
                )
            ),
            '{}'::jsonb
        ) INTO items_json
    FROM (
        SELECT DISTINCT ON (bs.item_slot)
            bs.item_slot, bs.item_name, bs.item_icon, bs.item_quality, bs.usage_count
        FROM build_statistics bs
        WHERE bs.class = p_class
        AND bs.spec = p_spec
        ORDER BY bs.item_slot, bs.usage_count DESC
    ) b;
    
    -- Return a single row with all values
    RETURN QUERY
    SELECT 
        top_talent_import,
        stat_priority,
        CASE WHEN items_json = '{}'::jsonb THEN NULL ELSE items_json::json END;
END;
$$ LANGUAGE plpgsql;