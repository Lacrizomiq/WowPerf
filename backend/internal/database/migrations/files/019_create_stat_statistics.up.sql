-- Migration: 19_create_stat_statistics.up.sql
CREATE TABLE stat_statistics (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP,
  
  class VARCHAR(255) NOT NULL,
  spec VARCHAR(255) NOT NULL,
  encounter_id INTEGER NOT NULL,
  
  stat_name VARCHAR(50) NOT NULL,
  stat_category VARCHAR(50) NOT NULL,
  avg_value FLOAT NOT NULL DEFAULT 0,
  median_value FLOAT NOT NULL DEFAULT 0,
  min_value FLOAT NOT NULL DEFAULT 0,
  max_value FLOAT NOT NULL DEFAULT 0
);

CREATE INDEX idx_stat_statistics_class_spec ON stat_statistics(class, spec);
CREATE INDEX idx_stat_statistics_encounter_id ON stat_statistics(encounter_id);
CREATE INDEX idx_stat_statistics_category ON stat_statistics(stat_category);
CREATE INDEX idx_stat_statistics_deleted_at ON stat_statistics(deleted_at);