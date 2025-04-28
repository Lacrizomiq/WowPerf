-- Migration UP: Add the workflow_states table and new columns to the existing tables

-- Creation of the workflow_states table
CREATE TABLE IF NOT EXISTS workflow_states (
    id VARCHAR(255) PRIMARY KEY, -- format: "workflow_type-run_id"
    workflow_type VARCHAR(50) NOT NULL, -- ex: "rankings", "reports", "builds"
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    items_processed INTEGER DEFAULT 0,
    status VARCHAR(20) NOT NULL, -- "running", "completed", "failed"
    error_message TEXT,
    last_processed_id VARCHAR(255), -- reference to the last processed element
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Add indexes for workflow_states
CREATE INDEX idx_workflow_states_type ON workflow_states(workflow_type);
CREATE INDEX idx_workflow_states_status ON workflow_states(status);
CREATE INDEX idx_workflow_states_type_created ON workflow_states(workflow_type, created_at);
CREATE INDEX idx_workflow_states_status_created ON workflow_states(status, created_at);
CREATE INDEX idx_workflow_states_duration ON workflow_states(started_at, completed_at);
CREATE INDEX idx_workflow_states_items_processed ON workflow_states(workflow_type, items_processed);

-- Add columns to class_rankings
ALTER TABLE class_rankings 
ADD COLUMN IF NOT EXISTS report_processing_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS report_processing_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS processing_batch_id VARCHAR(255);

-- Add indexes for class_rankings
CREATE INDEX idx_class_rankings_processing ON class_rankings(report_processing_status);
CREATE INDEX idx_rankings_batch ON class_rankings(processing_batch_id);
CREATE INDEX idx_rankings_processing_datetime ON class_rankings(report_processing_status, report_processing_at);

-- Add columns to warcraft_logs_reports
ALTER TABLE warcraft_logs_reports
ADD COLUMN IF NOT EXISTS build_extraction_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS build_extraction_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS processing_batch_id VARCHAR(255);

-- Add indexes for warcraft_logs_reports
CREATE INDEX idx_warcraft_logs_reports_extraction ON warcraft_logs_reports(build_extraction_status);
CREATE INDEX idx_reports_batch ON warcraft_logs_reports(processing_batch_id);
CREATE INDEX idx_reports_extraction_datetime ON warcraft_logs_reports(build_extraction_status, build_extraction_at);