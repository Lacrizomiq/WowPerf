package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
	"go.temporal.io/sdk/activity"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	buildsStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflowsModels "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// BuildsStatisticsActivity manages all operations related to builds statistics.
type BuildsStatisticsActivity struct {
	playerBuildsRepository     *playerBuildsRepository.PlayerBuildsRepository
	buildsStatisticsRepository *buildsStatisticsRepository.BuildsStatisticsRepository
}

// NewBuildsStatisticsActivity creates a new BuildsStatisticsActivity.
func NewBuildsStatisticsActivity(
	playerBuildsRepository *playerBuildsRepository.PlayerBuildsRepository,
	buildsStatisticsRepository *buildsStatisticsRepository.BuildsStatisticsRepository,
) *BuildsStatisticsActivity {
	return &BuildsStatisticsActivity{
		playerBuildsRepository:     playerBuildsRepository,
		buildsStatisticsRepository: buildsStatisticsRepository,
	}
}

// ProcessItemStatistics processes item statistics for a specific class, spec and encounter_id
func (a *BuildsStatisticsActivity) ProcessItemStatistics(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	batchSize int,
) (*workflowsModels.BuildsAnalysisResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflowsModels.BuildsAnalysisResult{
		ProcessedAt: time.Now(),
	}

	// 1. Delete existing build statistics
	if err := a.buildsStatisticsRepository.DeleteBuildStatistics(ctx, class, spec, encounterID); err != nil {
		return nil, fmt.Errorf("failed to delete existing build statistics: %w", err)
	}

	// 2. Get the total number of builds to process
	count, err := a.countBuilds(ctx, class, spec, encounterID)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		logger.Info("No builds found to analyze",
			"class", class,
			"spec", spec,
			"encounterID", encounterID)
		return result, nil
	}

	// 3. Process the builds by batches
	offset := 0
	totalProcessed := 0
	totalItems := 0

	for offset < int(count) {
		// Record heartbeat
		activity.RecordHeartbeat(ctx, map[string]interface{}{
			"status":     "processing",
			"class":      class,
			"spec":       spec,
			"encounter":  encounterID,
			"progress":   fmt.Sprintf("%d/%d", totalProcessed, count),
			"percentage": float64(totalProcessed) / float64(count) * 100,
		})

		// Get a batch of builds
		builds, err := a.getPlayerBuildsBatch(ctx, class, spec, encounterID, batchSize, offset)
		if err != nil {
			return nil, err
		}

		if len(builds) == 0 {
			break
		}

		// Process the batch
		buildStats, err := a.processBuildsBatch(builds)
		if err != nil {
			return nil, err
		}

		// Calculate the usage percentages
		a.calculateUsagePercentages(buildStats)

		// Persist the statistics
		if len(buildStats) > 0 {
			if err := a.buildsStatisticsRepository.StoreManyBuildStatistics(ctx, buildStats); err != nil {
				return nil, fmt.Errorf("failed to store build statistics: %w", err)
			}
			totalItems += len(buildStats)
		}

		totalProcessed += len(builds)
		offset += batchSize

		logger.Info("Processed builds batch",
			"class", class,
			"spec", spec,
			"encounter", encounterID,
			"batchSize", len(builds),
			"progress", fmt.Sprintf("%d/%d", totalProcessed, count))
	}

	result.BuildsProcessed = int32(totalProcessed)
	result.ItemsProcessed = int32(totalItems)

	logger.Info("Completed build statistics processing",
		"class", class,
		"spec", spec,
		"encounter", encounterID,
		"buildsProcessed", totalProcessed,
		"itemsProcessed", totalItems)

	return result, nil
}

// GearItem represents an item in the gear JSON field
type GearItem struct {
	ID                   int       `json:"id"`
	Slot                 int       `json:"slot"`
	Name                 string    `json:"name"`
	Icon                 string    `json:"icon"`
	Quality              int       `json:"quality"`
	ItemLevel            float64   `json:"itemLevel"`
	SetID                int       `json:"setID"`
	BonusIDs             []int64   `json:"bonusIDs"`
	Gems                 []GemItem `json:"gems"`
	PermanentEnchant     int       `json:"permanentEnchant"`
	PermanentEnchantName string    `json:"permanentEnchantName"`
	TemporaryEnchant     int       `json:"temporaryEnchant"`
	TemporaryEnchantName string    `json:"temporaryEnchantName"`
}

// GemItem represents a gem in an item
type GemItem struct {
	ID        int     `json:"id"`
	Icon      string  `json:"icon"`
	ItemLevel float64 `json:"itemLevel"`
}

// countBuilds sum the number of builds for a specific class, spec and encounter_id
func (a *BuildsStatisticsActivity) countBuilds(
	ctx context.Context,
	class, spec string,
	encounterID uint,
) (int64, error) {
	return a.playerBuildsRepository.CountPlayerBuildsByFilter(ctx, class, spec, encounterID)
}

