-- Migration UP: Add new tracking and specific fields to workflow_states table

-- Add new columns to workflow_states table
ALTER TABLE workflow_states
ADD COLUMN IF NOT EXISTS parent_workflow_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS continuation_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS total_items_to_process INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS progress_percentage FLOAT DEFAULT 0,
ADD COLUMN IF NOT EXISTS estimated_completion TIMESTAMP,
ADD COLUMN IF NOT EXISTS batch_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS class_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS api_requests_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS performance_metrics JSONB;

-- Add indexes for new columns
CREATE INDEX IF NOT EXISTS idx_workflow_states_parent_id ON workflow_states(parent_workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_states_batch_id ON workflow_states(batch_id);
CREATE INDEX IF NOT EXISTS idx_workflow_states_class_name ON workflow_states(class_name);
CREATE INDEX IF NOT EXISTS idx_workflow_states_type_status_progress ON workflow_states(workflow_type, status, progress_percentage);
