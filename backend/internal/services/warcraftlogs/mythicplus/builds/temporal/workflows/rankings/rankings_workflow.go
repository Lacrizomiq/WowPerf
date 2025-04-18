package warcraftlogsBuildsTemporalWorkflowsRankings

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// RankingsWorkflow implements the rankings workflow
type RankingsWorkflow struct{}

// NewRankingsWorkflow creates a new instance of the rankings workflow
func NewRankingsWorkflow() definitions.RankingsWorkflow {
	return &RankingsWorkflow{}
}

// Execute runs the rankings workflow
func (w *RankingsWorkflow) Execute(ctx workflow.Context, params models.RankingsWorkflowParams) (*models.RankingsWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting rankings workflow",
		"batchID", params.BatchID,
		"specCount", len(params.Specs),
		"dungeonCount", len(params.Dungeons))

	// Initialize result with start time and batch ID
	result := &models.RankingsWorkflowResult{
		StartedAt: workflow.Now(ctx),
		BatchID:   params.BatchID,
	}

	// Validate parameters
	if len(params.Specs) == 0 {
		return nil, fmt.Errorf("no specs found in parameters")
	}

	if len(params.Dungeons) == 0 {
		return nil, fmt.Errorf("no dungeons found in parameters")
	}

	// WorkflowState Tracking
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	workflowStateID := fmt.Sprintf("rankings-%s", workflowID)

	// Activity options for workflow state
	stateOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.5,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	stateCtx := workflow.WithActivityOptions(ctx, stateOpts)

	// Create workflow state
	err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, &warcraftlogsBuilds.WorkflowState{
		ID:              workflowStateID,
		WorkflowType:    "rankings",
		StartedAt:       workflow.Now(ctx),
		Status:          "running",
		ItemsProcessed:  0,
		LastProcessedID: "",
		CreatedAt:       workflow.Now(ctx),
		UpdatedAt:       workflow.Now(ctx),
	}).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to create workflow state", "error", err)
		// Continue execution even if state tracking fails
	}

	// Track processed specs and dungeons
	processedSpecs := make(map[string]bool)
	processedDungeons := make(map[string]bool)

	// Configure activity options for the rankings fetch
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 12,
		HeartbeatTimeout:    time.Minute * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * time.Duration(params.RetryDelay.Seconds()),
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    int32(params.MaxAttempts),
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Process each spec
	for _, spec := range params.Specs {
		specKey := common.GenerateSpecKey(spec)
		if processedSpecs[specKey] {
			logger.Info("Skipping already processed spec",
				"class", spec.ClassName,
				"spec", spec.SpecName)
			continue
		}

		logger.Info("Processing spec",
			"class", spec.ClassName,
			"spec", spec.SpecName)

		// Update workflow state with current spec
		err = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, &warcraftlogsBuilds.WorkflowState{
			ID:              workflowStateID,
			Status:          "running",
			LastProcessedID: specKey,
			UpdatedAt:       workflow.Now(ctx),
		}).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to update workflow state", "error", err)
		}

		// Process each dungeon for this spec
		for _, dungeon := range params.Dungeons {
			dungeonKey := common.GenerateDungeonKey(spec, dungeon)
			if processedDungeons[dungeonKey] {
				logger.Info("Skipping already processed dungeon",
					"dungeon", dungeon.Name,
					"spec", spec.SpecName)
				continue
			}

			logger.Info("Processing dungeon",
				"dungeon", dungeon.Name,
				"class", spec.ClassName,
				"spec", spec.SpecName)

			// Create batch config from params
			batchConfig := models.BatchConfig{
				Size:        params.BatchSize,
				RetryDelay:  params.RetryDelay,
				MaxAttempts: params.MaxAttempts,
			}

			// Execute the activity to fetch and store rankings
			var batchResult models.BatchResult
			err := workflow.ExecuteActivity(activityCtx,
				definitions.FetchRankingsActivity,
				spec, dungeon, batchConfig).Get(ctx, &batchResult)

			if err != nil {
				if common.IsRateLimitError(err) {
					// Update workflow state for rate limit
					workflowState := &warcraftlogsBuilds.WorkflowState{
						ID:           workflowStateID,
						Status:       "rate_limited",
						ErrorMessage: fmt.Sprintf("Rate limit reached: %v", err),
						UpdatedAt:    workflow.Now(ctx),
					}
					_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

					logger.Info("Rate limit reached during rankings processing",
						"class", spec.ClassName,
						"spec", spec.SpecName,
						"dungeon", dungeon.Name)

					result.CompletedAt = workflow.Now(ctx)
					return result, err
				}

				logger.Error("Failed to process rankings",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)

				// Update workflow state with error
				workflowState := &warcraftlogsBuilds.WorkflowState{
					ID:           workflowStateID,
					ErrorMessage: fmt.Sprintf("Error processing %s: %v", dungeonKey, err),
					UpdatedAt:    workflow.Now(ctx),
				}
				_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

				// Continue with next dungeon on error
				continue
			}

			// Mark rankings with batch ID and status
			if batchResult.RankingsCount > 0 {
				// This activity would be a new one created to support the workflow
				err = workflow.ExecuteActivity(activityCtx, definitions.MarkRankingsForReportActivity,
					batchResult.ClassName, batchResult.SpecName, batchResult.EncounterID, params.BatchID).Get(ctx, nil)

				if err != nil {
					logger.Error("Failed to mark rankings for report processing",
						"class", spec.ClassName,
						"spec", spec.SpecName,
						"error", err)
				}
			}

			// Mark dungeon as processed
			processedDungeons[dungeonKey] = true

			// Update result with data from this batch
			result.RankingsProcessed += batchResult.ProcessedItems

			logger.Info("Successfully processed rankings",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"itemsProcessed", batchResult.ProcessedItems,
				"totalProcessed", result.RankingsProcessed)

			// Update workflow state with progress
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:             workflowStateID,
				ItemsProcessed: int(result.RankingsProcessed),
				UpdatedAt:      workflow.Now(ctx),
			}
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			// Small delay between dungeons to avoid overwhelming the API
			workflow.Sleep(ctx, time.Second*2)
		}

		// Mark spec as processed
		processedSpecs[specKey] = true
		result.SpecsProcessed++

		logger.Info("Completed processing for spec",
			"class", spec.ClassName,
			"spec", spec.SpecName)
	}

	// Set final counts and timestamps
	result.DungeonsProcessed = int32(len(processedDungeons))
	result.CompletedAt = workflow.Now(ctx)

	// Complete workflow state
	workflowState := &warcraftlogsBuilds.WorkflowState{
		ID:             workflowStateID,
		Status:         "completed",
		CompletedAt:    workflow.Now(ctx),
		ItemsProcessed: int(result.RankingsProcessed),
		UpdatedAt:      workflow.Now(ctx),
	}
	_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

	logger.Info("Rankings workflow completed",
		"totalProcessed", result.RankingsProcessed,
		"specsProcessed", result.SpecsProcessed,
		"dungeonsProcessed", result.DungeonsProcessed,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}
