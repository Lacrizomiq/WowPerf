package warcraftlogsBuildsTemporalActivities

import (
	"context"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/activity"
)

// RateLimitActivity is an activity that handles rate limiting for WarcraftLogs API requests
type RateLimitActivity struct {
	client *warcraftlogs.WarcraftLogsClientService
}

// NewRateLimitActivity creates a new instance of RateLimitActivity
func NewRateLimitActivity(client *warcraftlogs.WarcraftLogsClientService) *RateLimitActivity {
	return &RateLimitActivity{
		client: client,
	}
}

// ReservePoints attempt to reserve points for a workflow execution
func (a *RateLimitActivity) ReservePoints(ctx context.Context, params workflows.WorkflowParams) error {
	logger := activity.GetLogger(ctx)

	// Calculate required points based on specs and dungeons
	requiredPoints := calculateRequiredPoints(params)

	// Try to reserve points
	if err := a.client.GetRateLimiter().ReserveWorkflowPoints(params.WorkflowID, requiredPoints); err != nil {
		logger.Error("Failed to reserve points",
			"workflowID", params.WorkflowID,
			"requiredPoints", requiredPoints,
			"error", err,
		)
		return err
	}
	logger.Info("Reserved points for workflow",
		"workflowID", params.WorkflowID,
		"points", requiredPoints)

	return nil
}

// ReleasePoints releases points reserved for a workflow
func (a *RateLimitActivity) ReleasePoints(ctx context.Context, params workflows.WorkflowParams) error {
	logger := activity.GetLogger(ctx)

	a.client.GetRateLimiter().ReleaseWorkflowPoints(params.WorkflowID)

	logger.Info("Releasing points for workflow",
		"workflowID", params.WorkflowID,
	)

	return nil
}

// calculateRequiredPoints estimates points needed for the workflow
func calculateRequiredPoints(params workflows.WorkflowParams) float64 {
	specsCount := len(params.Config.Specs)
	dungeonsCount := len(params.Config.Dungeons)

	// Base calculation:
	// 1 point per ranking request
	// 2 points per report (details + talents)
	pointsPerCombo := 3.0 // 1 + 2

	// Total combinations
	totalCombos := float64(specsCount * dungeonsCount)

	// Add 10% buffer for retries and overhead
	return pointsPerCombo * totalCombos * 1.1
}
