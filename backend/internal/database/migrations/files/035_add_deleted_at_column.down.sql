-- 035_add_deleted_at_column.down.sql
-- Delete the index
DROP INDEX IF EXISTS idx_daily_spec_metrics_mythic_plus_deleted_at;

-- Delete the deleted_at column
ALTER TABLE daily_spec_metrics_mythic_plus DROP COLUMN IF EXISTS deleted_at;