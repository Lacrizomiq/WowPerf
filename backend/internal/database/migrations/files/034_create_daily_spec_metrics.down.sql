-- 034_create_daily_spec_metrics.down.sql

-- Drop indexes for daily_spec_metrics_mythic_plus table
DROP INDEX IF EXISTS idx_daily_spec_metrics_mythic_plus_composite;
DROP INDEX IF EXISTS idx_daily_spec_metrics_mythic_plus_global;
DROP INDEX IF EXISTS idx_daily_spec_metrics_mythic_plus_spec_class_role;
DROP INDEX IF EXISTS idx_daily_spec_metrics_mythic_plus_date;

-- Drop daily_spec_metrics_mythic_plus table
DROP TABLE IF EXISTS daily_spec_metrics_mythic_plus;