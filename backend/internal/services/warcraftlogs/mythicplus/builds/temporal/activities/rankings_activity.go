package warcraftlogsBuildsTemporal

import (
	"context"
	"fmt"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	rankingsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

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

// FetchAndStore fetches and stores rankings for a given class spec and dungeon
func (a *RankingsActivity) FetchAndStore(ctx context.Context, spec workflows.ClassSpec, dungeon workflows.Dungeon, batchConfig workflows.BatchConfig) (*workflows.BatchResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Starting rankings fetch activity",
		"class", spec.ClassName,
		"spec", spec.SpecName,
		"dungeon", dungeon.Name)

	result := &workflows.BatchResult{
		ClassName:   spec.ClassName,
		SpecName:    spec.SpecName,
		EncounterID: dungeon.EncounterID,
		ProcessedAt: time.Now(),
	}

	// Check if an update is necessary
	lastRanking, err := a.repository.GetLastRankingForEncounter(ctx, dungeon.EncounterID)
	if err != nil {
		return nil, fmt.Errorf("failed to check last ranking: %w", err)
	}

	if lastRanking != nil {
		// if last ranking is not too old, we can skip the fetch
		if time.Since(lastRanking.UpdatedAt) < 7*24*time.Hour {
			logger.Info("Rankings recently updated, skipping",
				"lastUpdate", lastRanking.UpdatedAt)
			return result, nil
		}
	}

	// Configuration of the variables for the request
	variables := map[string]interface{}{
		"encounterId": dungeon.EncounterID,
		"className":   spec.ClassName,
		"specName":    spec.SpecName,
		"page":        1,
	}

	// Retrieve the rankings from WarcraftLogs with Temporal retry policy
	var rankings []*warcraftlogsBuilds.ClassRanking
	err = a.fetchRankings(ctx, variables, &rankings)
	if err != nil {
		return nil, err
	}

	// Store the rankings in the database
	if len(rankings) > 0 {
		logger.Info("Storing rankings", "count", len(rankings))
		if err := a.repository.StoreRankings(ctx, dungeon.EncounterID, rankings); err != nil {
			return nil, fmt.Errorf("failed to store rankings: %w", err)
		}
		result.Rankings = rankings
	}

	activity.RecordHeartbeat(ctx, fmt.Sprintf("Processed %d rankings", len(rankings)))

	return result, nil
}

func (a *RankingsActivity) fetchRankings(ctx context.Context, variables map[string]interface{}, rankings *[]*warcraftlogsBuilds.ClassRanking) error {
	response, err := a.client.MakeRequest(ctx, rankingsQueries.ClassRankingsQuery, variables)
	if err != nil {
		return temporal.NewApplicationError(
			fmt.Sprintf("API request failed: %v", err),
			"API_ERROR",
		)
	}

	fetchedRankings, hasMore, err := rankingsQueries.ParseRankingsResponse(
		response,
		uint(variables["encounterId"].(int)),
	)

	if err != nil {
		return temporal.NewNonRetryableApplicationError(
			"Failed to parse rankings",
			"PARSE_ERROR",
			err,
		)
	}

	*rankings = append(*rankings, fetchedRankings...)

	// Log progress
	activity.GetLogger(ctx).Info("Fetched rankings batch",
		"count", len(fetchedRankings),
		"hasMore", hasMore,
		"total", len(*rankings))

	return nil
}
