/*
 * MIGRATION: Fix Type Issues in Analytics Functions
 * 
 * This migration addresses type compatibility issues between function return types and actual column types:
 * 1. Adds explicit type casts to ensure compatibility
 * 2. Fixes array handling for gem data
 * 3. Updates return type definitions to match actual column types
 */

-- Function 1: Get Popular Items By Slot (Fixed Types)
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
            bs.encounter_id,
            bs.item_slot, 
            bs.item_id, 
            bs.item_name::TEXT, 
            bs.item_icon::TEXT,
            bs.item_quality,
            bs.item_level,
            bs.usage_count,
            bs.usage_percentage,
            bs.avg_keystone_level,
            ROW_NUMBER() OVER (PARTITION BY bs.encounter_id, bs.item_slot ORDER BY bs.usage_count DESC) as rank
        FROM build_statistics bs
        WHERE bs.class = p_class
        AND bs.spec = p_spec
    )
    SELECT * FROM ranked_items ri
    WHERE ri.rank <= 4
    ORDER BY ri.encounter_id, ri.item_slot, ri.rank;
END;
$$ LANGUAGE plpgsql;

-- Function 2: Get Enchant Usage (Fixed Types)
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
        bs.item_slot,
        bs.permanent_enchant_id,
        bs.permanent_enchant_name::TEXT,
        COUNT(*)::BIGINT as usage_count,
        AVG(bs.avg_keystone_level) as avg_keystone_level,
        AVG(bs.avg_item_level) as avg_item_level,
        MAX(bs.max_keystone_level)::BIGINT as max_keystone_level,
        ROW_NUMBER() OVER (PARTITION BY bs.item_slot ORDER BY COUNT(*) DESC)::BIGINT as rank
    FROM build_statistics bs
    WHERE bs.class = p_class
    AND bs.spec = p_spec
    AND bs.has_permanent_enchant = true
    GROUP BY bs.item_slot, bs.permanent_enchant_id, bs.permanent_enchant_name
    ORDER BY bs.item_slot, ROW_NUMBER() OVER (PARTITION BY bs.item_slot ORDER BY COUNT(*) DESC);
END;
$$ LANGUAGE plpgsql;

-- Function 3: Get Gem Usage (Fixed Types)
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

-- Function 4: Get Top Talent Builds (Fixed Types)
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
        ts.talent_import::TEXT,
        SUM(ts.usage_count)::BIGINT as total_usage,
        AVG(ts.usage_percentage) as avg_usage_percentage,
        AVG(ts.avg_keystone_level) as avg_keystone_level
    FROM talent_statistics ts
    WHERE ts.class = p_class
    AND ts.spec = p_spec
    GROUP BY ts.talent_import
    ORDER BY SUM(ts.usage_count) DESC
    LIMIT 3;
END;
$$ LANGUAGE plpgsql;

-- Function 5: Get Talent Builds By Dungeon (Fixed Types)
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
        ts.class::TEXT,
        ts.spec::TEXT,
        ts.encounter_id,
        d.name::TEXT as dungeon_name,
        ts.talent_import::TEXT,
        SUM(ts.usage_count)::BIGINT as total_usage,
        AVG(ts.usage_percentage) as avg_usage_percentage,
        AVG(ts.avg_keystone_level) as avg_keystone_level
    FROM talent_statistics ts
    JOIN dungeons d ON ts.encounter_id = d.encounter_id
    WHERE ts.class = p_class
    AND ts.spec = p_spec
    GROUP BY ts.encounter_id, d.name, ts.talent_import, ts.class, ts.spec
    ORDER BY ts.class, ts.spec, ts.encounter_id, SUM(ts.usage_count) DESC;
END;
$$ LANGUAGE plpgsql;

-- Function 6: Get Stat Priorities (Fixed Types)
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

-- Function 7: Get Optimal Build (Fixed Types)
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

-- Function 8: Get Spec Comparison (Fixed Types)
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
            ts.spec::TEXT,
            AVG(ts.avg_keystone_level) as avg_keystone_level,
            MAX(ts.max_keystone_level)::BIGINT as max_keystone_level,
            AVG(ts.avg_item_level) as avg_item_level,
            COUNT(DISTINCT ts.encounter_id)::BIGINT as dungeons_count
        FROM talent_statistics ts
        WHERE ts.class = p_class
        GROUP BY ts.spec
    ),
    spec_stat_priorities AS (
        SELECT 
            sv.spec,
            STRING_AGG(sv.stat_name, ' > ' ORDER BY sv.avg_stat_val DESC)::TEXT as stat_priority
        FROM (
            SELECT 
                ss.spec::TEXT, 
                ss.stat_name::TEXT,
                AVG(ss.avg_value) as avg_stat_val
            FROM stat_statistics ss
            WHERE ss.class = p_class
            AND ss.stat_category = 'secondary'
            GROUP BY ss.spec, ss.stat_name
        ) as sv
        GROUP BY sv.spec
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

-- Bonus Function: Get Class Spec Summary (Fixed Types)
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
            AVG(ts.avg_keystone_level) as avg_keystone_level,
            MAX(ts.max_keystone_level)::BIGINT as max_keystone_level,
            AVG(ts.avg_item_level) as avg_item_level,
            COUNT(DISTINCT ts.encounter_id)::BIGINT as dungeons_count
        FROM talent_statistics ts
        WHERE ts.class = p_class
        AND ts.spec = p_spec
    ),
    top_talent AS (
        SELECT ts.talent_import::TEXT
        FROM talent_statistics ts
        WHERE ts.class = p_class
        AND ts.spec = p_spec
        GROUP BY ts.talent_import
        ORDER BY SUM(ts.usage_count) DESC
        LIMIT 1
    ),
    stat_priority AS (
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
    )
    SELECT 
        tst.avg_keystone_level,
        tst.max_keystone_level,
        tst.avg_item_level,
        tt.talent_import,
        sp.stat_priority,
        tst.dungeons_count
    FROM talent_stats tst
    CROSS JOIN top_talent tt
    CROSS JOIN stat_priority sp;
END;
$$ LANGUAGE plpgsql;