-- Migration 036 DOWN : Return to original views

-- View 1 : spec_dungeon_max_key_levels (return to original version)
CREATE OR REPLACE VIEW spec_dungeon_max_key_levels AS
SELECT pr.class,
    pr.spec,
    d.name AS dungeon_name,
    d.slug AS dungeon_slug,
    max(pr.hard_mode_level) AS max_key_level
FROM player_rankings pr
JOIN dungeons d ON (pr.dungeon_id = d.encounter_id)
WHERE pr.deleted_at IS NULL
GROUP BY pr.class, pr.spec, d.name, d.slug
ORDER BY pr.class, pr.spec, max(pr.hard_mode_level) DESC;

-- View 2 : spec_global_score_averages (return to original version)
CREATE OR REPLACE VIEW spec_global_score_averages AS
WITH specplayerscores AS (
    SELECT player_rankings.class,
        player_rankings.spec,
        player_rankings.name,
        player_rankings.server_name,
        player_rankings.server_region,
        (sum(player_rankings.score))::numeric(10,2) AS total_score,
        min((player_rankings.role)::text) AS role
    FROM player_rankings
    WHERE player_rankings.deleted_at IS NULL
    GROUP BY player_rankings.class, player_rankings.spec, player_rankings.name, player_rankings.server_name, player_rankings.server_region
    HAVING (count(DISTINCT player_rankings.dungeon_id) = 8)
), specaverages AS (
    SELECT specplayerscores.class,
        specplayerscores.spec,
        (avg(specplayerscores.total_score))::numeric(10,2) AS avg_global_score,
        count(*) AS player_count,
        min(specplayerscores.role) AS role
    FROM specplayerscores
    GROUP BY specplayerscores.class, specplayerscores.spec
)
SELECT specaverages.class,
    specaverages.spec,
    lower(concat(specaverages.class, '-', replace((specaverages.spec)::text, ' '::text, '-'::text))) AS slug,
    specaverages.avg_global_score,
    specaverages.player_count,
    specaverages.role,
    rank() OVER (ORDER BY specaverages.avg_global_score DESC) AS overall_rank,
    rank() OVER (PARTITION BY specaverages.role ORDER BY specaverages.avg_global_score DESC) AS role_rank
FROM specaverages
ORDER BY specaverages.avg_global_score DESC;

-- View 3 : spec_dungeon_score_averages (deleted because it didn't exist before)
DROP VIEW IF EXISTS spec_dungeon_score_averages;