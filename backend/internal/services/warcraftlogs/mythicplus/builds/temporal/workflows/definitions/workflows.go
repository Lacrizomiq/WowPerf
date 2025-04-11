// workflows.go
package warcraftlogsBuildsTemporalWorkflowsDefinitions

import (
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	"go.temporal.io/sdk/workflow"
)

/* This file :

- Define all the interface of the workflow
- Specify all contract a workflow should implement

*/

// RankingsWorkflow defines the interface for the rankings retrieval workflow
type RankingsWorkflow interface {
	Execute(ctx workflow.Context, params models.RankingsWorkflowParams) (*models.RankingsWorkflowResult, error)
}

// ReportsWorkflow defines the interface for the reports processing workflow
type ReportsWorkflow interface {
	Execute(ctx workflow.Context, params models.ReportsWorkflowParams) (*models.ReportsWorkflowResult, error)
}

// BuildsWorkflow defines the interface for the builds extraction workflow
type BuildsWorkflow interface {
	Execute(ctx workflow.Context, params models.BuildsWorkflowParams) (*models.BuildsWorkflowResult, error)
}

// SyncWorkflow defines the interface for the main synchronization workflow
// This workflow is used to synchronize the rankings, reports and builds from WarcraftLogs
// Legacy workflow, not used anymore, will be removed in the future.
type SyncWorkflow interface {
	Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error)
}

// AnalyzeWorkflow defines the interface for the main analysis workflow
type AnalyzeWorkflow interface {
	Execute(ctx workflow.Context, config models.AnalysisWorkflowConfig) (*models.AnalysisWorkflowResult, error)
}
