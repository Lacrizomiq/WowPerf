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

// == Decoupled workflows ==

// RankingsWorkflow defines the interface for the rankings retrieval workflow
// It handles the orchestration of fetching and storing player rankings from Warcraft Logs.
type RankingsWorkflow interface {
	Execute(ctx workflow.Context, params models.RankingsWorkflowParams) (*models.RankingsWorkflowResult, error)
}

// ReportsWorkflow defines the interface for the reports processing workflow
// It orchestrates the retrieval and processing of detailed combat reports from Warcraft Logs.
type ReportsWorkflow interface {
	Execute(ctx workflow.Context, params models.ReportsWorkflowParams) (*models.ReportsWorkflowResult, error)
}

// BuildsWorkflow defines the interface for the builds extraction workflow
// It manages the extraction of player builds from combat reports and their storage.
type BuildsWorkflow interface {
	Execute(ctx workflow.Context, params models.BuildsWorkflowParams) (*models.BuildsWorkflowResult, error)
}

// EquipmentAnalysisWorkflow defines the interface for the equipment analysis workflow
// This workflow analyzes player equipment usage patterns for a specific spec and dungeon
type EquipmentAnalysisWorkflow interface {
	Execute(ctx workflow.Context, config models.EquipmentAnalysisWorkflowParams) (*models.EquipmentAnalysisWorkflowResult, error)
}

// TalentAnalysisWorkflow defines the interface for the talent analysis workflow
// This workflow analyzes player talent builds configuration for a specific spec and dungeon
type TalentAnalysisWorkflow interface {
	Execute(ctx workflow.Context, config models.TalentAnalysisWorkflowParams) (*models.TalentAnalysisWorkflowResult, error)
}

// StatAnalysisWorkflow defines the interface for the stats analysis workflow
// This workflow analyzes player stats distribution patterns for a specific spec and dungeon
type StatAnalysisWorkflow interface {
	Execute(ctx workflow.Context, config models.StatAnalysisWorkflowParams) (*models.StatAnalysisWorkflowResult, error)
}
