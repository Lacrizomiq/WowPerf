package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
	StatStatisticsRepository handles database operations for character stat statistics.

	Methods:
	- DeleteStatStatistics: Deletes stat statistics for a class and spec.
	- GetStatStatistics: Retrieves stat statistics from the database based on filter criteria.
	- StoreManyStatStatistics: Persists multiple stat statistics to the database.
	- CountStatStatistics: Returns the total count of stat statistics in the database.
	- GetStatPriorities: Returns stat statistics sorted by average value (highest first) for a specific category.
*/

// StatStatisticsRepository handles database operations for character stat statistics.
type StatStatisticsRepository struct {
	db *gorm.DB
}

// NewStatStatisticsRepository creates a new instance of StatStatisticsRepository.
func NewStatStatisticsRepository(db *gorm.DB) *StatStatisticsRepository {
	return &StatStatisticsRepository{
		db: db,
	}
}

// DeleteStatStatistics removes stat statistics for a class and spec.
func (r *StatStatisticsRepository) DeleteStatStatistics(ctx context.Context, class, spec string, encounterID uint) error {
	// Unscoped() for the hard delete
	query := r.db.WithContext(ctx).
		Unscoped().
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	result := query.Delete(&warcraftlogsBuilds.StatStatistic{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete stat statistics for %s-%s: %w", class, spec, result.Error)
	}

	log.Printf("[INFO] Deleted %d existing stat statistics for %s-%s (encounterID: %d)",
		result.RowsAffected, class, spec, encounterID)

	// reset of the status in the player_builds table
	resetQuery := r.db.WithContext(ctx).
		Model(&warcraftlogsBuilds.PlayerBuild{}).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		resetQuery = resetQuery.Where("encounter_id = ?", encounterID)
	}

	resetResult := resetQuery.
		Where("stat_status = 'processed'").
		Updates(map[string]interface{}{
			"stat_status":       "pending",
			"stat_processed_at": nil,
			"updated_at":        time.Now(),
		})

	if resetResult.Error != nil {
		return fmt.Errorf("failed to reset stats status for builds %s-%s: %w", class, spec, resetResult.Error)
	}

	log.Printf("[INFO] Reset stats status to 'pending' for %d builds of %s-%s (encounterID: %d)",
		resetResult.RowsAffected, class, spec, encounterID)

	return nil
}

// GetStatStatistics retrieves stat statistics from the database based on filter criteria.
func (r *StatStatisticsRepository) GetStatStatistics(ctx context.Context, class, spec string, encounterID uint) ([]*warcraftlogsBuilds.StatStatistic, error) {
	var stats []*warcraftlogsBuilds.StatStatistic

	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch stat statistics: %w", err)
	}

	log.Printf("[INFO] Retrieved %d stat statistics for %s-%s (encounterID: %d)",
		len(stats), class, spec, encounterID)
	return stats, nil
}

// GetStatStatisticsByCategory retrieves stat statistics for a specific category.
func (r *StatStatisticsRepository) GetStatStatisticsByCategory(ctx context.Context, class, spec string, encounterID uint, category string) ([]*warcraftlogsBuilds.StatStatistic, error) {
	var stats []*warcraftlogsBuilds.StatStatistic

	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ? AND stat_category = ?", class, spec, category)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch %s stat statistics: %w", category, err)
	}

	log.Printf("[INFO] Retrieved %d %s stat statistics for %s-%s (encounterID: %d)",
		len(stats), category, class, spec, encounterID)
	return stats, nil
}

