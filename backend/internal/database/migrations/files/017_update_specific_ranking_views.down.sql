-- Migration: 017_update_specific_ranking_views.down.sql
-- Purpose: Revert only the views that were modified
-- Date: February 26, 2025

-- Drop the modified views
DROP VIEW IF EXISTS spec_global_score_averages;
DROP VIEW IF EXISTS top_10_players_per_spec;

-- Recreate the original versions

-- Original spec_global_score_averages view
CREATE VIEW spec_global_score_averages AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        ROUND(SUM(score), 2) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, spec, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
)
SELECT 
    class,
    spec,
    ROUND(AVG(total_score), 2) AS avg_global_score,
    COUNT(*) AS player_count
FROM SpecPlayerScores
GROUP BY class, spec
ORDER BY avg_global_score DESC;

-- Original top_10_players_per_spec view
CREATE VIEW top_10_players_per_spec AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        ROUND(SUM(score), 2) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, spec, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
),
RankedPlayers AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        total_score,
        ROW_NUMBER() OVER (PARTITION BY class, spec ORDER BY total_score DESC) AS rank
    FROM SpecPlayerScores
)
SELECT 
    class,
    spec,
    name,
    server_name,
    server_region,
    total_score,
    rank
FROM RankedPlayers
WHERE rank <= 10 AND server_region <> 'CN'
ORDER BY class, spec, total_score DESC;