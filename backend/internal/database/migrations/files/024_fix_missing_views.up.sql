-- Migration: 024_fix_missing_views.up.sql
-- Purpose: Recreate missing views that were dropped but not recreated in previous migrations
--          and update top_5_players_per_role to exclude CN players
-- Date: April 04, 2025

-- Recreate spec_dungeon_max_key_levels view
CREATE OR REPLACE VIEW spec_dungeon_max_key_levels AS
SELECT 
    pr.class,
    pr.spec,
    d.name AS dungeon_name,
    d.slug AS dungeon_slug,
    MAX(pr.hard_mode_level) AS max_key_level
FROM player_rankings pr
JOIN dungeons d ON pr.dungeon_id = d.encounter_id
WHERE pr.deleted_at IS NULL
GROUP BY pr.class, pr.spec, d.name, d.slug
ORDER BY pr.class, pr.spec, max_key_level DESC;

-- Recreate class_global_score_averages view if it's also missing
CREATE OR REPLACE VIEW class_global_score_averages AS
WITH ClassPlayerScores AS (
    SELECT 
        class,
        name,
        server_name,
        server_region,
        CAST(SUM(score) AS numeric(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
)
SELECT 
    class,
    CAST(AVG(total_score) AS numeric(10,2)) AS avg_global_score,
    COUNT(*) AS player_count
FROM ClassPlayerScores
GROUP BY class
ORDER BY avg_global_score DESC;

-- Recreate dungeon_avg_key_levels view if it's also missing
CREATE OR REPLACE VIEW dungeon_avg_key_levels AS
SELECT 
    d.name AS dungeon_name,
    d.slug AS dungeon_slug,
    CAST(AVG(pr.hard_mode_level) AS numeric(10,2)) AS avg_key_level,
    COUNT(*) AS run_count
FROM player_rankings pr
JOIN dungeons d ON pr.dungeon_id = d.encounter_id
WHERE pr.deleted_at IS NULL
GROUP BY d.name, d.slug
ORDER BY avg_key_level DESC;

-- Recreate top_5_players_per_role view with CN players excluded
CREATE OR REPLACE VIEW top_5_players_per_role AS
WITH RolePlayerScores AS (
    SELECT 
        name,
        server_name,
        server_region,
        class,
        spec,
        role,
        CAST(SUM(score) AS numeric(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL 
    AND server_region <> 'CN' -- Exclude CN players
    GROUP BY name, server_name, server_region, class, spec, role
    HAVING COUNT(DISTINCT dungeon_id) = 8
),
RankedPlayers AS (
    SELECT 
        name,
        server_name,
        server_region,
        class,
        spec,
        role,
        total_score,
        ROW_NUMBER() OVER (PARTITION BY role ORDER BY total_score DESC) AS rank
    FROM RolePlayerScores
)
SELECT 
    name,
    server_name,
    server_region,
    class,
    spec,
    role,
    total_score,
    rank
FROM RankedPlayers
WHERE rank <= 5
ORDER BY role, total_score DESC;