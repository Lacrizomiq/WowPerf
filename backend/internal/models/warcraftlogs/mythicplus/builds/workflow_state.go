package warcraftlogsBuilds

import "time"

// WorkflowState represents the state of a workflow
// It is used to track the progress of a workflow and to resume it in case of failure
// It is also useful to get information about a specific workflow for the other workflows
type WorkflowState struct {
	ID              string `gorm:"primaryKey"`
	WorkflowType    string `gorm:"column:workflow_type"`
	StartedAt       time.Time
	CompletedAt     time.Time
	ItemsProcessed  int `gorm:"default:0"`
	Status          string
	ErrorMessage    string `gorm:"column:error_message"`
	LastProcessedID string `gorm:"column:last_processed_id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (WorkflowState) TableName() string {
	return "workflow_states"
}
