-- Migration DOWN: Remove new tracking and specific fields from workflow_states table

ALTER TABLE workflow_states
DROP COLUMN IF EXISTS parent_workflow_id,
DROP COLUMN IF EXISTS continuation_count,
DROP COLUMN IF EXISTS total_items_to_process,
DROP COLUMN IF EXISTS progress_percentage,
DROP COLUMN IF EXISTS estimated_completion,
DROP COLUMN IF EXISTS batch_id,
DROP COLUMN IF EXISTS class_name,
DROP COLUMN IF EXISTS api_requests_count,
DROP COLUMN IF EXISTS performance_metrics;

-- Drop indexes for removed columns (PostgreSQL drops indexes automatically when columns are dropped, but good practice to be explicit if needed for other DBs or specific scenarios)
-- For PostgreSQL, these DROP INDEX statements might be redundant if the columns are dropped.
-- However, if columns were not dropped but indexes needed removal, they would be here.
-- DROP INDEX IF EXISTS idx_workflow_states_parent_id;
-- DROP INDEX IF EXISTS idx_workflow_states_batch_id;
-- DROP INDEX IF EXISTS idx_workflow_states_class_name;
-- DROP INDEX IF EXISTS idx_workflow_states_type_status_progress;