// result.go
package warcraftlogsBuildsTemporalWorkflowsModels

import (
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

/*

This file contains the models for the results of the workflows.

*/

// RankingsWorkflowResult holds the results of the rankings workflow
type RankingsWorkflowResult struct {
	RankingsProcessed int32     `json:"rankings_processed"` // Number of rankings processed
	SpecsProcessed    int32     `json:"specs_processed"`    // Number of specs processed
	DungeonsProcessed int32     `json:"dungeons_processed"` // Number of dungeons processed
	BatchID           string    `json:"batch_id"`           // Batch ID for tracking
	StartedAt         time.Time `json:"started_at"`         // Timestamp when the workflow started
	CompletedAt       time.Time `json:"completed_at"`       // Timestamp when the workflow completed
}

// ReportsWorkflowResult holds the results of the reports workflow
type ReportsWorkflowResult struct {
	ReportsProcessed  int32     `json:"reports_processed"`  // Number of reports processed
	RankingsProcessed int32     `json:"rankings_processed"` // Number of rankings processed
	FailedReports     int32     `json:"failed_reports"`     // Number of failed reports
	APIRequestsCount  int32     `json:"api_requests_count"` // Number of API requests
	BatchID           string    `json:"batch_id"`           // Batch ID for tracking
	StartedAt         time.Time `json:"started_at"`         // Timestamp when the workflow started
	CompletedAt       time.Time `json:"completed_at"`       // Timestamp when the workflow completed
}

// BuildsWorkflowResult holds the results of the builds workflow
type BuildsWorkflowResult struct {
	BuildsProcessed   int32            `json:"builds_processed"`     // Number of builds processed
	ReportsProcessed  int32            `json:"reports_processed"`    // Number of reports processed
	BuildsByClassSpec map[string]int32 `json:"builds_by_class_spec"` // Number of builds by class+spec
	BatchID           string           `json:"batch_id"`             // Batch ID for tracking
	StartedAt         time.Time        `json:"started_at"`           // Timestamp when the workflow started
	CompletedAt       time.Time        `json:"completed_at"`         // Timestamp when the workflow completed
}

// ReportProcessingResult holds the results of processing a batch of rankings for reports
type ReportProcessingResult struct {
	ProcessedCount   int32                        `json:"processed_count"`   // Number of reports processed in this batch
	SuccessCount     int32                        `json:"success_count"`     // Number of successful report processing
	FailureCount     int32                        `json:"failure_count"`     // Number of failed report processing
	ProcessedReports []*warcraftlogsBuilds.Report `json:"processed_reports"` // Reports processed in this batch
	ProcessedAt      time.Time                    `json:"processed_at"`      // Timestamp for this batch activity completion
}

// BuildsActivityResult holds the results of processing a batch of reports for builds
type BuildsActivityResult struct {
	ProcessedBuildsCount int32            `json:"processed_builds_count"` // Builds successfully stored in this batch
	SuccessCount         int32            `json:"success_count"`          // Reports successfully processed in this batch
	FailureCount         int32            `json:"failure_count"`          // Reports that failed processing in this batch
	BuildsByClassSpec    map[string]int32 `json:"builds_by_class_spec"`   // Builds counted by class+spec FOR THIS BATCH
	ProcessedAt          time.Time        `json:"processed_at"`           // Timestamp for this batch activity completion
}

// WorkflowResult represents the final outcome of a workflow execution
// This is the result of the main workflow, which is the synchronization of the rankings, reports and builds
// Legacy workflow, not used anymore, will be removed in the future.
type WorkflowResult struct {
	RankingsProcessed int32     `json:"rankings_processed"` // Number of rankings processed
	ReportsProcessed  int32     `json:"reports_processed"`  // Number of reports processed
	BuildsProcessed   int32     `json:"builds_processed"`   // Number of builds processed
	BatchID           string    `json:"batch_id"`           // Batch ID for tracking
	WorkflowType      string    `json:"workflow_type"`      // Workflow type (e.g., "rankings", "reports", "builds")
	StartedAt         time.Time `json:"started_at"`         // Timestamp when the workflow started
	CompletedAt       time.Time `json:"completed_at"`       // Timestamp when the workflow completed
}

// BuildsAnalysisResult represents the result of the builds analysis workflow
// It is used to analyze the builds and extract the best builds for each class and spec
// Should refacto the analysis activities cuz they should use the one in analysis_result.go
type BuildsAnalysisResult struct {
	BuildsProcessed int32     `json:"builds_processed"` // Number of builds processed
	ItemsProcessed  int32     `json:"items_processed"`  // Number of items processed
	ProcessedAt     time.Time `json:"processed_at"`     // Timestamp when the workflow completed
}
