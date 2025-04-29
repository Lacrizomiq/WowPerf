-- UP MIGRATION

-- Add an encounter_id column to the get_popular_items_by_slot function
-- This function will be used to filter the items by encounter_id

-- Add the function to get the global popular items by slot
-- This function will be used to get the popular items across all encounters

DROP FUNCTION IF EXISTS get_popular_items_by_slot(TEXT, TEXT, INTEGER);
DROP FUNCTION IF EXISTS get_global_popular_items_by_slot(TEXT, TEXT);

-- Function: get_popular_items_by_slot (with encounter_id filtering)
-- This function will be used to get the popular items by slot and encounter_id
CREATE OR REPLACE FUNCTION get_popular_items_by_slot(p_class TEXT, p_spec TEXT, p_encounter_id INTEGER DEFAULT NULL) 
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
            ROW_NUMBER() OVER (PARTITION BY bs.encounter_id, bs.item_slot ORDER BY bs.usage_count DESC) as item_rank
        FROM build_statistics bs
        WHERE bs.class = p_class
        AND bs.spec = p_spec
        AND (p_encounter_id IS NULL OR bs.encounter_id = p_encounter_id)
    )
    SELECT 
        ri.encounter_id,
        ri.item_slot, 
        ri.item_id, 
        ri.item_name, 
        ri.item_icon,
        ri.item_quality,
        ri.item_level,
        ri.usage_count,
        ri.usage_percentage,
        ri.avg_keystone_level,
        ri.item_rank as rank
    FROM ranked_items ri
    WHERE ri.item_rank <= 4
    ORDER BY ri.encounter_id, ri.item_slot, ri.item_rank;
END;
$$ LANGUAGE plpgsql;

-- Function: get_global_popular_items_by_slot (global view across all encounters)
-- This function will be used to get the popular items across all encounters
CREATE OR REPLACE FUNCTION get_global_popular_items_by_slot(p_class TEXT, p_spec TEXT) 
RETURNS TABLE (
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
    WITH item_stats AS (
        -- Group all items with the same ID and sum their usages
        SELECT
            bs.item_slot,
            bs.item_id,
            MAX(bs.item_name)::TEXT AS item_name,
            MAX(bs.item_icon)::TEXT AS item_icon,
            MAX(bs.item_quality) AS item_quality,
            -- Use the max level for display
            MAX(bs.item_level) AS item_level,
            SUM(bs.usage_count) AS total_usage,
            SUM(bs.usage_count * bs.avg_keystone_level) / NULLIF(SUM(bs.usage_count), 0) AS avg_key_level
        FROM build_statistics bs
        WHERE bs.class = p_class AND bs.spec = p_spec
        GROUP BY bs.item_slot, bs.item_id
    ),
    slot_totals AS (
        -- Calculate the total for each slot
        SELECT 
            is_tot.item_slot, 
            SUM(is_tot.total_usage) AS slot_total
        FROM item_stats is_tot
        GROUP BY is_tot.item_slot
    ),
    final_ranked AS (
        -- Calculate the percentages and rank the items
        SELECT 
            is_rank.item_slot,
            is_rank.item_id,
            is_rank.item_name,
            is_rank.item_icon,
            is_rank.item_quality,
            is_rank.item_level,
            is_rank.total_usage::INTEGER AS usage_count,
            ROUND((is_rank.total_usage * 100.0 / NULLIF(st.slot_total, 0))::NUMERIC, 2) AS usage_percentage,
            ROUND(is_rank.avg_key_level::NUMERIC, 2) AS avg_keystone_level,
            ROW_NUMBER() OVER (PARTITION BY is_rank.item_slot ORDER BY is_rank.total_usage DESC) AS item_rank
        FROM item_stats is_rank
        JOIN slot_totals st ON is_rank.item_slot = st.item_slot
    )
    SELECT 
        fr.item_slot,
        fr.item_id,
        fr.item_name,
        fr.item_icon,
        fr.item_quality,
        fr.item_level,
        fr.usage_count,
        fr.usage_percentage,
        fr.avg_keystone_level,
        fr.item_rank as rank
    FROM final_ranked fr
    WHERE fr.item_rank <= 4
    ORDER BY fr.item_slot, fr.item_rank;
END;
$$ LANGUAGE plpgsql;