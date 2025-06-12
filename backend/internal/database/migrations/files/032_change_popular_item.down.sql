-- DOWN MIGRATION

-- Drop the new functions
DROP FUNCTION IF EXISTS get_popular_items_by_slot(TEXT, TEXT, INTEGER);
DROP FUNCTION IF EXISTS get_global_popular_items_by_slot(TEXT, TEXT);

-- Restore the original get_popular_items_by_slot function without encounter_id parameter
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