// getPlayerBuildsBatch get a batch of builds for a specific class, spec and encounter_id with pagination
func (a *BuildsStatisticsActivity) getPlayerBuildsBatch(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	limit, offset int,
) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	return a.playerBuildsRepository.GetPlayerBuildsByFilter(ctx, class, spec, encounterID, limit, offset)
}

// processBuildsBatch process a batch of builds to extract the equipment statistics
// with parallel processing
func (a *BuildsStatisticsActivity) processBuildsBatch(
	builds []*warcraftlogsBuilds.PlayerBuild,
) ([]*warcraftlogsBuilds.BuildStatistic, error) {
	if len(builds) == 0 {
		return nil, nil
	}

	// Map to store the aggregated results
	itemStats := make(map[string]*warcraftlogsBuilds.BuildStatistic)
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
				// Extract the equipment data from the JSON
				var gearItems []GearItem
				if err := json.Unmarshal([]byte(build.Gear), &gearItems); err != nil {
					errCh <- fmt.Errorf("error parsing gear JSON for build %d: %w", build.ID, err)
					continue
				}

				// Process each equipment item
				for _, item := range gearItems {
					// Ignore empty slots
					if item.ID == 0 {
						continue
					}

					// Unique key for each combination of class/spec/encounter/slot/item
					key := fmt.Sprintf("%s_%s_%d_%d_%d",
						build.Class, build.Spec, build.EncounterID, item.Slot, item.ID)

					// Acquire the lock to update the shared map
					mutex.Lock()

					// Create or update the statistic
					stat, exists := itemStats[key]
					if !exists {
						stat = &warcraftlogsBuilds.BuildStatistic{
							Class:       build.Class,
							Spec:        build.Spec,
							EncounterID: build.EncounterID,
							ItemSlot:    item.Slot,
							ItemID:      item.ID,
							ItemName:    item.Name,
							ItemIcon:    item.Icon,
							ItemQuality: item.Quality,
							ItemLevel:   item.ItemLevel,
							// Initialize the fields for min/max item level and keystone level
							MinItemLevel:     build.ItemLevel,
							MaxItemLevel:     build.ItemLevel,
							MinKeystoneLevel: build.KeystoneLevel,
							MaxKeystoneLevel: build.KeystoneLevel,
						}

						// Process the set bonus
						if item.SetID > 0 {
							stat.HasSetBonus = true
							stat.SetID = item.SetID
						}

						// Process the bonus IDs
						if len(item.BonusIDs) > 0 {
							stat.BonusIDs = pq.Int64Array(item.BonusIDs)
						}

						// Process the gems
						if len(item.Gems) > 0 {
							stat.HasGems = true
							stat.GemsCount = len(item.Gems)

							gemIDs := make([]int64, len(item.Gems))
							gemIcons := make([]string, len(item.Gems))
							gemLevels := make([]float64, len(item.Gems))

							for i, gem := range item.Gems {
								gemIDs[i] = int64(gem.ID)
								gemIcons[i] = gem.Icon
								gemLevels[i] = gem.ItemLevel
							}

							stat.GemIDs = pq.Int64Array(gemIDs)
							stat.GemIcons = pq.StringArray(gemIcons)
							stat.GemLevels = pq.Float64Array(gemLevels)
						}

						// Process the enchantments
						if item.PermanentEnchant > 0 {
							stat.HasPermanentEnchant = true
							stat.PermanentEnchantID = item.PermanentEnchant
							stat.PermanentEnchantName = item.PermanentEnchantName
						}

						if item.TemporaryEnchant > 0 {
							stat.HasTemporaryEnchant = true
							stat.TemporaryEnchantID = item.TemporaryEnchant
							stat.TemporaryEnchantName = item.TemporaryEnchantName
						}

						itemStats[key] = stat
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
	result := make([]*warcraftlogsBuilds.BuildStatistic, 0, len(itemStats))
	for _, stat := range itemStats {
		// Calculate the averages
		if stat.UsageCount > 0 {
			stat.AvgItemLevel /= float64(stat.UsageCount)
			stat.AvgKeystoneLevel /= float64(stat.UsageCount)
		}

		result = append(result, stat)
	}

	return result, nil
}

// calculateUsagePercentages calculates the usage percentages
func (a *BuildsStatisticsActivity) calculateUsagePercentages(
	stats []*warcraftlogsBuilds.BuildStatistic,
) {
	// Group by slot
	slotStats := make(map[int]int)
	for _, stat := range stats {
		slotStats[stat.ItemSlot] += stat.UsageCount
	}

	// Calculate the percentages
	for _, stat := range stats {
		totalForSlot := slotStats[stat.ItemSlot]
		if totalForSlot > 0 {
			stat.UsagePercentage = float64(stat.UsageCount) / float64(totalForSlot) * 100
		}
	}
}
