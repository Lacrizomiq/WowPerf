package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// Processor handles the builds processing logic
type Processor struct {
	totalProcessed int32
}

// NewProcessor creates a new builds processor
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessBuilds processes a batch of builds from reports
func (p *Processor) ProcessBuilds(
	ctx workflow.Context,
	reports []*warcraftlogsBuilds.Report,
	workerConfig models.WorkerConfig,
) (*models.BatchResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Processing builds from reports", "count", len(reports))

	// Configure activity options with parallel processing
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 2,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Execute builds processing activity using definition constant
	var result models.BatchResult
	err := workflow.ExecuteActivity(ctx,
		definitions.ProcessBuildsActivity,
		reports,
	).Get(ctx, &result)

	if err != nil {
		return nil, err
	}

	// Update total processed
	p.totalProcessed += result.ProcessedItems

	return &result, nil
}

// GetStoredReports retrieves stored reports for processing
func (p *Processor) GetStoredReports(ctx workflow.Context) ([]*warcraftlogsBuilds.Report, error) {
	var reports []*warcraftlogsBuilds.Report
	err := workflow.ExecuteActivity(ctx,
		definitions.GetReportsBatchActivity,
		100, // batchSize
		0,   // offset
	).Get(ctx, &reports)

	if err != nil {
		return nil, err
	}

	return reports, nil
}

// CountStoredReports returns the total number of reports available
func (p *Processor) CountStoredReports(ctx workflow.Context) (int64, error) {
	var count int64
	err := workflow.ExecuteActivity(ctx,
		definitions.CountAllReportsActivity,
	).Get(ctx, &count)

	return count, err
}

// GetTotalProcessed returns the total number of builds processed
func (p *Processor) GetTotalProcessed() int32 {
	return p.totalProcessed
}
