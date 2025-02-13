package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"fmt"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	rankingsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
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

	// Add validation for spec and dungeon
	if spec.ClassName == "" || spec.SpecName == "" {
		return nil, fmt.Errorf("invalid spec configuration: class=%s, spec=%s", spec.ClassName, spec.SpecName)
	}

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

	// Fetch rankings with retry handling
	rankings, err := a.fetchRankingsWithRetry(ctx, spec, dungeon, batchConfig)
	if err != nil {
		return nil, err
	}

	// Store rankings if we got any
	if len(rankings) > 0 {
		logger.Info("Storing rankings",
			"count", len(rankings),
			"class", spec.ClassName,
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

// fetchRankingsWithRetry fetches rankings with retry handling
func (a *RankingsActivity) fetchRankingsWithRetry(ctx context.Context, spec workflows.ClassSpec, dungeon workflows.Dungeon, batchConfig workflows.BatchConfig) ([]*warcraftlogsBuilds.ClassRanking, error) {
	logger := activity.GetLogger(ctx)
	var rankings []*warcraftlogsBuilds.ClassRanking

	variables := map[string]interface{}{
		"encounterId": int(dungeon.EncounterID),
		"className":   spec.ClassName,
		"specName":    spec.SpecName,
		"page":        1,
	}

	logger.Debug("Making GraphQL request",
		"class", spec.ClassName,
		"spec", spec.SpecName,
		"variables", variables)

	response, err := a.client.MakeRequest(ctx, rankingsQueries.ClassRankingsQuery, variables)
	if err != nil {
		if wlErr, ok := err.(*warcraftlogsTypes.WarcraftLogsError); ok {
			switch wlErr.Type {
			case warcraftlogsTypes.ErrorTypeRateLimit, warcraftlogsTypes.ErrorTypeQuotaExceeded:
				logger.Info("Rate limit reached",
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

	// Log progress with class/spec info
	logger.Info("Fetched rankings",
		"class", spec.ClassName,
		"spec", spec.SpecName,
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

	if className == "" || specName == "" {
		return nil, fmt.Errorf("invalid class/spec parameters: class=%s, spec=%s", className, specName)
	}

	rankings, err := a.repository.GetRankingsForSpec(ctx, className, specName, encounterID)
	if err != nil {
		logger.Error("Failed to get stored rankings",
			"class", className,
			"spec", specName,
			"error", err)
		return nil, err
	}

	logger.Info("Retrieved stored rankings",
		"class", className,
		"spec", specName,
		"count", len(rankings))

	return rankings, nil
}
