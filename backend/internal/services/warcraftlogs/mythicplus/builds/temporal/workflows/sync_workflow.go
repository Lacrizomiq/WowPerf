package warcraftlogsBuildsTemporalWorkflows

import (
	"errors"
	"fmt"
	"strings"
	"time"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ProcessState tracks the detailed progress of spec and dungeon processing
type ProcessState struct {
	CurrentSpec     ClassSpec
	CurrentDungeon  Dungeon
	RemainingPoints float64
	LastCheckTime   time.Time
	RetryCount      int
	ProcessedCount  int
}

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

	// Initialize or recover progress
	progress := params.Progress
	if progress == nil {
		progress = &WorkflowProgress{
			CompletedSpecs:    make(map[string]bool),
			CompletedDungeons: make(map[string]bool),
			PartialResults: &WorkflowResult{
				StartedAt: workflow.Now(ctx),
			},
		}
	}

	// Configure activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30,
		HeartbeatTimeout:    time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Initialize state tracking
	state := &ProcessState{
		LastCheckTime: workflow.Now(ctx),
	}

	// Initial points check
	if err := workflow.ExecuteActivity(ctx, CheckRemainingPointsActivity, params).Get(ctx, &state.RemainingPoints); err != nil {
		logger.Error("Failed to check points", "error", err)
		return nil, err
	}

	// Reserve points for workflow execution
	if err := workflow.ExecuteActivity(ctx, ReserveRateLimitPointsActivity, params).Get(ctx, nil); err != nil {
		logger.Error("Failed to reserve points", "error", err)
		if isQuotaExceeded(err) {
			workflow.Sleep(ctx, time.Minute*15)
			params.Progress = progress
			return nil, workflow.NewContinueAsNewError(ctx, "SyncWorkflow", params)
		}
		return nil, err
	}

	defer func() {
		if err := workflow.ExecuteActivity(ctx, ReleaseRateLimitPointsActivity, params).Get(ctx, nil); err != nil {
			logger.Error("Failed to release points", "error", err)
		}
	}()

	// Check if builds table is empty
	var buildsCount int64
	if err := workflow.ExecuteActivity(ctx, CountPlayerBuildsActivityName).Get(ctx, &buildsCount); err != nil {
		logger.Error("Failed to count player builds", "error", err)
		return nil, err
	}

	// Handle empty builds table rebuilding
	if buildsCount == 0 && progress.PartialResults.BuildsProcessed == 0 {
		logger.Info("No player builds found, attempting rebuild from existing reports")

		rebuildResult, err := rebuildFromExistingReports(ctx, params.Config)
		if err != nil {
			logger.Error("Failed to rebuild from existing reports", "error", err)
			return nil, err
		}

		if rebuildResult.TotalBuildsProcessed > 0 {
			progress.PartialResults.BuildsProcessed = rebuildResult.TotalBuildsProcessed
			progress.PartialResults.CompletedAt = workflow.Now(ctx)
			return progress.PartialResults, nil
		}
	}

	// Process specs and dungeons
	for i := progress.CurrentSpecIndex; i < len(params.Config.Specs); i++ {
		spec := params.Config.Specs[i]
		state.CurrentSpec = spec

		if progress.CompletedSpecs[spec.ClassName+spec.SpecName] {
			continue
		}

		for j := progress.CurrentDungeonIndex; j < len(params.Config.Dungeons); j++ {
			dungeon := params.Config.Dungeons[j]
			state.CurrentDungeon = dungeon

			dungeonKey := fmt.Sprintf("%d", dungeon.ID)
			if progress.CompletedDungeons[dungeonKey] {
				continue
			}

			// Check points if needed
			if time.Since(state.LastCheckTime) > time.Minute*5 || state.RemainingPoints < 5.0 {
				if err := workflow.ExecuteActivity(ctx, CheckRemainingPointsActivity, params).Get(ctx, &state.RemainingPoints); err != nil {
					logger.Error("Failed to check remaining points", "error", err)
					params.Progress = progress
					return nil, workflow.NewContinueAsNewError(ctx, "SyncWorkflow", params)
				}
				state.LastCheckTime = workflow.Now(ctx)
			}

			requiredPoints := estimateRequiredPoints(spec, dungeon)
			if state.RemainingPoints < requiredPoints {
				logger.Info("Insufficient points remaining, continuing as new workflow",
					"remaining", state.RemainingPoints,
					"required", requiredPoints)
				progress.CurrentSpecIndex = i
				progress.CurrentDungeonIndex = j
				params.Progress = progress
				return nil, workflow.NewContinueAsNewError(ctx, "SyncWorkflow", params)
			}

			if err := processSpecAndDungeon(ctx, spec, dungeon, params, progress); err != nil {
				logger.Error("Failed to process spec and dungeon",
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)
				if isQuotaExceeded(err) {
					progress.CurrentSpecIndex = i
					progress.CurrentDungeonIndex = j
					params.Progress = progress
					return nil, workflow.NewContinueAsNewError(ctx, "SyncWorkflow", params)
				}
				continue
			}

			progress.CompletedDungeons[dungeonKey] = true
			state.ProcessedCount++

			// Update remaining points estimate after processing
			state.RemainingPoints -= requiredPoints
		}

		progress.CompletedSpecs[spec.ClassName+spec.SpecName] = true
		progress.CurrentDungeonIndex = 0
	}

	progress.PartialResults.CompletedAt = workflow.Now(ctx)
	return progress.PartialResults, nil
}

