-- Migration 036 UP : Update views to exclude CN players, adding a new view and various adjustments

-- View 1 : spec_dungeon_max_key_levels (Updated)
-- This view calculates the maximum key level for each spec, excluding CN players
CREATE OR REPLACE VIEW spec_dungeon_max_key_levels AS
SELECT 
    pr.class,
    pr.spec,
    d.name AS dungeon_name,
    d.slug AS dungeon_slug,
    pr.dungeon_id AS encounter_id,
    max(pr.hard_mode_level) AS max_key_level
FROM player_rankings pr
JOIN dungeons d ON (pr.dungeon_id = d.encounter_id)
WHERE pr.deleted_at IS NULL
  AND pr.server_region != 'CN'
GROUP BY pr.class, pr.spec, d.name, d.slug, pr.dungeon_id
ORDER BY pr.class, pr.spec, max(pr.hard_mode_level) DESC;

-- View 2 : spec_global_score_averages (Updated)
-- This view calculates the average, max, min, and total score for each spec, excluding CN players
CREATE OR REPLACE VIEW spec_global_score_averages AS
WITH specplayerscores AS (
    SELECT 
        player_rankings.class,
        player_rankings.spec,
        player_rankings.name,
        player_rankings.server_name,
        player_rankings.server_region,
        (sum(player_rankings.score))::numeric(10,2) AS total_score,
        min((player_rankings.role)::text) AS role
    FROM player_rankings
    WHERE player_rankings.deleted_at IS NULL
      AND player_rankings.server_region != 'CN'
    GROUP BY player_rankings.class, player_rankings.spec, player_rankings.name, 
             player_rankings.server_name, player_rankings.server_region
    HAVING (count(DISTINCT player_rankings.dungeon_id) = 8)
), 
top10players AS (
    SELECT 
        class,
        spec,
        total_score,
        role,
        ROW_NUMBER() OVER (PARTITION BY class, spec ORDER BY total_score DESC) as rn
    FROM specplayerscores
),
specaverages AS (
    SELECT 
        class,
        spec,
        (avg(total_score))::numeric(10,2) AS avg_global_score,
        (max(total_score))::numeric(10,2) AS max_global_score,
        (min(total_score))::numeric(10,2) AS min_global_score,
        count(*) AS player_count,
        min(role) AS role
    FROM top10players
    WHERE rn <= 10
    GROUP BY class, spec
)
SELECT 
    specaverages.class,
    specaverages.spec,
    lower(concat(specaverages.class, '-', replace((specaverages.spec)::text, ' '::text, '-'::text))) AS slug,
    specaverages.avg_global_score,
    specaverages.max_global_score,
    specaverages.min_global_score,
    specaverages.player_count,
    specaverages.role,
    rank() OVER (ORDER BY specaverages.avg_global_score DESC) AS overall_rank,
    rank() OVER (PARTITION BY specaverages.role ORDER BY specaverages.avg_global_score DESC) AS role_rank
FROM specaverages
ORDER BY specaverages.avg_global_score DESC;

-- View 3 : spec_dungeon_score_averages (New view)
-- This view calculates the average, max, min, and total score for each spec for each dungeon, excluding CN players
CREATE OR REPLACE VIEW spec_dungeon_score_averages AS
WITH dungeon_player_scores AS (
    SELECT 
        pr.class,
        pr.spec,
        pr.dungeon_id,
        pr.name,
        pr.server_name,
        pr.server_region,
        pr.score,
        pr.role
    FROM player_rankings pr
    WHERE pr.deleted_at IS NULL
      AND pr.server_region != 'CN'
), 
top10_by_dungeon AS (
    SELECT 
        class,
        spec,
        dungeon_id,
        score,
        role,
        ROW_NUMBER() OVER (PARTITION BY class, spec, dungeon_id ORDER BY score DESC) as rn
    FROM dungeon_player_scores
),
spec_dungeon_averages AS (
    SELECT 
        class,
        spec,
        dungeon_id as encounter_id,
        (avg(score))::numeric(10,2) AS avg_dungeon_score,
        (max(score))::numeric(10,2) AS max_score,
        (min(score))::numeric(10,2) AS min_score,
        count(*) AS player_count,
        min(role) AS role
    FROM top10_by_dungeon
    WHERE rn <= 10
    GROUP BY class, spec, dungeon_id
),
ranked_results AS (
    SELECT 
        sda.class,
        sda.spec,
        lower(concat(sda.class, '-', replace((sda.spec)::text, ' '::text, '-'::text))) AS slug,
        sda.encounter_id,
        sda.avg_dungeon_score,
        sda.max_score,
        sda.min_score,
        sda.player_count,
        sda.role,
        rank() OVER (PARTITION BY sda.encounter_id ORDER BY sda.avg_dungeon_score DESC) AS overall_rank,
        rank() OVER (PARTITION BY sda.encounter_id, sda.role ORDER BY sda.avg_dungeon_score DESC) AS role_rank
    FROM spec_dungeon_averages sda
)
SELECT * FROM ranked_results
ORDER BY encounter_id, avg_dungeon_score DESC;