package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

// BuildsService encapsulates the builds processing logic
type BuildsService struct {
	processor *Processor
}

// NewBuildsService creates a new builds service
func NewBuildsService() *BuildsService {
	return &BuildsService{
		processor: NewProcessor(),
	}
}

// ProcessAllBuilds processes all builds from reports
// It updates the state with the processing results
func (s *BuildsService) ProcessAllBuilds(
	ctx workflow.Context,
	params models.WorkflowConfig,
	state *state.WorkflowState,
) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting builds processing")

	// Get reports to process
	reports, err := s.processor.GetStoredReports(ctx)
	if err != nil {
		logger.Error("Failed to get stored reports", "error", err)
		return err
	}

	if len(reports) > 0 {
		logger.Info("Processing builds from reports", "count", len(reports))

		// Process builds in batches
		batchResult, err := s.processor.ProcessBuilds(ctx, reports, params.Worker)
		if err != nil {
			if common.IsRateLimitError(err) {
				// Rate limit reached, bail out and let main workflow handle continuation
				logger.Info("Rate limit reached during builds processing")
				return err
			}
			logger.Error("Failed to process builds", "error", err)
			return err
		}

		// Update state with results
		state.PartialResults.BuildsProcessed += batchResult.ProcessedItems

		logger.Info("Completed processing builds",
			"processedCount", batchResult.ProcessedItems,
			"totalProcessed", state.PartialResults.BuildsProcessed)
	} else {
		logger.Info("No reports found to process builds")
	}

	return nil
}
