package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	talentStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflowsModels "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// TalentStatisticActivity manages all operations related to talent statistics.
type TalentStatisticActivity struct {
	playerBuildsRepository     *playerBuildsRepository.PlayerBuildsRepository
	talentStatisticsRepository *talentStatisticsRepository.TalentStatisticsRepository
}

// NewTalentStatisticActivity creates a new TalentStatisticActivity.
func NewTalentStatisticActivity(
	playerBuildsRepository *playerBuildsRepository.PlayerBuildsRepository,
	talentStatisticsRepository *talentStatisticsRepository.TalentStatisticsRepository,
) *TalentStatisticActivity {
	return &TalentStatisticActivity{
		playerBuildsRepository:     playerBuildsRepository,
		talentStatisticsRepository: talentStatisticsRepository,
	}
}

// ProcessTalentStatistics analyzes talent configurations for a specific class, spec and dungeon
func (a *TalentStatisticActivity) ProcessTalentStatistics(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	batchSize int,
) (*workflowsModels.TalentAnalysisWorkflowResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflowsModels.TalentAnalysisWorkflowResult{
		StartedAt: time.Now(),
	}

	// 1. Delete existing statistics
	if err := a.talentStatisticsRepository.DeleteTalentStatistics(ctx, class, spec, encounterID); err != nil {
		return nil, fmt.Errorf("failed to delete existing talent statistics: %w", err)
	}

	// 2. Get the total number of builds to process
	count, err := a.playerBuildsRepository.CountPlayerBuildsNeedingTalentAnalysis(ctx, class, spec, encounterID)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		logger.Info("No builds found to analyze for talents",
			"class", class,
			"spec", spec,
			"encounterID", encounterID)
		return result, nil
	}

	// 3. Process the builds by batches
	offset := 0
	totalProcessed := 0
	totalTalentConfigs := 0

	// For storing the IDs of successfully processed builds
	processedBuildIDs := make([]uint, 0)

	for offset < int(count) {
		// Record heartbeat
		activity.RecordHeartbeat(ctx, map[string]interface{}{
			"status":     "processing_talents",
			"class":      class,
			"spec":       spec,
			"encounter":  encounterID,
			"progress":   fmt.Sprintf("%d/%d", totalProcessed, count),
			"percentage": float64(totalProcessed) / float64(count) * 100,
		})

		// Get a batch of builds
		builds, err := a.playerBuildsRepository.GetPlayerBuildsNeedingTalentAnalysis(
			ctx, class, spec, encounterID, batchSize, offset)
		if err != nil {
			return nil, err
		}

		if len(builds) == 0 {
			break
		}

		// Collect the IDs of the builds to mark them later
		batchBuildIDs := make([]uint, 0, len(builds))
		for _, build := range builds {
			batchBuildIDs = append(batchBuildIDs, build.ID)
		}

		// Process the batch
		talentStats, err := a.ProcessTalentsBatch(builds)
		if err != nil {
			// In case of error, mark these builds as failed
			if len(batchBuildIDs) > 0 {
				_ = a.playerBuildsRepository.MarkPlayerBuildsAsProcessedForTalent(
					ctx, batchBuildIDs, "failed")
			}
			return nil, err
		}

		// Calculate the usage percentages
		a.CalculateUsagePercentages(talentStats, len(builds))

		// Persist the statistics
		if len(talentStats) > 0 {
			if err := a.talentStatisticsRepository.StoreManyTalentStatistics(ctx, talentStats); err != nil {
				// In case of error, mark these builds as failed
				if len(batchBuildIDs) > 0 {
					_ = a.playerBuildsRepository.MarkPlayerBuildsAsProcessedForTalent(
						ctx, batchBuildIDs, "failed")
				}
				return nil, fmt.Errorf("failed to store talent statistics: %w", err)
			}
			totalTalentConfigs += len(talentStats)

			// Add the IDs of successfully processed builds
			processedBuildIDs = append(processedBuildIDs, batchBuildIDs...)
		}

		totalProcessed += len(builds)
		offset += batchSize

		logger.Info("Processed talents batch",
			"class", class,
			"spec", spec,
			"encounter", encounterID,
			"batchSize", len(builds),
			"uniqueTalentConfigs", len(talentStats),
			"progress", fmt.Sprintf("%d/%d", totalProcessed, count))
	}

	// Mark all successfully processed builds
	if len(processedBuildIDs) > 0 {
		if err := a.playerBuildsRepository.MarkPlayerBuildsAsProcessedForTalent(
			ctx, processedBuildIDs, "processed"); err != nil {
			logger.Error("Failed to mark builds as processed for talents",
				"error", err,
				"buildsCount", len(processedBuildIDs))
			// Continue despite the error
		}
	}

	result.TotalBuilds = int32(totalProcessed)
	result.TalentsAnalyzed = int32(totalTalentConfigs)
	result.CompletedAt = time.Now()

	logger.Info("Completed talent analysis",
		"class", class,
		"spec", spec,
		"encounter", encounterID,
		"buildsProcessed", totalProcessed,
		"talentConfigsProcessed", totalTalentConfigs,
		"duration", result.CompletedAt.Sub(result.StartedAt))

	return result, nil
}

