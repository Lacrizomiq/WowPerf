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
	BuildsStatisticsRepository handles all database operations for builds statistics.

	Methods:
	- DeleteBuildStatistics: Deletes build statistics for a class and spec.
	- GetBuildStatistics: Retrieves build statistics from the database based on filter criteria.
	- StoreManyBuildStatistics: Persists multiple build statistics to the database.
	- CountBuildStatistics: Returns the total count of build statistics in the database.

*/

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
	// Unscoped() for the hard delete
	query := r.db.WithContext(ctx).
		Unscoped().
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

	// reset of the status in the player_builds table
	resetQuery := r.db.WithContext(ctx).
		Model(&warcraftlogsBuilds.PlayerBuild{}).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		resetQuery = resetQuery.Where("encounter_id = ?", encounterID)
	}

	resetResult := resetQuery.
		Where("equipment_status = 'processed'").
		Updates(map[string]interface{}{
			"equipment_status":       "pending",
			"equipment_processed_at": nil,
			"updated_at":             time.Now(),
		})

	if resetResult.Error != nil {
		return fmt.Errorf("failed to reset equipment status for builds %s-%s: %w", class, spec, resetResult.Error)
	}

	log.Printf("[INFO] Reset equipment status to 'pending' for %d builds of %s-%s (encounterID: %d)",
		resetResult.RowsAffected, class, spec, encounterID)

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

				// Create/update statistic
				result := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "class"},
						{Name: "spec"},
						{Name: "encounter_id"},
						{Name: "item_slot"},
						{Name: "item_id"},
					},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"item_name":              stat.ItemName,
						"item_icon":              stat.ItemIcon,
						"item_quality":           stat.ItemQuality,
						"item_level":             stat.ItemLevel,
						"has_set_bonus":          stat.HasSetBonus,
						"set_id":                 stat.SetID,
						"bonus_ids":              stat.BonusIDs,
						"has_gems":               stat.HasGems,
						"gems_count":             stat.GemsCount,
						"gem_ids":                stat.GemIDs,
						"gem_icons":              stat.GemIcons,
						"gem_levels":             stat.GemLevels,
						"has_permanent_enchant":  stat.HasPermanentEnchant,
						"permanent_enchant_id":   stat.PermanentEnchantID,
						"permanent_enchant_name": stat.PermanentEnchantName,
						"has_temporary_enchant":  stat.HasTemporaryEnchant,
						"temporary_enchant_id":   stat.TemporaryEnchantID,
						"temporary_enchant_name": stat.TemporaryEnchantName,
						"usage_count":            gorm.Expr("build_statistics.usage_count + ?", stat.UsageCount),
						"usage_percentage":       stat.UsagePercentage,
						"avg_item_level": gorm.Expr(
							"(CAST(build_statistics.avg_item_level AS numeric) * CAST(build_statistics.usage_count AS numeric) + CAST(? AS numeric) * CAST(? AS numeric)) / CAST((build_statistics.usage_count + ?) AS numeric)",
							stat.AvgItemLevel, stat.UsageCount, stat.UsageCount),
						"min_item_level": gorm.Expr("LEAST(build_statistics.min_item_level, ?)", stat.MinItemLevel),
						"max_item_level": gorm.Expr("GREATEST(build_statistics.max_item_level, ?)", stat.MaxItemLevel),
						"avg_keystone_level": gorm.Expr(
							"(CAST(build_statistics.avg_keystone_level AS numeric) * CAST(build_statistics.usage_count AS numeric) + CAST(? AS numeric) * CAST(? AS numeric)) / CAST((build_statistics.usage_count + ?) AS numeric)",
							stat.AvgKeystoneLevel, stat.UsageCount, stat.UsageCount),
						"min_keystone_level": gorm.Expr("LEAST(build_statistics.min_keystone_level, ?)", stat.MinKeystoneLevel),
						"max_keystone_level": gorm.Expr("GREATEST(build_statistics.max_keystone_level, ?)", stat.MaxKeystoneLevel),
						"updated_at":         time.Now(),
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
