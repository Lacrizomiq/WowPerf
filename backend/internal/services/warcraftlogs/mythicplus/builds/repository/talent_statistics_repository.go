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
	TalentStatisticsRepository handles database operations for talent statistics.

	Methods:
	- DeleteTalentStatistics: Deletes talent statistics for a class and spec.
	- GetTalentStatistics: Retrieves talent statistics from the database based on filter criteria.
	- StoreManyTalentStatistics: Persists multiple talent statistics to the database.
	- CountTalentStatistics: Returns the total count of talent statistics in the database.
	- GetMostPopularTalentImport: Returns the most frequently used talent import for a class/spec.
*/

// TalentStatisticsRepository handles database operations for talent statistics.
type TalentStatisticsRepository struct {
	db *gorm.DB
}

// NewTalentStatisticsRepository creates a new instance of TalentStatisticsRepository.
func NewTalentStatisticsRepository(db *gorm.DB) *TalentStatisticsRepository {
	return &TalentStatisticsRepository{
		db: db,
	}
}

// DeleteTalentStatistics removes talent statistics for a class and spec.
func (r *TalentStatisticsRepository) DeleteTalentStatistics(ctx context.Context, class, spec string, encounterID uint) error {
	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	result := query.Delete(&warcraftlogsBuilds.TalentStatistic{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete talent statistics for %s-%s: %w", class, spec, result.Error)
	}

	log.Printf("[INFO] Deleted %d existing talent statistics for %s-%s (encounterID: %d)",
		result.RowsAffected, class, spec, encounterID)
	return nil
}

// GetTalentStatistics retrieves talent statistics from the database based on filter criteria.
func (r *TalentStatisticsRepository) GetTalentStatistics(ctx context.Context, class, spec string, encounterID uint) ([]*warcraftlogsBuilds.TalentStatistic, error) {
	var stats []*warcraftlogsBuilds.TalentStatistic

	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch talent statistics: %w", err)
	}

	log.Printf("[INFO] Retrieved %d talent statistics for %s-%s (encounterID: %d)",
		len(stats), class, spec, encounterID)
	return stats, nil
}

// StoreManyTalentStatistics persists multiple talent statistics to the database.
// It handles batching to avoid memory issues and uses UPSERT for conflict resolution.
func (r *TalentStatisticsRepository) StoreManyTalentStatistics(ctx context.Context, talentStats []*warcraftlogsBuilds.TalentStatistic) error {
	if len(talentStats) == 0 {
		log.Printf("[DEBUG] No talent statistics to store")
		return nil
	}

	// Process by batch to avoid memory issues
	const batchSize = 5
	for i := 0; i < len(talentStats); i += batchSize {
		end := i + batchSize
		if end > len(talentStats) {
			end = len(talentStats)
		}

		batch := talentStats[i:end]
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
						{Name: "talent_import"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"usage_count",
						"usage_percentage",
						"avg_keystone_level",
						"min_keystone_level",
						"max_keystone_level",
						"avg_item_level",
						"min_item_level",
						"max_item_level",
						"updated_at",
					}),
				}).Create(stat)

				if result.Error != nil {
					return fmt.Errorf("failed to store talent statistic for import '%s': %w",
						stat.TalentImport, result.Error)
				}

				log.Printf("[TRACE] Stored talent statistic for import '%s' (usage: %d)",
					stat.TalentImport, stat.UsageCount)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to store talent statistics batch %d-%d: %w", i, end, err)
		}

		log.Printf("[DEBUG] Processed batch %d/%d, total talent statistics processed: %d",
			(i/batchSize)+1, (len(talentStats)+batchSize-1)/batchSize, end)

		// Small delay between batches to avoid overload
		time.Sleep(time.Millisecond * 100)
	}

	log.Printf("[INFO] Successfully processed all talent statistics: stored %d statistics", len(talentStats))
	return nil
}

// CountTalentStatistics returns the total count of talent statistics in the database.
func (r *TalentStatisticsRepository) CountTalentStatistics(ctx context.Context, class, spec string, encounterID uint) (int64, error) {
	var count int64

	query := r.db.Model(&warcraftlogsBuilds.TalentStatistic{}).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count talent statistics: %w", err)
	}

	return count, nil
}

// GetMostPopularTalentImport returns the most frequently used talent import for a class/spec.
func (r *TalentStatisticsRepository) GetMostPopularTalentImport(ctx context.Context, class, spec string, encounterID uint) (*warcraftlogsBuilds.TalentStatistic, error) {
	var stat warcraftlogsBuilds.TalentStatistic

	query := r.db.WithContext(ctx).
		Where("class = ? AND spec = ?", class, spec)

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Order("usage_count DESC").First(&stat).Error; err != nil {
		return nil, fmt.Errorf("failed to find most popular talent import: %w", err)
	}

	return &stat, nil
}
