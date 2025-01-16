package warcraftlogsBuildsTemporalWorkflows

import (
	"fmt"
	"time"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// RebuildResult represents the result of rebuilding player builds from existing reports
type RebuildResult struct {
	TotalBuildsProcessed int
	SuccessfulBatches    int
	StartedAt            time.Time
	CompletedAt          time.Time
}

// BuildBatchParams contains parameters for processing a batch of builds
type BuildBatchParams struct {
	DungeonID  uint
	BatchSize  int
	Offset     int
	TotalCount int
}

// BuildBatchResult contains the result of processing a batch of builds
type BuildBatchResult struct {
	ProcessedBuilds int
	StartedAt       time.Time
	CompletedAt     time.Time
}

// SyncWorkflow is the main workflow for syncing builds
// It handles both initial data population and regular updates
func SyncWorkflow(ctx workflow.Context, params WorkflowParams) (*WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &WorkflowResult{
		StartedAt: workflow.Now(ctx),
	}

	// Configure activity options with retry policy
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

	// 1. Check if player_builds table is empty
	var buildsCount int64
	if err := workflow.ExecuteActivity(ctx, CountPlayerBuildsActivityName).Get(ctx, &buildsCount); err != nil {
		logger.Error("Failed to count player builds", "error", err)
		return nil, err
	}

	// 2. If table is empty, attempt reconstruction from existing reports
	if buildsCount == 0 {
		logger.Info("No player builds found, attempting rebuild from existing reports")

		rebuildResult, err := rebuildFromExistingReports(ctx, params.Dungeons)
		if err != nil {
			logger.Error("Failed to rebuild from existing reports", "error", err)
			return nil, err
		}

		result.BuildsProcessed = rebuildResult.TotalBuildsProcessed

		// If rebuild was successful, we can finish here
		if rebuildResult.TotalBuildsProcessed > 0 {
			result.CompletedAt = workflow.Now(ctx)
			logger.Info("Successfully rebuilt player builds",
				"totalBuilds", rebuildResult.TotalBuildsProcessed,
				"duration", result.CompletedAt.Sub(result.StartedAt))
			return result, nil
		}
	}

	// 3. Proceed with normal workflow (either table wasn't empty or rebuild produced no results)
	logger.Info("Starting normal workflow process")

	// Fetch and process rankings for each spec and dungeon
	var rankingsResult []*BatchResult
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

	// Process reports from rankings
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

	// Process builds for new reports
	logger.Info("Processing builds for new reports")

	var reports []*warcraftlogsBuilds.Report
	if err := workflow.ExecuteActivity(ctx, GetProcessedReportsActivityName, allRankings).Get(ctx, &reports); err != nil {
		logger.Error("Failed to retrieve processed reports", "error", err)
		return nil, err
	}

	// Process new reports in batches
	const newReportsBatchSize = 5
	for i := 0; i < len(reports); i += newReportsBatchSize {
		end := i + newReportsBatchSize
		if end > len(reports) {
			end = len(reports)
		}

		batchReports := reports[i:end]
		var buildsResult *BuildsProcessingResult
		if err := workflow.ExecuteActivity(ctx, ProcessPlayerBuildsActivity, batchReports).Get(ctx, &buildsResult); err != nil {
			logger.Error("Failed to process new reports builds",
				"startIndex", i,
				"endIndex", end,
				"error", err)
			continue
		}

		if buildsResult != nil {
			result.BuildsProcessed += buildsResult.ProcessedBuilds
		}

		// Add small delay between batches
		workflow.Sleep(ctx, time.Second*2)
	}

	result.CompletedAt = workflow.Now(ctx)
	logger.Info("Workflow completed successfully",
		"rankingsProcessed", result.RankingsProcessed,
		"reportsProcessed", result.ReportsProcessed,
		"buildsProcessed", result.BuildsProcessed,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}

// rebuildFromExistingReports handles the reconstruction of player builds from existing reports
func rebuildFromExistingReports(ctx workflow.Context, dungeons []Dungeon) (*RebuildResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &RebuildResult{
		StartedAt: workflow.Now(ctx),
	}

	// Configure activity options specifically for rebuild process
	rebuildActivityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
		HeartbeatTimeout:    time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
	}
	rebuildCtx := workflow.WithActivityOptions(ctx, rebuildActivityOpts)

	// Process each dungeon sequentially
	for _, dungeon := range dungeons {
		logger.Info("Processing dungeon for builds reconstruction",
			"dungeonName", dungeon.Name,
			"encounterID", dungeon.EncounterID)

		// Get total report count
		var reportsCount int64
		if err := workflow.ExecuteActivity(rebuildCtx,
			CountReportsForEncounterName,
			dungeon.EncounterID).Get(rebuildCtx, &reportsCount); err != nil {
			logger.Error("Failed to count reports for dungeon",
				"dungeonID", dungeon.EncounterID,
				"error", err)
			continue
		}

		if reportsCount == 0 {
			logger.Info("No reports found for dungeon",
				"dungeonName", dungeon.Name)
			continue
		}

		// Process reports in batches
		const batchSize = 5
		totalProcessed := 0

		for offset := 0; offset < int(reportsCount); offset += batchSize {
			// Process batch in child workflow
			batchParams := BuildBatchParams{
				DungeonID:  dungeon.EncounterID,
				BatchSize:  batchSize,
				Offset:     offset,
				TotalCount: int(reportsCount),
			}

			var batchResult BuildBatchResult
			err := workflow.ExecuteChildWorkflow(ctx,
				ProcessBuildBatch,
				batchParams).Get(ctx, &batchResult)

			if err != nil {
				logger.Error("Failed to process batch",
					"dungeonName", dungeon.Name,
					"offset", offset,
					"error", err)
				continue
			}

			totalProcessed += batchResult.ProcessedBuilds
			result.TotalBuildsProcessed += batchResult.ProcessedBuilds
			result.SuccessfulBatches++

			logger.Info("Processed batch of builds",
				"dungeonName", dungeon.Name,
				"batchProcessed", batchResult.ProcessedBuilds,
				"totalProcessed", totalProcessed,
				"progress", fmt.Sprintf("%d/%d", offset+batchSize, reportsCount))

			// Add delay between batches
			workflow.Sleep(ctx, time.Second*2)
		}
	}

	result.CompletedAt = workflow.Now(ctx)
	return result, nil
}

// ProcessBuildBatch is a child workflow that processes a single batch of builds
func ProcessBuildBatch(ctx workflow.Context, params BuildBatchParams) (*BuildBatchResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &BuildBatchResult{
		StartedAt: workflow.Now(ctx),
	}

	// Activity options for batch processing
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
		HeartbeatTimeout:    time.Minute * 2,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Fetch reports batch
	var reportsBatch []warcraftlogsBuilds.Report
	err := workflow.ExecuteActivity(activityCtx,
		GetReportsForEncounterBatchName,
		params.DungeonID,
		params.BatchSize,
		params.Offset).Get(ctx, &reportsBatch)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch reports batch: %w", err)
	}

	// Convert reports to pointers
	reportPtrs := make([]*warcraftlogsBuilds.Report, len(reportsBatch))
	for i := range reportsBatch {
		reportPtrs[i] = &reportsBatch[i]
	}

	// Process builds
	var buildsResult *BuildsProcessingResult
	if err := workflow.ExecuteActivity(activityCtx,
		ProcessPlayerBuildsActivity,
		reportPtrs).Get(ctx, &buildsResult); err != nil {
		return nil, fmt.Errorf("failed to process builds: %w", err)
	}

	logger.Info("Processed builds", "processedBuilds", buildsResult.ProcessedBuilds)

	result.ProcessedBuilds = buildsResult.ProcessedBuilds
	result.CompletedAt = workflow.Now(ctx)
	return result, nil
}
