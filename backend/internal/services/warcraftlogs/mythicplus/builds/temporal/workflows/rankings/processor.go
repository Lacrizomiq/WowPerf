package warcraftlogsBuildsTemporalWorkflowsRankings

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// Processor handles the rankings processing logic
type Processor struct {
	totalProcessed int32
}

// NewProcessor creates a new rankings processor
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessRankings processes rankings for a specific spec and dungeon
func (p *Processor) ProcessRankings(
	ctx workflow.Context,
	spec models.ClassSpec,
	dungeon models.Dungeon,
	batchConfig models.BatchConfig,
) (*models.BatchResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Processing rankings",
		"class", spec.ClassName,
		"spec", spec.SpecName,
		"dungeon", dungeon.Name)

	// Validate input
	if err := common.ValidateSpec(spec); err != nil {
		return nil, err
	}

	// Configure activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    batchConfig.RetryDelay,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    int32(batchConfig.MaxAttempts),
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Execute rankings activity using definition constant
	var result models.BatchResult
	err := workflow.ExecuteActivity(ctx,
		definitions.FetchRankingsActivity,
		spec,
		dungeon,
	).Get(ctx, &result)

	if err != nil {
		return nil, err
	}

	// Update total processed
	p.totalProcessed += result.ProcessedItems

	return &result, nil
}
