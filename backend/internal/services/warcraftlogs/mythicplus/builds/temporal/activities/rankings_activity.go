package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"fmt"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	rankingsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

type RankingsActivity struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *rankingsRepository.RankingsRepository
}

func NewRankingsActivity(client *warcraftlogs.WarcraftLogsClientService, repository *rankingsRepository.RankingsRepository) *RankingsActivity {
	return &RankingsActivity{
		client:     client,
		repository: repository,
	}
}

func (a *RankingsActivity) FetchAndStore(ctx context.Context, spec workflows.ClassSpec, dungeon workflows.Dungeon, batchConfig workflows.BatchConfig) (*workflows.BatchResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting rankings fetch activity",
		"class", spec.ClassName,
		"spec", spec.SpecName,
		"dungeon", dungeon.Name,
		"encounterId", dungeon.EncounterID,
	)

	result := &workflows.BatchResult{
		ClassName:   spec.ClassName,
		SpecName:    spec.SpecName,
		EncounterID: dungeon.EncounterID,
		ProcessedAt: time.Now(),
	}

	// Check if an update is necessary for this specific spec
	lastRanking, err := a.repository.GetLastRankingForEncounter(ctx, dungeon.EncounterID, spec.ClassName, spec.SpecName)
	if err != nil {
		return nil, fmt.Errorf("failed to check last ranking: %w", err)
	}

	if lastRanking != nil {
		if time.Since(lastRanking.UpdatedAt) < 7*24*time.Hour {
			logger.Info("Rankings recently updated, skipping",
				"lastUpdate", lastRanking.UpdatedAt,
				"spec", spec.SpecName)
			return result, nil
		}
	}

	// Fetch rankings with retry handling
	rankings, err := a.fetchRankingsWithRetry(ctx, spec, dungeon, batchConfig)
	if err != nil {
		return nil, err
	}

	// Store rankings if we got any
	if len(rankings) > 0 {
		logger.Info("Storing rankings",
			"count", len(rankings),
			"spec", spec.SpecName)
		if err := a.repository.StoreRankings(ctx, dungeon.EncounterID, rankings); err != nil {
			return nil, fmt.Errorf("failed to store rankings: %w", err)
		}
		result.Rankings = rankings
	}

	activity.RecordHeartbeat(ctx, fmt.Sprintf("Processed %d rankings for %s %s",
		len(rankings), spec.ClassName, spec.SpecName))
	return result, nil
}

func (a *RankingsActivity) fetchRankingsWithRetry(ctx context.Context, spec workflows.ClassSpec, dungeon workflows.Dungeon, batchConfig workflows.BatchConfig) ([]*warcraftlogsBuilds.ClassRanking, error) {
	logger := activity.GetLogger(ctx)
	var rankings []*warcraftlogsBuilds.ClassRanking

	variables := map[string]interface{}{
		"encounterId": int(dungeon.EncounterID),
		"className":   spec.ClassName,
		"specName":    spec.SpecName,
		"page":        1,
	}

	logger.Debug("Making GraphQL request", "variables", variables)

	response, err := a.client.MakeRequest(ctx, rankingsQueries.ClassRankingsQuery, variables)
	if err != nil {
		if wlErr, ok := err.(*warcraftlogsTypes.WarcraftLogsError); ok {
			switch wlErr.Type {
			case warcraftlogsTypes.ErrorTypeRateLimit, warcraftlogsTypes.ErrorTypeQuotaExceeded:
				logger.Info("Rate limit reached, will retry",
					"resetIn", wlErr.RetryIn)
				return nil, temporal.NewApplicationError(
					fmt.Sprintf("Rate limit reached: %v", wlErr),
					"RATE_LIMIT_ERROR",
				)
			case warcraftlogsTypes.ErrorTypeAPI:
				if !wlErr.Retryable {
					return nil, temporal.NewNonRetryableApplicationError(
						fmt.Sprintf("Non-retryable API error: %v", wlErr),
						"API_ERROR",
						wlErr,
					)
				}
			}
		}
		return nil, temporal.NewApplicationError(
			fmt.Sprintf("Failed to fetch rankings: %v", err),
			"FETCH_ERROR",
		)
	}

	fetchedRankings, err := rankingsQueries.ParseRankingsResponse(
		response,
		dungeon.EncounterID,
	)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError(
			"Failed to parse rankings",
			"PARSE_ERROR",
			err,
		)
	}

	rankings = append(rankings, fetchedRankings...)

	// Log progress
	logger.Info("Fetched rankings",
		"count", len(fetchedRankings),
		"total", len(rankings))

	return rankings, nil
}

// GetStoredRankings retrieves rankings from the database
func (a *RankingsActivity) GetStoredRankings(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting stored rankings",
		"class", className,
		"spec", specName,
		"encounterID", encounterID)

	rankings, err := a.repository.GetRankingsForSpec(ctx, className, specName, encounterID)
	if err != nil {
		logger.Error("Failed to get stored rankings", "error", err)
		return nil, err
	}

	logger.Info("Retrieved stored rankings", "count", len(rankings))
	return rankings, nil
}
