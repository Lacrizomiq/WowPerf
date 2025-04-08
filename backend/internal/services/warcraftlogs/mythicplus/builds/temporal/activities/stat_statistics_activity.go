package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	statStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflowsModels "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// StatStatisticsActivity manages all operations related to stat statistics.
type StatStatisticsActivity struct {
	playerBuildsRepository   *playerBuildsRepository.PlayerBuildsRepository
	statStatisticsRepository *statStatisticsRepository.StatStatisticsRepository
}

// NewStatStatisticsActivity creates a new StatStatisticsActivity.
func NewStatStatisticsActivity(
	playerBuildsRepository *playerBuildsRepository.PlayerBuildsRepository,
	statStatisticsRepository *statStatisticsRepository.StatStatisticsRepository,
) *StatStatisticsActivity {
	return &StatStatisticsActivity{
		playerBuildsRepository:   playerBuildsRepository,
		statStatisticsRepository: statStatisticsRepository,
	}
}

// List of stats per category (secondary and minor)
var (
	secondaryStats = map[string]bool{
		"Crit":        true,
		"Haste":       true,
		"Mastery":     true,
		"Versatility": true,
	}

	minorStats = map[string]bool{
		"Leech":     true,
		"Avoidance": true,
		"Speed":     true,
	}
)

// Stat represent a stat with his min/max value in the JSON
type Stat struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// ProcessStatStatistics analyze the stats for a class/spec/dungeon
func (a *StatStatisticsActivity) ProcessStatStatistics(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	batchSize int,
) (*workflowsModels.BuildsAnalysisResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflowsModels.BuildsAnalysisResult{
		ProcessedAt: time.Now(),
	}

	// 1. Delete existing statistics
	if err := a.statStatisticsRepository.DeleteStatStatistics(ctx, class, spec, encounterID); err != nil {
		return nil, fmt.Errorf("failed to delete existing stat statistics: %w", err)
	}

	// 2. Get the total number of builds to process
	count, err := a.countBuilds(ctx, class, spec, encounterID)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		logger.Info("No builds found to analyze stats",
			"class", class,
			"spec", spec,
			"encounterID", encounterID)
		return result, nil
	}

	// 3. Process the builds by batches
	offset := 0
	totalProcessed := 0

	// Structures to store the aggregated statistics
	statData := make(map[string]*StatAggregation)

	for offset < int(count) {
		// Record heartbeat
		activity.RecordHeartbeat(ctx, map[string]interface{}{
			"status":     "processing_stats",
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
		err = a.ProcessStatsBatch(builds, statData)
		if err != nil {
			return nil, err
		}

		totalProcessed += len(builds)
		offset += batchSize

		logger.Info("Processed stats batch",
			"class", class,
			"spec", spec,
			"encounter", encounterID,
			"batchSize", len(builds),
			"progress", fmt.Sprintf("%d/%d", totalProcessed, count))
	}

	// Convert the aggregated data to final statistics
	statStats := a.ConvertToStatStatistics(statData, class, spec, encounterID)

	// Persist the statistics
	if len(statStats) > 0 {
		if err := a.statStatisticsRepository.StoreManyStatStatistics(ctx, statStats); err != nil {
			return nil, fmt.Errorf("failed to store stat statistics: %w", err)
		}
	}

	result.BuildsProcessed = int32(totalProcessed)
	result.ItemsProcessed = int32(len(statStats))

	logger.Info("Completed stat statistics processing",
		"class", class,
		"spec", spec,
		"encounter", encounterID,
		"buildsProcessed", totalProcessed,
		"statsProcessed", len(statStats))

	return result, nil
}

// Structure to aggregate statistics
type StatAggregation struct {
	Category   string
	TotalValue float64
	MinValue   float64
	MaxValue   float64
	Count      int

	// For correlations with levels
	TotalItemLevel     float64
	MinItemLevel       float64
	MaxItemLevel       float64
	TotalKeystoneLevel float64
	MinKeystoneLevel   int
	MaxKeystoneLevel   int
}

// countBuilds count the number of builds for a specific class, spec and encounter_id
func (a *StatStatisticsActivity) countBuilds(
	ctx context.Context,
	class, spec string,
	encounterID uint,
) (int64, error) {
	return a.playerBuildsRepository.CountPlayerBuildsByFilter(ctx, class, spec, encounterID)
}

// getPlayerBuildsBatch get a batch of builds for a specific class, spec and encounter_id with pagination
func (a *StatStatisticsActivity) getPlayerBuildsBatch(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	limit, offset int,
) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	return a.playerBuildsRepository.GetPlayerBuildsByFilter(ctx, class, spec, encounterID, limit, offset)
}