// StoreManyStatStatistics persists multiple stat statistics to the database.
// It handles batching to avoid memory issues and uses UPSERT for conflict resolution.
func (r *StatStatisticsRepository) StoreManyStatStatistics(ctx context.Context, statStats []*warcraftlogsBuilds.StatStatistic) error {
	if len(statStats) == 0 {
		log.Printf("[DEBUG] No stat statistics to store")
		return nil
	}

	// Process by batch to avoid memory issues
	const batchSize = 5
	for i := 0; i < len(statStats); i += batchSize {
		end := i + batchSize
		if end > len(statStats) {
			end = len(statStats)
		}

		batch := statStats[i:end]
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, stat := range batch {
				// Set timestamps
				now := time.Now()
				stat.CreatedAt = now
				stat.UpdatedAt = now

				// Create/update statistic with UPSERT
				result := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "class"},
						{Name: "spec"},
						{Name: "encounter_id"},
						{Name: "stat_name"},
					},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"avg_value": gorm.Expr(
							"(CAST(stat_statistics.avg_value AS numeric) * CAST(stat_statistics.sample_size AS numeric) + CAST(? AS numeric) * CAST(? AS numeric)) / CAST((stat_statistics.sample_size + ?) AS numeric)",
							stat.AvgValue, stat.SampleSize, stat.SampleSize),
						"min_value":   gorm.Expr("LEAST(stat_statistics.min_value, ?)", stat.MinValue),
						"max_value":   gorm.Expr("GREATEST(stat_statistics.max_value, ?)", stat.MaxValue),
						"sample_size": gorm.Expr("stat_statistics.sample_size + ?", stat.SampleSize),
						"avg_keystone_level": gorm.Expr(
							"(CAST(stat_statistics.avg_keystone_level AS numeric) * CAST(stat_statistics.sample_size AS numeric) + CAST(? AS numeric) * CAST(? AS numeric)) / CAST((stat_statistics.sample_size + ?) AS numeric)",
							stat.AvgKeystoneLevel, stat.SampleSize, stat.SampleSize),
						"min_keystone_level": gorm.Expr("LEAST(stat_statistics.min_keystone_level, ?)", stat.MinKeystoneLevel),
						"max_keystone_level": gorm.Expr("GREATEST(stat_statistics.max_keystone_level, ?)", stat.MaxKeystoneLevel),
						"avg_item_level": gorm.Expr(
							"(CAST(stat_statistics.avg_item_level AS numeric) * CAST(stat_statistics.sample_size AS numeric) + CAST(? AS numeric) * CAST(? AS numeric)) / CAST((stat_statistics.sample_size + ?) AS numeric)",
							stat.AvgItemLevel, stat.SampleSize, stat.SampleSize),
						"min_item_level": gorm.Expr("LEAST(stat_statistics.min_item_level, ?)", stat.MinItemLevel),
						"max_item_level": gorm.Expr("GREATEST(stat_statistics.max_item_level, ?)", stat.MaxItemLevel),
						"updated_at":     time.Now(),
					}),
				}).Create(stat)

				if result.Error != nil {
					return fmt.Errorf("failed to store stat statistic for %s (%s): %w",
						stat.StatName, stat.StatCategory, result.Error)
				}

				log.Printf("[TRACE] Stored stat statistic for %s (%s), avg value: %.2f",
					stat.StatName, stat.StatCategory, stat.AvgValue)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to store stat statistics batch %d-%d: %w", i, end, err)
		}

		log.Printf("[DEBUG] Processed batch %d/%d, total stat statistics processed: %d",
			(i/batchSize)+1, (len(statStats)+batchSize-1)/batchSize, end)

		// Small delay between batches to avoid overload
		time.Sleep(time.Millisecond * 100)
	}

	log.Printf("[INFO] Successfully processed all stat statistics: stored %d statistics", len(statStats))
	return nil
}

// CountStatStatistics returns the total count of stat statistics in the database.
func (r *StatStatisticsRepository) CountStatStatistics(ctx context.Context, class, spec string, encounterID uint) (int64, error) {
	var count int64

	query := r.db.Model(&warcraftlogsBuilds.StatStatistic{}).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count stat statistics: %w", err)
	}

	return count, nil
}

// GetStatPriorities returns stat statistics sorted by average value (highest first)
// for a specific category, typically used to determine stat priorities.
func (r *StatStatisticsRepository) GetStatPriorities(ctx context.Context, class, spec string, encounterID uint, category string) ([]*warcraftlogsBuilds.StatStatistic, error) {
	var stats []*warcraftlogsBuilds.StatStatistic

	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ? AND stat_category = ?", class, spec, category).
		Order("avg_value DESC")

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch stat priorities: %w", err)
	}

	log.Printf("[INFO] Retrieved stat priorities for %s-%s (%s stats, encounterID: %d)",
		class, spec, category, encounterID)
	return stats, nil
}
