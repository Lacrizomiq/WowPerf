-- Average Global Score per Spec
-- This is the spec average
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

-- Max Key Level per Spec for Each Dungeon
-- This finds the highest hard_mode_level for each spec in each dungeon, joining with the dungeons table to include name and slug. 
-- No 8-dungeon restriction here, as it’s spec/dungeon-specific
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

-- Average Global Score per Class
-- Mean score across all specs and dungeons per class, only for players with 8 dungeons
WITH ClassPlayerScores AS (
    SELECT 
        class,
        name,
        server_name,
        server_region,
        ROUND(SUM(score), 2) AS total_score
    FROM player_rankings
    WHERE deleted_at IS NULL
    GROUP BY class, name, server_name, server_region
    HAVING COUNT(DISTINCT dungeon_id) = 8
)
SELECT 
    class,
    ROUND(AVG(total_score), 2) AS avg_global_score,
    COUNT(*) AS player_count
FROM ClassPlayerScores
GROUP BY class
ORDER BY avg_global_score DESC;

-- Average Key Level per Dungeon with Dungeon Details
-- Average hard_mode_level per dungeon, joined with dungeons for name and slug:
SELECT 
    d.name AS dungeon_name,
    d.slug AS dungeon_slug,
    ROUND(AVG(pr.hard_mode_level), 2) AS avg_key_level,
    COUNT(*) AS run_count
FROM player_rankings pr
JOIN dungeons d ON pr.dungeon_id = d.encounter_id
WHERE pr.deleted_at IS NULL
GROUP BY d.name, d.slug
ORDER BY avg_key_level DESC;

-- Top 10 Players per Spec with Global Score
-- This finds the top 10 players for each spec, calculating their global score (summed across all 8 dungeons), ensuring they have 8 dungeon entries
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

-- Top 5 Players per Role with Global Score
-- Not representative that much but could be use leaderboard role based
WITH RolePlayerScores AS (
    SELECT 
        name,
        server_name,
        server_region,
        class,
        spec,
        role,
        ROUND(SUM(score), 2) AS total_score
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