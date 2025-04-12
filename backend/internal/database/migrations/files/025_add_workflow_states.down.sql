-- Migration DOWN: Drop all changes made by the migration

-- Drop indexes for warcraft_logs_reports
DROP INDEX IF EXISTS idx_warcraft_logs_reports_extraction;
DROP INDEX IF EXISTS idx_reports_batch;
DROP INDEX IF EXISTS idx_reports_extraction_datetime;

-- Drop columns from warcraft_logs_reports
ALTER TABLE warcraft_logs_reports
DROP COLUMN IF EXISTS build_extraction_status,
DROP COLUMN IF EXISTS build_extraction_at,
DROP COLUMN IF EXISTS processing_batch_id;

-- Drop indexes for class_rankings
DROP INDEX IF EXISTS idx_class_rankings_processing;
DROP INDEX IF EXISTS idx_rankings_batch;
DROP INDEX IF EXISTS idx_rankings_processing_datetime;

-- Drop columns from class_rankings
ALTER TABLE class_rankings
DROP COLUMN IF EXISTS report_processing_status,
DROP COLUMN IF EXISTS report_processing_at,
DROP COLUMN IF EXISTS processing_batch_id;

-- Drop indexes for workflow_states
DROP INDEX IF EXISTS idx_workflow_states_type;
DROP INDEX IF EXISTS idx_workflow_states_status;
DROP INDEX IF EXISTS idx_workflow_states_type_created;
DROP INDEX IF EXISTS idx_workflow_states_status_created;
DROP INDEX IF EXISTS idx_workflow_states_duration;
DROP INDEX IF EXISTS idx_workflow_states_items_processed;

-- Drop the workflow_states table
DROP TABLE IF EXISTS workflow_states;