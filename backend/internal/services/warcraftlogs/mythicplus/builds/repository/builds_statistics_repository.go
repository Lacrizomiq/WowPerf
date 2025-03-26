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

// BuildsStatisticsRepository handles all database operations for builds statistics.
type BuildsStatisticsRepository struct {
	db *gorm.DB
}

// NewBuildsStatisticsRepository creates a new instance of BuildsStatisticsRepository.
func NewBuildsStatisticsRepository(db *gorm.DB) *BuildsStatisticsRepository {
	return &BuildsStatisticsRepository{
		db: db,
	}
}

// DeleteBuildStatistics removes build statistics for a class and spec.
func (r *BuildsStatisticsRepository) DeleteBuildStatistics(ctx context.Context, class, spec string, encounterID uint) error {
	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	result := query.Delete(&warcraftlogsBuilds.BuildStatistic{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete build statistics for %s %s: %w", class, spec, result.Error)
	}

	log.Printf("[INFO] Deleted %d existing build statistics for %s-%s (encounterID: %d)",
		result.RowsAffected, class, spec, encounterID)
	return nil
}

// GetBuildStatistics retrieves build statistics from the database based on filter criteria.
func (r *BuildsStatisticsRepository) GetBuildStatistics(ctx context.Context, class, spec string, encounterID uint) ([]*warcraftlogsBuilds.BuildStatistic, error) {
	var stats []*warcraftlogsBuilds.BuildStatistic

	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch build statistics: %w", err)
	}

	log.Printf("[INFO] Retrieved %d build statistics for %s-%s (encounterID: %d)",
		len(stats), class, spec, encounterID)
	return stats, nil
}

// StoreManyBuildStatistics persists multiple build statistics to the database.
// It handles batching to avoid memory issues and uses UPSERT for conflict resolution.
func (r *BuildsStatisticsRepository) StoreManyBuildStatistics(ctx context.Context, buildStats []*warcraftlogsBuilds.BuildStatistic) error {
	if len(buildStats) == 0 {
		log.Printf("[DEBUG] No build statistics to store")
		return nil
	}

	// Process by batch to avoid memory issues
	const batchSize = 5
	for i := 0; i < len(buildStats); i += batchSize {
		end := i + batchSize
		if end > len(buildStats) {
			end = len(buildStats)
		}

		batch := buildStats[i:end]
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, stat := range batch {
				// Set timestamp to current time
				now := time.Now()
				stat.CreatedAt = now
				stat.UpdatedAt = now

				// Create/update statistic with UPSERT
				result := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "class"},
						{Name: "spec"},
						{Name: "encounter_id"},
						{Name: "item_slot"},
						{Name: "item_id"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"item_name", "item_icon", "item_quality", "item_level",
						"has_set_bonus", "set_id", "bonus_ids",
						"has_gems", "gems_count", "gem_ids", "gem_icons", "gem_levels",
						"has_permanent_enchant", "permanent_enchant_id", "permanent_enchant_name",
						"has_temporary_enchant", "temporary_enchant_id", "temporary_enchant_name",
						"usage_count", "usage_percentage",
						"avg_item_level", "min_item_level", "max_item_level",
						"avg_keystone_level", "min_keystone_level", "max_keystone_level",
						"updated_at",
					}),
				}).Create(stat)

				if result.Error != nil {
					return fmt.Errorf("failed to store build statistic for item %d (slot %d): %w",
						stat.ItemID, stat.ItemSlot, result.Error)
				}

				log.Printf("[TRACE] Stored build statistic for item %d in slot %d",
					stat.ItemID, stat.ItemSlot)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to store build statistics batch %d-%d: %w", i, end, err)
		}

		log.Printf("[DEBUG] Processed batch %d/%d, total build statistics processed: %d",
			(i/batchSize)+1, (len(buildStats)+batchSize-1)/batchSize, end)

		// Small delay between batches to avoid overload
		time.Sleep(time.Millisecond * 100)
	}

	log.Printf("[INFO] Successfully processed all build statistics: stored %d statistics", len(buildStats))
	return nil
}

// CountBuildStatistics returns the total count of build statistics in the database.
func (r *BuildsStatisticsRepository) CountBuildStatistics(ctx context.Context, class, spec string, encounterID uint) (int64, error) {
	var count int64

	query := r.db.Model(&warcraftlogsBuilds.BuildStatistic{}).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count build statistics: %w", err)
	}

	return count, nil
}
