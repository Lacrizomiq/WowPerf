package warcraftlogsBuildsTemporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

type SyncWorkflow struct{}

func NewSyncWorkflow() *SyncWorkflow {
	return &SyncWorkflow{}
}

// Execute is the entry point for the sync workflow
func (w *SyncWorkflow) Execute(ctx workflow.Context, params WorkflowParams) (*WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &WorkflowResult{
		StartedAt: workflow.Now(ctx),
	}

	// Options for the activities with retry policy
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 1,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// 1. Process rankings
	var rankingsResult []*warcraftlogsBuilds.ClassRanking
	logger.Info("Starting rankings processing")

	for _, spec := range params.Specs {
		for _, dungeon := range params.Dungeons {
			var BatchResult *BatchResult
			err := workflow.ExecuteActivity(ctx, FetchRankingsActivityName, spec, dungeon, params.BatchConfig).Get(ctx, &BatchResult)
			if err != nil {
				logger.Error("Failed to process rankings",
					"spec", spec,
					"dungeon", dungeon,
					"error", err)
				continue
			}
			rankingsResult = append(rankingsResult, BatchResult.Rankings...)
			result.RankingsProcessed += len(BatchResult.Rankings)
		}
	}

	if len(rankingsResult) == 0 {
		logger.Info("No rankings found to process")
		result.CompletedAt = workflow.Now(ctx)
		return result, nil
	}

	// 2. Process reports
	var reportsResult *ReportProcessingResult
	logger.Info("Starting reports processing", "rankingsCount", len(rankingsResult))

	err := workflow.ExecuteActivity(ctx, ProcessReportsActivityName, rankingsResult).Get(ctx, &reportsResult)
	if err != nil {
		logger.Error("Failed to process reports", "error", err)
		result.CompletedAt = workflow.Now(ctx)
		return result, err
	}
	result.ReportsProcessed = reportsResult.ProcessedReports

	// Retrieve processed reports
	var reports []*warcraftlogsBuilds.Report
	err = workflow.ExecuteActivity(ctx, "get-processed-reports", rankingsResult).Get(ctx, &reports)
	if err != nil {
		logger.Error("Failed to retrieve processed reports", "error", err)
		result.CompletedAt = workflow.Now(ctx)
		return result, err
	}

	// 3. Process Player Builds
	var buildsResult *BuildsProcessingResult
	logger.Info("Starting player builds processing", "reportsCount", len(reports))

	err = workflow.ExecuteActivity(ctx, ProcessPlayerBuildsActivity, reports).Get(ctx, &buildsResult)
	if err != nil {
		logger.Error("Failed to process player builds", "error", err)
		result.CompletedAt = workflow.Now(ctx)
		return result, err
	}
	result.BuildsProcessed = buildsResult.ProcessedBuilds

	result.CompletedAt = workflow.Now(ctx)
	logger.Info("Workflow completed successfully",
		"rankingsProcessed", result.RankingsProcessed,
		"reportsProcessed", result.ReportsProcessed,
		"buildsProcessed", result.BuildsProcessed,
		"duration", result.CompletedAt.Sub(result.StartedAt),
	)

	return result, nil
}