// ProcessStatsBatch process a batch of builds to extract the stat statistics and aggregate them
func (a *StatStatisticsActivity) ProcessStatsBatch(
	builds []*warcraftlogsBuilds.PlayerBuild,
	statData map[string]*StatAggregation,
) error {
	if len(builds) == 0 {
		return nil
	}

	var mutex sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(builds))

	// Process each build in parallel
	for _, build := range builds {
		wg.Add(1)
		go func(b *warcraftlogsBuilds.PlayerBuild) {
			defer wg.Done()

			// Check that the stats field is not empty
			if len(b.Stats) == 0 {
				return
			}

			// Parse the stats JSON
			var statsMap map[string]Stat
			if err := json.Unmarshal([]byte(b.Stats), &statsMap); err != nil {
				errCh <- fmt.Errorf("error parsing stats JSON for build %d: %w", b.ID, err)
				return
			}

			// For each stat, check its category and aggregate the data
			mutex.Lock()
			defer mutex.Unlock()

			for statName, statValue := range statsMap {
				// Determine the category
				var category string
				if secondaryStats[statName] {
					category = "secondary"
				} else if minorStats[statName] {
					category = "minor"
				} else {
					// Ignore stats that are neither secondary nor minor
					continue
				}

				// Use the average value of min/max
				value := (statValue.Min + statValue.Max) / 2

				// Create or update the aggregation
				agg, exists := statData[statName]
				if !exists {
					agg = &StatAggregation{
						Category:         category,
						MinValue:         value,
						MaxValue:         value,
						MinItemLevel:     b.ItemLevel,
						MaxItemLevel:     b.ItemLevel,
						MinKeystoneLevel: b.KeystoneLevel,
						MaxKeystoneLevel: b.KeystoneLevel,
					}
					statData[statName] = agg
				}

				// Update the values
				agg.TotalValue += value
				agg.Count++

				if value < agg.MinValue {
					agg.MinValue = value
				}
				if value > agg.MaxValue {
					agg.MaxValue = value
				}

				// Update the correlations
				agg.TotalItemLevel += b.ItemLevel
				if b.ItemLevel < agg.MinItemLevel {
					agg.MinItemLevel = b.ItemLevel
				}
				if b.ItemLevel > agg.MaxItemLevel {
					agg.MaxItemLevel = b.ItemLevel
				}

				agg.TotalKeystoneLevel += float64(b.KeystoneLevel)
				if b.KeystoneLevel < agg.MinKeystoneLevel {
					agg.MinKeystoneLevel = b.KeystoneLevel
				}
				if b.KeystoneLevel > agg.MaxKeystoneLevel {
					agg.MaxKeystoneLevel = b.KeystoneLevel
				}
			}
		}(build)
	}

	// Wait for all the goroutines to finish
	wg.Wait()
	close(errCh)

	// Check if there were any errors
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

// ConvertToStatStatistics convert the aggregated data to StatStatistic objects
func (a *StatStatisticsActivity) ConvertToStatStatistics(
	statData map[string]*StatAggregation,
	class, spec string,
	encounterID uint,
) []*warcraftlogsBuilds.StatStatistic {
	result := make([]*warcraftlogsBuilds.StatStatistic, 0, len(statData))

	for statName, agg := range statData {
		if agg.Count == 0 {
			continue
		}

		stat := &warcraftlogsBuilds.StatStatistic{
			Class:            class,
			Spec:             spec,
			EncounterID:      encounterID,
			StatName:         statName,
			StatCategory:     agg.Category,
			AvgValue:         agg.TotalValue / float64(agg.Count),
			MinValue:         agg.MinValue,
			MaxValue:         agg.MaxValue,
			SampleSize:       agg.Count,
			AvgItemLevel:     agg.TotalItemLevel / float64(agg.Count),
			MinItemLevel:     agg.MinItemLevel,
			MaxItemLevel:     agg.MaxItemLevel,
			AvgKeystoneLevel: agg.TotalKeystoneLevel / float64(agg.Count),
			MinKeystoneLevel: agg.MinKeystoneLevel,
			MaxKeystoneLevel: agg.MaxKeystoneLevel,
		}

		result = append(result, stat)
	}

	return result
}
