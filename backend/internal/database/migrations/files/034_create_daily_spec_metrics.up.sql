-- 034_create_daily_spec_metrics.up.sql
-- This migration creates the daily_spec_metrics_mythic_plus table and indexes for better performance.
-- It also adds a unique constraint on the capture_date, spec, class, role, and encounter_id columns.

-- daily_spec_metrics_mythic_plus table
CREATE TABLE IF NOT EXISTS daily_spec_metrics_mythic_plus (
    id SERIAL PRIMARY KEY,
    capture_date DATE NOT NULL,
    spec VARCHAR(50) NOT NULL,
    class VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL,
    encounter_id INTEGER NOT NULL,
    is_global BOOLEAN NOT NULL DEFAULT FALSE,
    avg_score FLOAT NOT NULL,
    max_score FLOAT NOT NULL,
    min_score FLOAT NOT NULL,
    avg_key_level FLOAT NOT NULL,
    max_key_level INTEGER NOT NULL,
    min_key_level INTEGER NOT NULL,
    role_rank INTEGER NOT NULL,
    overall_rank INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(capture_date, spec, class, role, encounter_id)
);

-- Create indexes for better performance
CREATE INDEX idx_daily_spec_metrics_date ON daily_spec_metrics_mythic_plus(capture_date);
CREATE INDEX idx_daily_spec_metrics_spec_class_role ON daily_spec_metrics_mythic_plus(spec, class, role);
CREATE INDEX idx_daily_spec_metrics_global ON daily_spec_metrics_mythic_plus(is_global);
CREATE INDEX idx_daily_spec_metrics_composite ON daily_spec_metrics_mythic_plus(capture_date, spec, is_global, encounter_id);