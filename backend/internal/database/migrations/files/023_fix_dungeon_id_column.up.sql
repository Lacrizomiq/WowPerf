-- First drop the dependent views
DROP VIEW IF EXISTS spec_global_score_averages;
DROP VIEW IF EXISTS top_10_players_per_spec;
DROP VIEW IF EXISTS spec_dungeon_max_key_levels;
DROP VIEW IF EXISTS class_global_score_averages;
DROP VIEW IF EXISTS dungeon_avg_key_levels;
DROP VIEW IF EXISTS top_5_players_per_role;

-- Now alter the column
ALTER TABLE player_rankings ALTER COLUMN dungeon_id TYPE bigint;

-- Recreate the views
CREATE VIEW spec_global_score_averages AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        CAST(SUM(score) AS numeric(10,2)) AS total_score,
        MIN(role) AS role
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, spec, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
),
SpecAverages AS (
    SELECT 
        class,
        spec,
        CAST(AVG(total_score) AS numeric(10,2)) AS avg_global_score,
        COUNT(*) AS player_count,
        MIN(role) AS role
    FROM SpecPlayerScores
    GROUP BY class, spec
)
SELECT 
    class,
    spec,
    LOWER(CONCAT(class, '-', REPLACE(spec, ' ', '-'))) AS slug,
    avg_global_score,
    player_count,
    role,
    RANK() OVER (ORDER BY avg_global_score DESC) AS overall_rank,
    RANK() OVER (PARTITION BY role ORDER BY avg_global_score DESC) AS role_rank
FROM SpecAverages
ORDER BY avg_global_score DESC;

-- Recreate top 10 players per spec view
CREATE VIEW top_10_players_per_spec AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        CAST(SUM(score) AS numeric(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL 
    AND server_region <> 'CN'
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

-- Recreate spec_dungeon_max_key_levels view 
CREATE VIEW spec_dungeon_max_key_levels AS
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