package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"time"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"

	"go.temporal.io/sdk/activity"
)

// RateLimitActivity handles rate limiting for WarcraftLogs API requests
type RateLimitActivity struct {
	client *warcraftlogs.WarcraftLogsClientService
}

// NewRateLimitActivity creates a new instance of RateLimitActivity
func NewRateLimitActivity(client *warcraftlogs.WarcraftLogsClientService) *RateLimitActivity {
	return &RateLimitActivity{
		client: client,
	}
}

// CheckRemainingPoints returns the remaining points based on local tracking
func (a *RateLimitActivity) CheckRemainingPoints(ctx context.Context, _ workflows.WorkflowParams) (float64, error) {
	logger := activity.GetLogger(ctx)
	info := a.client.GetRateLimiter().GetRateLimitInfo()

	logger.Info("Checking remaining points",
		"remainingPoints", info.RemainingPoints,
		"resetIn", info.ResetIn)

	return info.RemainingPoints, nil
}

// ReservePoints checks if we have enough points
func (a *RateLimitActivity) ReservePoints(ctx context.Context, params workflows.WorkflowParams) error {
	logger := activity.GetLogger(ctx)
	info := a.client.GetRateLimiter().GetRateLimitInfo()

	required := estimateRequiredPoints(&params)

	// Simple check
	if info.RemainingPoints < required {
		logger.Warn("Insufficient points for workflow",
			"available", info.RemainingPoints,
			"required", required,
			"resetIn", info.ResetIn)

		// If the reset is close, wait
		if info.ResetIn < time.Minute*15 {
			return nil // Let the workflow continue
		}

		return warcraftlogsTypes.NewQuotaExceededError(info)
	}

	return nil
}

// ReleasePoints is now a monitoring operation
func (a *RateLimitActivity) ReleasePoints(ctx context.Context, params workflows.WorkflowParams) error {
	logger := activity.GetLogger(ctx)
	info := a.client.GetRateLimiter().GetRateLimitInfo()

	logger.Info("Workflow completion status",
		"remainingPoints", info.RemainingPoints,
		"resetIn", info.ResetIn,
		"workflowID", params.WorkflowID)

	return nil
}

// estimateRequiredPoints provides a rough estimate of points needed for a workflow
func estimateRequiredPoints(params *workflows.WorkflowParams) float64 {
	if params == nil || params.Config == nil {
		return 1.0
	}

	// Base cost for workflow
	totalPoints := 1.0

	// Calculate points needed for each spec/dungeon combination
	numSpecs := len(params.Config.Specs)
	numDungeons := len(params.Config.Dungeons)

	// Points per combination:
	// - 1 point for rankings query
	// - ~2 points for processing reports (average)
	pointsPerCombo := 3.0

	totalPoints += float64(numSpecs*numDungeons) * pointsPerCombo

	// Add 20% buffer for unexpected operations
	return totalPoints * 1.2
}
