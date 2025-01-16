package warcraftlogsBuildsTemporalWorkflows

import (
	"time"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// SyncWorkflow est maintenant une fonction au lieu d'une struct/méthode
func SyncWorkflow(ctx workflow.Context, params WorkflowParams) (*WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &WorkflowResult{
		StartedAt: workflow.Now(ctx),
	}

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 2,
		HeartbeatTimeout:    time.Minute * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Check the number of player builds
	var buildsCount int64
	if err := workflow.ExecuteActivity(ctx, CountPlayerBuildsActivityName).Get(ctx, &buildsCount); err != nil {
		logger.Error("Failed to count player builds", "error", err)
		return nil, err
	}

	// If no player builds, retrieve existing reports and generate builds
	if buildsCount == 0 {
		logger.Info("No player builds found, attempting to rebuild from existing reports")

		for _, dungeon := range params.Dungeons {
			var reports []warcraftlogsBuilds.Report
			err := workflow.ExecuteActivity(ctx, GetReportsForEncounterName, dungeon.EncounterID).Get(ctx, &reports)
			if err != nil {
				logger.Error("Failed to get reports for dungeon",
					"dungeonID", dungeon.EncounterID,
					"error", err)
				continue
			}

			if len(reports) > 0 {
				// Traiter par lots de 10 reports
				const batchSize = 10
				for i := 0; i < len(reports); i += batchSize {
					end := i + batchSize
					if end > len(reports) {
						end = len(reports)
					}

					batch := reports[i:end]
					reportPtrs := make([]*warcraftlogsBuilds.Report, len(batch))
					for j := range batch {
						reportPtrs[j] = &batch[j]
					}

					var buildsResult *BuildsProcessingResult
					if err := workflow.ExecuteActivity(ctx, ProcessPlayerBuildsActivity, reportPtrs).Get(ctx, &buildsResult); err != nil {
						logger.Error("Failed to process batch of reports",
							"startIndex", i,
							"endIndex", end,
							"error", err)
						continue
					}

					result.BuildsProcessed += buildsResult.ProcessedBuilds
				}
			}
		}
	}

	// 1. Rankings processing
	var rankingsResult []*BatchResult
	logger.Info("Starting rankings processing")

	for _, spec := range params.Specs {
		for _, dungeon := range params.Dungeons {
			var batchResult BatchResult
			err := workflow.ExecuteActivity(ctx, FetchRankingsActivityName,
				spec, dungeon, params.BatchConfig).Get(ctx, &batchResult)

			if err != nil {
				logger.Error("Failed to process rankings",
					"spec", spec,
					"dungeon", dungeon,
					"error", err)
				continue
			}

			rankingsResult = append(rankingsResult, &batchResult)
			result.RankingsProcessed += len(batchResult.Rankings)
		}
	}

	if len(rankingsResult) == 0 {
		logger.Info("No rankings found to process")
		result.CompletedAt = workflow.Now(ctx)
		return result, nil
	}

	// 2. Reports processing
	logger.Info("Starting reports processing", "rankingsCount", len(rankingsResult))

	var allRankings []*warcraftlogsBuilds.ClassRanking
	for _, batch := range rankingsResult {
		allRankings = append(allRankings, batch.Rankings...)
	}

	var reportsResult *ReportProcessingResult
	if err := workflow.ExecuteActivity(ctx, ProcessReportsActivityName, allRankings).Get(ctx, &reportsResult); err != nil {
		logger.Error("Failed to process reports", "error", err)
		return nil, err
	}
	result.ReportsProcessed = reportsResult.ProcessedReports

	var reports []*warcraftlogsBuilds.Report
	if err := workflow.ExecuteActivity(ctx, GetProcessedReportsActivityName, allRankings).Get(ctx, &reports); err != nil {
		logger.Error("Failed to retrieve processed reports", "error", err)
		return nil, err
	}

	// 3. Player builds processing
	logger.Info("Starting player builds processing", "reportsCount", len(reports))

	var buildsResult *BuildsProcessingResult
	if err := workflow.ExecuteActivity(ctx, ProcessPlayerBuildsActivity, reports).Get(ctx, &buildsResult); err != nil {
		logger.Error("Failed to process player builds", "error", err)
		return nil, err
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
