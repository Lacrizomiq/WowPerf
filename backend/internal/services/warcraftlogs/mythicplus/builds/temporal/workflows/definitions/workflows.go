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

// RankingsWorkflow defines the interface for rankings synchronization workflow
type RankingsWorkflow interface {
	Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error)
}

// ProcessBuildBatchWorkflow defines the interface for build batch processing workflow
type ProcessBuildBatchWorkflow interface {
	Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error)
}

// SyncWorkflow defines the interface for the main synchronization workflow
type SyncWorkflow interface {
	Execute(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error)
}
