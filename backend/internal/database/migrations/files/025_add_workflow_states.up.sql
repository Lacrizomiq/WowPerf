-- migration to add workflow_states table and columns to class_rankings and reports tables

-- Create the workflow_states table
CREATE TABLE workflow_states (
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

-- Add columns to the class_rankings table
ALTER TABLE class_rankings 
ADD COLUMN report_processing_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN report_processing_at TIMESTAMP,
ADD COLUMN processing_batch_id VARCHAR(255);

-- Add columns to the reports table
ALTER TABLE reports
ADD COLUMN build_extraction_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN build_extraction_at TIMESTAMP,
ADD COLUMN processing_batch_id VARCHAR(255);