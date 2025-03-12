package warcraftlogsBuildsTemporalWorkflowsReports

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// Processor handles the reports processing logic
type Processor struct {
	totalProcessed int32
}

// NewProcessor creates a new reports processor
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessReports processes a batch of reports for given rankings
func (p *Processor) ProcessReports(
	ctx workflow.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
	workerConfig models.WorkerConfig,
) (*models.BatchResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Processing reports for rankings", "count", len(rankings))

	// Configure activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Execute reports processing activity using definition constant
	var result models.BatchResult
	err := workflow.ExecuteActivity(ctx,
		definitions.ProcessReportsActivity,
		rankings,
	).Get(ctx, &result)

	if err != nil {
		return nil, err
	}

	// Update total processed
	p.totalProcessed += result.ProcessedItems

	return &result, nil
}

// GetStoredRankings retrieves stored rankings for a spec
func (p *Processor) GetStoredRankings(ctx workflow.Context, spec *models.ClassSpec) ([]*warcraftlogsBuilds.ClassRanking, error) {
	if spec == nil {
		return nil, &common.WorkflowError{
			Type:      common.ErrorTypeConfiguration,
			Message:   "spec cannot be nil",
			Retryable: false,
		}
	}

	var rankings []*warcraftlogsBuilds.ClassRanking
	err := workflow.ExecuteActivity(ctx,
		definitions.GetStoredRankingsActivity,
		spec.ClassName,
		spec.SpecName,
	).Get(ctx, &rankings)

	return rankings, err
}

// GetTotalProcessed returns the total number of reports processed
func (p *Processor) GetTotalProcessed() int32 {
	return p.totalProcessed
}
