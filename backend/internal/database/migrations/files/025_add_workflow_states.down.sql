-- migration to drop workflow_states table and columns from class_rankings and reports tables

-- Drop the workflow_states table
DROP TABLE workflow_states;

-- Drop columns from the reports table
ALTER TABLE reports 
DROP COLUMN build_extraction_status,
DROP COLUMN build_extraction_at,
DROP COLUMN processing_batch_id;

-- Drop columns from the class_rankings table
ALTER TABLE class_rankings
DROP COLUMN report_processing_status,
DROP COLUMN report_processing_at,
DROP COLUMN processing_batch_id;