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

// CheckRemainingPoints gets real-time points status from API
func (a *RateLimitActivity) CheckRemainingPoints(ctx context.Context, _ workflows.WorkflowParams) (float64, error) {
	logger := activity.GetLogger(ctx)

	// Trigger an update if necessary
	limiter := a.client.GetRateLimiter()
	info := limiter.GetRateLimitInfo()

	logger.Info("Rate limit check",
		"remainingPoints", info.RemainingPoints,
		"resetIn", info.ResetIn)

	return info.RemainingPoints, nil
}

// ReservePoints checks if we have enough points for the workflow
func (a *RateLimitActivity) ReservePoints(ctx context.Context, params workflows.WorkflowParams) error {
	logger := activity.GetLogger(ctx)

	// Calculer les points requis
	required := estimateRequiredPoints(&params)

	// Obtenir l'état actuel
	info := a.client.GetRateLimiter().GetRateLimitInfo()

	// Vérifier si nous avons assez de points
	if info.RemainingPoints < required {
		logger.Warn("Insufficient points for workflow",
			"available", info.RemainingPoints,
			"required", required,
			"resetIn", info.ResetIn)

		// Si le reset est proche, attendre
		if info.ResetIn < time.Minute*15 {
			return nil // Laisser le workflow continuer
		}

		return warcraftlogsTypes.NewQuotaExceededError(info)
	}

	logger.Info("Points reserved for workflow",
		"required", required,
		"available", info.RemainingPoints,
		"workflowID", params.WorkflowID)

	return nil
}

// ReleasePoints monitors workflow completion status
func (a *RateLimitActivity) ReleasePoints(ctx context.Context, params workflows.WorkflowParams) error {
	logger := activity.GetLogger(ctx)
	info := a.client.GetRateLimiter().GetRateLimitInfo()

	logger.Info("Workflow completion monitoring",
		"remainingPoints", info.RemainingPoints,
		"resetIn", info.ResetIn,
		"workflowID", params.WorkflowID)

	// Add monitoring metrics if needed
	activity.RecordHeartbeat(ctx, map[string]interface{}{
		"remainingPoints": info.RemainingPoints,
		"resetIn":         info.ResetIn,
		"workflowID":      params.WorkflowID,
	})

	return nil
}

// estimateRequiredPoints calculates points needed including rate limit checks
func estimateRequiredPoints(params *workflows.WorkflowParams) float64 {
	if params == nil || params.Config == nil {
		return 1.0
	}

	// Calculate operations per spec/dungeon combo
	numSpecs := len(params.Config.Specs)
	numDungeons := len(params.Config.Dungeons)

	// Points breakdown per spec/dungeon combination:
	// - Rankings query: ~13 points + 2 points (rate limit check) = 15 points
	// - Reports queries (x2): (~13-16 points + 2 points check) × 2 = ~36 points
	pointsPerCombo := 51.0 // Total: 15 + 36 = 51 points per combination

	totalPoints := float64(numSpecs*numDungeons) * pointsPerCombo

	// Add 20% buffer for unexpected variations
	return totalPoints * 1.2
}
