-- +migrate Down

-- Remove unique constraint first
ALTER TABLE player_builds
    DROP CONSTRAINT IF EXISTS unique_player_builds_report_fight_actor;

-- Remove actor_id index
DROP INDEX IF EXISTS idx_player_builds_actor_id;