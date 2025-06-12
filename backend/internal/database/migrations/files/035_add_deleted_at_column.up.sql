-- 035_add_deleted_at_column.up.sql

-- Add the deleted_at column to support GORM's soft delete
ALTER TABLE daily_spec_metrics_mythic_plus ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- Add an index on deleted_at to optimize queries filtering by this field
CREATE INDEX idx_daily_spec_metrics_mythic_plus_deleted_at ON daily_spec_metrics_mythic_plus(deleted_at);