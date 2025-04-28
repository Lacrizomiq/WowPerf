package warcraftlogsBuildsTemporalWorkflowsBuildsStatisticsEquipmentStatistics

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

// EquipmentAnalysisWorkflow implements the equipment analysis workflow
type EquipmentAnalysisWorkflow struct{}

// NewEquipmentAnalysisWorkflow creates a new instance of the equipment analysis workflow
func NewEquipmentAnalysisWorkflow() definitions.EquipmentAnalysisWorkflow {
	return &EquipmentAnalysisWorkflow{}
}

// Execute runs the equipment analysis workflow
func (w *EquipmentAnalysisWorkflow) Execute(ctx workflow.Context, params models.EquipmentAnalysisWorkflowParams) (*models.EquipmentAnalysisWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting equipment analysis workflow",
		"specCount", len(params.Spec),
		"dungeonCount", len(params.Dungeon),
		"batchSize", params.BatchSize)

	// Initialize the result
	result := &models.EquipmentAnalysisWorkflowResult{
		StartedAt: workflow.Now(ctx),
		BatchID:   params.BatchID,
	}

	// Validate the parameters
	if len(params.Spec) == 0 {
		return nil, fmt.Errorf("no specs found in parameters")
	}

	if len(params.Dungeon) == 0 {
		return nil, fmt.Errorf("no dungeons found in parameters")
	}

	// Generate a unique ID for the workflow
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	workflowStateID := fmt.Sprintf("equipment-analysis-%s", workflowID)

	// Options for the state management activities
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

	// Create the initial workflow state
	err := workflow.ExecuteActivity(stateCtx, definitions.CreateWorkflowStateActivity, &warcraftlogsBuilds.WorkflowState{
		ID:              workflowStateID,
		WorkflowType:    "equipment-analysis",
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

	// Options for the analysis activities
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 24,
		HeartbeatTimeout:    time.Minute * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * time.Duration(params.RetryDelay.Seconds()),
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    int32(params.RetryAttempts),
		},
	}
	activityCtx := workflow.WithActivityOptions(ctx, activityOpts)

	// Process each spec/dungeon combination
	processedCombinations := make(map[string]bool)
	totalBuilds := int32(0)
	totalItems := int32(0)
	specsProcessed := int32(0)
	dungeonsProcessed := make(map[string]bool)

	for _, spec := range params.Spec {
		specProcessed := false

		for _, dungeon := range params.Dungeon {
			// Identify the combination
			combinationKey := fmt.Sprintf("%s_%s_%d", spec.ClassName, spec.SpecName, dungeon.EncounterID)

			// Check if already processed
			if processedCombinations[combinationKey] {
				logger.Info("Skipping already processed combination",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name)
				continue
			}

			// Update the workflow state
			err = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, &warcraftlogsBuilds.WorkflowState{
				ID:              workflowStateID,
				Status:          "running",
				LastProcessedID: combinationKey,
				UpdatedAt:       workflow.Now(ctx),
			}).Get(ctx, nil)

			if err != nil {
				logger.Error("Failed to update workflow state", "error", err)
			}

			logger.Info("Processing equipment analysis",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name)

			// Execute the equipment analysis activity
			var activityResult models.EquipmentAnalysisWorkflowResult
			err := workflow.ExecuteActivity(activityCtx,
				definitions.ProcessBuildStatisticsActivity,
				spec.ClassName,
				spec.SpecName,
				uint(dungeon.EncounterID),
				int(params.BatchSize),
			).Get(ctx, &activityResult)

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

					logger.Info("Rate limit reached during equipment analysis",
						"class", spec.ClassName,
						"spec", spec.SpecName,
						"dungeon", dungeon.Name)

					result.CompletedAt = workflow.Now(ctx)
					return result, err
				}

				logger.Error("Failed to process equipment analysis",
					"class", spec.ClassName,
					"spec", spec.SpecName,
					"dungeon", dungeon.Name,
					"error", err)

				// Update workflow state with error
				workflowState := &warcraftlogsBuilds.WorkflowState{
					ID:           workflowStateID,
					ErrorMessage: fmt.Sprintf("Error processing %s: %v", combinationKey, err),
					UpdatedAt:    workflow.Now(ctx),
				}
				_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

				// Continue with the next combination on error
				continue
			}

			// Mark the combination as processed
			processedCombinations[combinationKey] = true
			dungeonsProcessed[dungeon.Name] = true
			specProcessed = true

			// Update the counters
			totalBuilds += activityResult.TotalBuilds
			totalItems += activityResult.ItemsAnalyzed

			// Update the workflow state with progress
			workflowState := &warcraftlogsBuilds.WorkflowState{
				ID:             workflowStateID,
				ItemsProcessed: int(totalItems),
				UpdatedAt:      workflow.Now(ctx),
			}
			_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

			logger.Info("Successfully processed equipment analysis",
				"class", spec.ClassName,
				"spec", spec.SpecName,
				"dungeon", dungeon.Name,
				"buildsProcessed", activityResult.TotalBuilds,
				"itemsProcessed", activityResult.ItemsAnalyzed)

			// Small delay between combinations to avoid overloading the system
			workflow.Sleep(ctx, time.Second*2)
		}

		// Increment the specs counter
		if specProcessed {
			specsProcessed++
		}
	}

	// Finalize the result
	result.TotalBuilds = totalBuilds
	result.ItemsAnalyzed = totalItems
	result.SpecsProcessed = specsProcessed
	result.DungeonsProcessed = int32(len(dungeonsProcessed))
	result.CompletedAt = workflow.Now(ctx)

	// Complete the workflow state
	workflowState := &warcraftlogsBuilds.WorkflowState{
		ID:             workflowStateID,
		Status:         "completed",
		CompletedAt:    workflow.Now(ctx),
		ItemsProcessed: int(totalItems),
		UpdatedAt:      workflow.Now(ctx),
	}
	_ = workflow.ExecuteActivity(stateCtx, definitions.UpdateWorkflowStateActivity, workflowState).Get(ctx, nil)

	logger.Info("Equipment analysis workflow completed",
		"totalBuilds", totalBuilds,
		"itemsAnalyzed", totalItems,
		"specsProcessed", specsProcessed,
		"dungeonsProcessed", len(dungeonsProcessed),
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}
