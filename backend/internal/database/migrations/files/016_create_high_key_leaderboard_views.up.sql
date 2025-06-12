-- +migrate Up

-- Create view for average global score per spec (high-key players with 8 dungeons)
CREATE OR REPLACE VIEW spec_global_score_averages AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        CAST(SUM(score) AS NUMERIC(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, spec, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
)
SELECT 
    class,
    spec,
    CAST(AVG(total_score) AS NUMERIC(10,2)) AS avg_global_score,
    COUNT(*) AS player_count
FROM SpecPlayerScores
GROUP BY class, spec
ORDER BY avg_global_score DESC;

-- Create view for max key level per spec for each dungeon (high-key performance)
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

-- Create view for average global score per class (high-key players with 8 dungeons)
CREATE OR REPLACE VIEW class_global_score_averages AS
WITH ClassPlayerScores AS (
    SELECT 
        class,
        name,
        server_name,
        server_region,
        CAST(SUM(score) AS NUMERIC(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
)
SELECT 
    class,
    CAST(AVG(total_score) AS NUMERIC(10,2)) AS avg_global_score,
    COUNT(*) AS player_count
FROM ClassPlayerScores
GROUP BY class
ORDER BY avg_global_score DESC;

-- Create view for average key level per dungeon (high-key performance)
CREATE OR REPLACE VIEW dungeon_avg_key_levels AS
SELECT 
    d.name AS dungeon_name,
    d.slug AS dungeon_slug,
    CAST(AVG(pr.hard_mode_level) AS NUMERIC(10,2)) AS avg_key_level,
    COUNT(*) AS run_count
FROM player_rankings pr
JOIN dungeons d ON pr.dungeon_id = d.encounter_id
WHERE pr.deleted_at IS NULL
GROUP BY d.name, d.slug
ORDER BY avg_key_level DESC;

-- Create view for top 10 players per spec with global score (high-key players, excluding CN region)
CREATE OR REPLACE VIEW top_10_players_per_spec AS
WITH SpecPlayerScores AS (
    SELECT 
        class,
        spec,
        name,
        server_name,
        server_region,
        CAST(SUM(score) AS NUMERIC(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL AND server_region <> 'CN'
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

-- Create view for top 5 players per role with global score (high-key players)
CREATE OR REPLACE VIEW top_5_players_per_role AS
WITH RolePlayerScores AS (
    SELECT 
        name,
        server_name,
        server_region,
        class,
        spec,
        role,
        CAST(SUM(score) AS NUMERIC(10,2)) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
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