// rebuildFromExistingReports handles the reconstruction of player builds from existing reports
func rebuildFromExistingReports(ctx workflow.Context, cfg *Config) (*RebuildResult, error) {
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
	for _, dungeon := range cfg.Dungeons {
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

// Helper functions

// newWorkflowProgress creates a new WorkflowProgress instance
func newWorkflowProgress(startTime time.Time) WorkflowProgress {
	return WorkflowProgress{
		CompletedSpecs:    make(map[string]bool),
		CompletedDungeons: make(map[string]bool),
		PartialResults: &WorkflowResult{
			StartedAt: startTime,
		},
	}
}

// processSpecAndDungeon processes a single spec and dungeon combination
func processSpecAndDungeon(ctx workflow.Context,
	spec ClassSpec,
	dungeon Dungeon,
	params WorkflowParams,
	progress *WorkflowProgress) error {

	logger := workflow.GetLogger(ctx)
	startTime := workflow.Now(ctx)

	logger.Info("Processing ranking combination",
		"spec", spec.SpecName,
		"dungeon", dungeon.Name,
		"startTime", startTime)

	if progress.CurrentSpecIndex > 0 {
		workflow.Sleep(ctx, time.Second*2)
	}

	var batchResult BatchResult
	err := workflow.ExecuteActivity(ctx, FetchRankingsActivityName,
		spec, dungeon, params.Config.Rankings.Batch).Get(ctx, &batchResult)

	if err != nil {
		logger.Error("Failed to fetch rankings",
			"spec", spec.SpecName,
			"dungeon", dungeon.Name,
			"error", err)
		return err
	}

	progress.PartialResults.RankingsProcessed += len(batchResult.Rankings)

	if len(batchResult.Rankings) > 0 {
		var reportsResult *ReportProcessingResult
		if err := workflow.ExecuteActivity(ctx, ProcessReportsActivityName,
			batchResult.Rankings).Get(ctx, &reportsResult); err != nil {
			logger.Error("Failed to process reports",
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"error", err)
			return err
		}

		progress.PartialResults.ReportsProcessed += reportsResult.ProcessedReports
	}

	return nil
}

// estimateRequiredPoints estimates the number of points needed for a spec and dungeon
func estimateRequiredPoints(spec ClassSpec, dungeon Dungeon) float64 {
	basePoints := 1.0
	reportsPoints := 2.0
	estimatedReports := 20.0
	totalPoints := basePoints + (reportsPoints * estimatedReports)

	// Buffer increased to 20% for more safety
	return totalPoints * 1.2
}

// Error implements the error interface
func (e *QuotaExceededError) Error() string {
	return e.Message
}

func isQuotaExceeded(err error) bool {
	if err == nil {
		return false
	}

	var quotaErr *QuotaExceededError
	return errors.As(err, &quotaErr) || strings.Contains(err.Error(), "quota exceeded")
}

func countPointsUsed(rankings []*warcraftlogsBuilds.ClassRanking) int {
	if len(rankings) == 0 {
		return 1 // Base cost for rankings request
	}
	// Base cost + (2 points per ranking for report processing)
	return 1 + (len(rankings) * 2)
}

// checkPointsAndWait checks if we have enough points and waits if not
func checkPointsAndWait(ctx workflow.Context,
	state *ProcessState,
	params WorkflowParams) error {

	logger := workflow.GetLogger(ctx)

	if time.Since(state.LastCheckTime) > time.Minute*5 || state.RemainingPoints < 5.0 {
		if err := workflow.ExecuteActivity(ctx, CheckRemainingPointsActivity, params).Get(ctx, &state.RemainingPoints); err != nil {
			return err
		}
		state.LastCheckTime = workflow.Now(ctx)

		if state.RemainingPoints < 1.0 {
			return &QuotaExceededError{
				Message: "Insufficient points available",
				ResetIn: time.Minute * 15,
			}
		}
	}

	logger.Info("Points available for requests",
		"remainingPoints", state.RemainingPoints)

	return nil
}
