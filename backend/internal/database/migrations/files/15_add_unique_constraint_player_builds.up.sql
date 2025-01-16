-- +migrate Up

-- Add index for actor_id (needed for the unique constraint)
CREATE INDEX IF NOT EXISTS idx_player_builds_actor_id 
    ON player_builds(actor_id);

-- Add unique constraint for UPSERT operations support
ALTER TABLE player_builds
    ADD CONSTRAINT unique_player_builds_report_fight_actor 
    UNIQUE (report_code, fight_id, actor_id);