package warcraftlogsBuilds

import (
	"time"

	"gorm.io/datatypes"
)

// WorkflowState represents the state of a workflow
// It is used to track metrics and status of the workflow
type WorkflowState struct {
	ID           string `gorm:"primaryKey"`
	WorkflowType string `gorm:"column:workflow_type"`

	// basic tracking
	StartedAt       time.Time
	CompletedAt     time.Time
	Status          string
	ErrorMessage    string `gorm:"column:error_message"`
	LastProcessedID string `gorm:"column:last_processed_id"`
	ItemsProcessed  int    `gorm:"default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// advanced tracking
	ParentWorkflowID    string    `gorm:"column:parent_workflow_id"`
	ContinuationCount   int       `gorm:"column:continuation_count;default:0"`
	TotalItemsToProcess int       `gorm:"column:total_items_to_process;default:0"`
	ProgressPercentage  float64   `gorm:"column:progress_percentage;default:0"`
	EstimatedCompletion time.Time `gorm:"column:estimated_completion"`

	// Specific fields for workflow types
	BatchID          string `gorm:"column:batch_id"`
	ClassName        string `gorm:"column:class_name"`
	ApiRequestsCount int    `gorm:"column:api_requests_count;default:0"`

	// Detailed metrics in JSON format
	PerformanceMetrics datatypes.JSON `gorm:"column:performance_metrics"`
}

func (WorkflowState) TableName() string {
	return "workflow_states"
}