// countBuilds sum the number of builds for a specific class, spec and encounter_id
func (a *TalentStatisticActivity) countBuilds(
	ctx context.Context,
	class, spec string,
	encounterID uint,
) (int64, error) {
	return a.playerBuildsRepository.CountPlayerBuildsByFilter(ctx, class, spec, encounterID)
}

// getPlayerBuildsBatch get a batch of builds for a specific class, spec and encounter_id with pagination
func (a *TalentStatisticActivity) getPlayerBuildsBatch(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	limit, offset int,
) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	return a.playerBuildsRepository.GetPlayerBuildsByFilter(ctx, class, spec, encounterID, limit, offset)
}

// ProcessTalentsBatch process a batch of builds to extract the talent statistics
// with parallel processing
func (a *TalentStatisticActivity) ProcessTalentsBatch(
	builds []*warcraftlogsBuilds.PlayerBuild,
) ([]*warcraftlogsBuilds.TalentStatistic, error) {
	if len(builds) == 0 {
		return nil, nil
	}

	// Map to store the aggregated results
	talentStats := make(map[string]*warcraftlogsBuilds.TalentStatistic)
	var mutex sync.Mutex

	// Number of workers to use for parallel processing
	numWorkers := 4
	if len(builds) < numWorkers {
		numWorkers = len(builds)
	}

	// Divide the work between the workers
	buildCh := make(chan *warcraftlogsBuilds.PlayerBuild)
	errCh := make(chan error, numWorkers)
	var wg sync.WaitGroup

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for build := range buildCh {
				// Check that talent_import is not empty
				if build.TalentImport == "" {
					continue
				}

				// Acquire the lock to update the shared map
				mutex.Lock()

				// Create or update the statistic
				stat, exists := talentStats[build.TalentImport]
				if !exists {
					stat = &warcraftlogsBuilds.TalentStatistic{
						Class:        build.Class,
						Spec:         build.Spec,
						EncounterID:  build.EncounterID,
						TalentImport: build.TalentImport,
						// Initialize the fields for min/max item level and keystone level
						MinItemLevel:     build.ItemLevel,
						MaxItemLevel:     build.ItemLevel,
						MinKeystoneLevel: build.KeystoneLevel,
						MaxKeystoneLevel: build.KeystoneLevel,
					}
					talentStats[build.TalentImport] = stat
				}

				// Update the usage statistics
				stat.UsageCount++

				// Update the min/max item level
				if build.ItemLevel < stat.MinItemLevel {
					stat.MinItemLevel = build.ItemLevel
				}
				if build.ItemLevel > stat.MaxItemLevel {
					stat.MaxItemLevel = build.ItemLevel
				}

				if build.KeystoneLevel < stat.MinKeystoneLevel {
					stat.MinKeystoneLevel = build.KeystoneLevel
				}
				if build.KeystoneLevel > stat.MaxKeystoneLevel {
					stat.MaxKeystoneLevel = build.KeystoneLevel
				}

				// Add to the total to calculate the average later
				stat.AvgItemLevel += build.ItemLevel
				stat.AvgKeystoneLevel += float64(build.KeystoneLevel)

				// Release the lock
				mutex.Unlock()
			}
		}()
	}

	// Send the builds to the workers
	for _, build := range builds {
		buildCh <- build
	}
	close(buildCh)

	// Wait for all the workers to finish
	wg.Wait()
	close(errCh)

	// Check if there were any errors
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	// Convert the map to a slice and calculate the final statistics
	result := make([]*warcraftlogsBuilds.TalentStatistic, 0, len(talentStats))
	for _, stat := range talentStats {
		// Calculate the averages
		if stat.UsageCount > 0 {
			stat.AvgItemLevel /= float64(stat.UsageCount)
			stat.AvgKeystoneLevel /= float64(stat.UsageCount)
		}

		result = append(result, stat)
	}

	return result, nil

}

// calculateUsagePercentages calculate the usage percentages
func (a *TalentStatisticActivity) CalculateUsagePercentages(
	stats []*warcraftlogsBuilds.TalentStatistic,
	totalBuilds int,
) {
	// Check that we have builds
	if totalBuilds <= 0 {
		return
	}

	// Calculate the percentage for each configuration
	for _, stat := range stats {
		stat.UsagePercentage = float64(stat.UsageCount) / float64(totalBuilds) * 100
	}
}
