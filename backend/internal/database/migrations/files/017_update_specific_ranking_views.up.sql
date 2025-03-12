-- Migration: 017_update_specific_ranking_views.up.sql
-- Purpose: Update only the views that need enhancements
-- Date: February 26, 2025

-- Drop and recreate only the views that need updates
DROP VIEW IF EXISTS spec_global_score_averages;
DROP VIEW IF EXISTS top_10_players_per_spec;

-- Create enhanced view for average global score per spec
-- Additions: 
-- - Role information
-- - Slug for easy URL creation
-- - Overall rank and rank within role
CREATE VIEW spec_global_score_averages AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        ROUND(SUM(score), 2) AS total_score,
        MIN(role) AS role -- Role should be consistent for a spec
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, spec, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
),
SpecAverages AS (
    SELECT 
        class,
        spec,
        ROUND(AVG(total_score), 2) AS avg_global_score,
        COUNT(*) AS player_count,
        MIN(role) AS role
    FROM SpecPlayerScores
    GROUP BY class, spec
)
SELECT 
    class,
    spec,
    LOWER(CONCAT(class, '-', REPLACE(spec, ' ', '-'))) AS slug, -- Create slug for URLs
    avg_global_score,
    player_count,
    role,
    RANK() OVER (ORDER BY avg_global_score DESC) AS overall_rank, -- Overall ranking
    RANK() OVER (PARTITION BY role ORDER BY avg_global_score DESC) AS role_rank -- Ranking within role
FROM SpecAverages
ORDER BY avg_global_score DESC;

-- Create improved view for top 10 players per spec
-- Fixed: CN players are filtered at the beginning to ensure full 10 players per spec
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
    AND server_region <> 'CN' -- Filter CN players early to get complete top 10
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
WHERE rank <= 10
ORDER BY class, spec, total_score DESC;