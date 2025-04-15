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

/*
	This file is not used anymore.
	All the logic is now in the models package.
	Some type might still be used in other files, but they should be removed
	or migrated to the models package.
*/

// ProcessState tracks the detailed progress of spec and dungeon processing
type ProcessState struct {
	CurrentSpec    ClassSpec
	CurrentDungeon Dungeon
	ProcessedCount int
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

func SyncWorkflow(ctx workflow.Context, params WorkflowParams) (*WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting sync workflow", "workflowID", params.WorkflowID)

	// Get the class name from workflow ID
	className := getClassFromWorkflowID(params.WorkflowID)
	if className == "" {
		return nil, fmt.Errorf("invalid workflow ID format: %s", params.WorkflowID)
	}

	// Initialize or recover progress
	progress := params.Progress
	if progress == nil {
		// Filter specs for just this class
		classSpecs := FilterSpecsForClass(params.Config.Specs, className)
		if len(classSpecs) == 0 {
			return nil, fmt.Errorf("no specs found for class: %s", className)
		}

		progress = &WorkflowProgress{
			CompletedSpecs:    make(map[string]bool),
			CompletedDungeons: make(map[string]bool),
			PartialResults: &WorkflowResult{
				StartedAt: workflow.Now(ctx),
			},
			Stats: &ProgressStats{
				TotalSpecs:    len(classSpecs),
				TotalDungeons: len(params.Config.Dungeons),
				StartedAt:     workflow.Now(ctx),
			},
		}
		logger.Info("Initialized new workflow progress",
			"class", className,
			"totalSpecs", progress.Stats.TotalSpecs,
			"totalDungeons", progress.Stats.TotalDungeons)
	} else {
		logger.Info("Resuming workflow progress",
			"class", className,
			"completedSpecs", len(progress.CompletedSpecs),
			"completedDungeons", len(progress.CompletedDungeons),
			"processedSpecs", progress.Stats.ProcessedSpecs,
			"processedDungeons", progress.Stats.ProcessedDungeons)
	}

	// Configure activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 2,
		HeartbeatTimeout:    time.Minute * 15,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Get filtered specs for this class
	classSpecs := FilterSpecsForClass(params.Config.Specs, className)

	// Phase 1 & 2: Sequential API Processing for Rankings and Reports
	logger.Info("Starting sequential API processing phase")
	for i, spec := range classSpecs {
		progress.Stats.CurrentSpec = &spec
		progress.Stats.LastSpecUpdate = workflow.Now(ctx)

		logger.Info("Processing spec",
			"specIndex", i,
			"totalSpecs", len(classSpecs),
			"className", spec.ClassName,
			"specName", spec.SpecName,
			"processedSpecs", progress.Stats.ProcessedSpecs)

		if progress.CompletedSpecs[spec.ClassName+spec.SpecName] {
			logger.Info("Skipping completed spec",
				"className", spec.ClassName,
				"specName", spec.SpecName)
			continue
		}

		for j := progress.CurrentDungeonIndex; j < len(params.Config.Dungeons); j++ {
			dungeon := params.Config.Dungeons[j]
			progress.Stats.CurrentDungeon = &dungeon
			progress.Stats.LastDungeonUpdate = workflow.Now(ctx)

			logger.Info("Processing dungeon",
				"dungeonIndex", j,
				"totalDungeons", progress.Stats.TotalDungeons,
				"dungeonName", dungeon.Name,
				"processedDungeons", progress.Stats.ProcessedDungeons)

			dungeonKey := generateDungeonKey(spec, dungeon)
			if progress.CompletedDungeons[dungeonKey] {
				logger.Info("Skipping completed dungeon",
					"dungeonName", dungeon.Name)
				continue
			}

			if err := processSpecAndDungeon(ctx, spec, dungeon, params, progress); err != nil {
				if isQuotaExceeded(err) {
					progress.CurrentSpecIndex = i
					progress.CurrentDungeonIndex = j
					params.Progress = progress
					return nil, workflow.NewContinueAsNewError(ctx, SyncWorkflowName, params)
				}
				logger.Error("Failed to process spec and dungeon",
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)
				continue
			}

			progress.CompletedDungeons[dungeonKey] = true
			progress.Stats.ProcessedDungeons++

			logger.Info("Completed dungeon processing",
				"dungeonName", dungeon.Name,
				"processedDungeons", progress.Stats.ProcessedDungeons,
				"remainingDungeons", progress.Stats.TotalDungeons-progress.Stats.ProcessedDungeons)

			// Add small delay between combinations
			workflow.Sleep(ctx, time.Second*2)
		}

		progress.CompletedSpecs[spec.ClassName+spec.SpecName] = true
		progress.Stats.ProcessedSpecs++
		progress.CurrentDungeonIndex = 0

		logger.Info("Completed spec processing",
			"className", spec.ClassName,
			"specName", spec.SpecName,
			"processedSpecs", progress.Stats.ProcessedSpecs,
			"remainingSpecs", len(classSpecs)-progress.Stats.ProcessedSpecs)
	}

	// Phase 3: Parallel Processing of Builds
	logger.Info("Starting parallel builds processing phase")
	rebuildResult, err := rebuildFromExistingReports(ctx, params.Config)
	if err != nil {
		logger.Error("Failed to process builds", "error", err)
		return nil, err
	}

	// Update final results
	progress.PartialResults.BuildsProcessed += rebuildResult.TotalBuildsProcessed
	progress.PartialResults.CompletedAt = workflow.Now(ctx)

	logger.Info("Workflow completed successfully",
		"class", className,
		"duration", progress.PartialResults.CompletedAt.Sub(progress.PartialResults.StartedAt),
		"processedSpecs", progress.Stats.ProcessedSpecs,
		"processedDungeons", progress.Stats.ProcessedDungeons,
		"totalBuildsProcessed", rebuildResult.TotalBuildsProcessed)

	return progress.PartialResults, nil
}

// rebuildFromExistingReports handles rebuilding player builds from stored reports
func rebuildFromExistingReports(ctx workflow.Context, cfg *Config) (*RebuildResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &RebuildResult{
		StartedAt: workflow.Now(ctx),
	}

	// Get total count of reports to process
	var reportsCount int64
	if err := workflow.ExecuteActivity(ctx,
		CountAllReportsActivity).Get(ctx, &reportsCount); err != nil {
		return nil, fmt.Errorf("failed to count reports: %w", err)
	}

	// Early return if no reports to process
	if reportsCount == 0 {
		logger.Info("No reports found for rebuild")
		result.CompletedAt = workflow.Now(ctx)
		return result, nil
	}

	logger.Info("Starting parallel rebuild from reports",
		"totalReports", reportsCount)

	// Process reports in optimized batches
	const batchSize = 50
	totalProcessed := 0
	successfulBatches := 0
	failedBatches := 0

	for offset := 0; offset < int(reportsCount); offset += batchSize {
		// Create batch parameters
		batchParams := BuildBatchParams{
			BatchSize:  batchSize,
			Offset:     offset,
			TotalCount: int(reportsCount),
		}

		// Process batch with child workflow
		var batchResult BuildBatchResult
		err := workflow.ExecuteChildWorkflow(ctx,
			ProcessBuildBatchWorkflowName,
			batchParams).Get(ctx, &batchResult)

		if err != nil {
			logger.Error("Failed to process batch",
				"offset", offset,
				"error", err)
			failedBatches++
			continue
		}

		// Update progress only if builds were processed
		if batchResult.ProcessedBuilds > 0 {
			totalProcessed += batchResult.ProcessedBuilds
			successfulBatches++

			// Record progress
			logger.Info("Batch processing progress",
				"batchProcessed", batchResult.ProcessedBuilds,
				"totalProcessed", totalProcessed,
				"progress", fmt.Sprintf("%d/%d", offset+batchSize, reportsCount),
				"successfulBatches", successfulBatches,
				"failedBatches", failedBatches)
		} else {
			// If no builds processed and not at the start, we might be done
			if offset > 0 {
				logger.Info("No more builds to process, ending rebuild",
					"totalProcessed", totalProcessed,
					"successfulBatches", successfulBatches)
				break
			}
		}

		// Small delay between batches to prevent database overload
		workflow.Sleep(ctx, time.Millisecond*500)
	}

	result.CompletedAt = workflow.Now(ctx)
	result.TotalBuildsProcessed = totalProcessed
	result.SuccessfulBatches = successfulBatches

	logger.Info("Completed rebuild",
		"totalBuildsProcessed", totalProcessed,
		"successfulBatches", successfulBatches,
		"failedBatches", failedBatches,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}

// ProcessBuildBatch processes a batch of reports and creates player builds from them
func ProcessBuildBatch(ctx workflow.Context, params BuildBatchParams) (*BuildBatchResult, error) {
	logger := workflow.GetLogger(ctx)
	result := &BuildBatchResult{
		StartedAt: workflow.Now(ctx),
	}

	// Configure activity options for batch processing
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 20,
		HeartbeatTimeout:    time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Fetch reports batch with improved error handling
	var reportsBatch []*warcraftlogsBuilds.Report
	err := workflow.ExecuteActivity(activityCtx,
		GetReportsBatchActivity,
		params.BatchSize,
		params.Offset).Get(ctx, &reportsBatch)

	if err != nil {
		logger.Error("Failed to fetch reports batch",
			"error", err,
			"offset", params.Offset,
			"batchSize", params.BatchSize)
		return nil, fmt.Errorf("failed to fetch reports batch: %w", err)
	}

	// Early return if no reports found
	if len(reportsBatch) == 0 {
		if params.Offset > 0 { // Only log if not the first batch
			logger.Info("No more reports to process - ending batch processing",
				"processedSoFar", params.Offset,
				"totalExpected", params.TotalCount)
		}
		result.CompletedAt = workflow.Now(ctx)
		return result, nil
	}

	// Process the batch of reports
	var buildsResult *BuildsProcessingResult
	if err := workflow.ExecuteActivity(activityCtx,
		ProcessBuildsActivity,
		reportsBatch).Get(ctx, &buildsResult); err != nil {
		logger.Error("Failed to process builds",
			"error", err,
			"reportsCount", len(reportsBatch))
		return nil, fmt.Errorf("failed to process builds: %w", err)
	}

	result.ProcessedBuilds = buildsResult.ProcessedBuilds
	result.CompletedAt = workflow.Now(ctx)

	logger.Info("Completed build batch processing",
		"processedBuilds", buildsResult.ProcessedBuilds,
		"batchProgress", fmt.Sprintf("%d/%d", params.Offset+params.BatchSize, params.TotalCount),
		"duration", result.CompletedAt.Sub(result.StartedAt))

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
// It handles only the API operations: Rankings -> Reports
func processSpecAndDungeon(ctx workflow.Context, spec ClassSpec, dungeon Dungeon, params WorkflowParams, progress *WorkflowProgress) error {
	logger := workflow.GetLogger(ctx)
	startTime := workflow.Now(ctx)

	// Validate spec matches workflow
	expectedClass := getClassFromWorkflowID(params.WorkflowID)
	if spec.ClassName != expectedClass {
		return fmt.Errorf("spec class %s does not match workflow class %s", spec.ClassName, expectedClass)
	}

	logger.Info("Processing spec and dungeon",
		"className", spec.ClassName,
		"spec", spec.SpecName,
		"dungeon", dungeon.Name,
		"startTime", startTime)

	// Step 1: Get rankings
	var batchResult BatchResult
	err := workflow.ExecuteActivity(ctx, FetchRankingsActivity,
		spec, dungeon, params.Config.Rankings.Batch).Get(ctx, &batchResult)

	if err != nil {
		logger.Error("Failed to fetch rankings",
			"spec", spec.SpecName,
			"dungeon", dungeon.Name,
			"error", err)
		return err
	}

	var rankingsToProcess []*warcraftlogsBuilds.ClassRanking

	if len(batchResult.Rankings) > 0 {
		rankingsToProcess = batchResult.Rankings
		progress.PartialResults.RankingsProcessed += len(batchResult.Rankings)
		logger.Info("Using newly fetched rankings", "count", len(batchResult.Rankings))
	} else {
		if err := workflow.ExecuteActivity(ctx, GetStoredRankingsActivity,
			spec.ClassName, spec.SpecName, dungeon.EncounterID).Get(ctx, &rankingsToProcess); err != nil {
			logger.Error("Failed to get stored rankings", "error", err)
			return err
		}
		logger.Info("Using stored rankings", "count", len(rankingsToProcess))
	}

	// Step 2: Process and sync reports if we have rankings
	if len(rankingsToProcess) > 0 {
		logger.Info("Starting reports processing",
			"spec", spec.SpecName,
			"dungeon", dungeon.Name,
			"rankingsCount", len(rankingsToProcess))

		var reportsResult *ReportProcessingResult
		if err := workflow.ExecuteActivity(ctx, ProcessReportsActivity,
			rankingsToProcess).Get(ctx, &reportsResult); err != nil {
			logger.Error("Failed to process reports",
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"error", err)
			return err
		}

		progress.PartialResults.ReportsProcessed += reportsResult.ProcessedCount

		logger.Info("Completed API processing phase",
			"spec", spec.SpecName,
			"dungeon", dungeon.Name,
			"rankingsProcessed", len(rankingsToProcess),
			"reportsProcessed", reportsResult.ProcessedCount,
			"duration", time.Since(startTime))
	} else {
		logger.Info("No rankings to process",
			"spec", spec.SpecName,
			"dungeon", dungeon.Name)
	}

	return nil
}

func estimateRequiredPoints(spec ClassSpec, dungeon Dungeon) float64 {
	basePoints := 1.0
	reportsPoints := 2.0
	estimatedReports := 20.0 // Average reports per spec/dungeon
	totalPoints := basePoints + (reportsPoints * estimatedReports)

	return totalPoints * 1.2 // 20% buffer
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

// generateDungeonKey generate a key for the completed dungeon so all spec can be processed
func generateDungeonKey(spec ClassSpec, dungeon Dungeon) string {
	return fmt.Sprintf("%s_%s_%d", spec.ClassName, spec.SpecName, dungeon.ID)
}

// Helper function to extract class name from workflow ID
func getClassFromWorkflowID(workflowID string) string {
	// Example workflowID: "warcraft-logs-Druid-workflow-2025-02-11T14:55:15Z"
	parts := strings.Split(workflowID, "-")
	if len(parts) >= 4 {
		return parts[2]
	}
	return ""
}

// Helper function to filter specs for a specific class
func FilterSpecsForClass(specs []ClassSpec, className string) []ClassSpec {
	filtered := make([]ClassSpec, 0)
	for _, spec := range specs {
		if spec.ClassName == className {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}
