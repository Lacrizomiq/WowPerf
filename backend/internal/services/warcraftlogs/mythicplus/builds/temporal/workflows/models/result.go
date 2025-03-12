// result.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// WorkflowResult represents the final outcome of a workflow execution
type WorkflowResult struct {
	RankingsProcessed int32     `json:"rankings_processed"`
	ReportsProcessed  int32     `json:"reports_processed"`
	BuildsProcessed   int32     `json:"builds_processed"`
	StartedAt         time.Time `json:"started_at"`
	CompletedAt       time.Time `json:"completed_at"`
}
