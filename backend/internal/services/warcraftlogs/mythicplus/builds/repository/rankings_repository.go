package warcraftlogsBuildsRepository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

type RankingsRepository struct {
	db *gorm.DB
}

func NewRankingsRepository(db *gorm.DB) *RankingsRepository {
	return &RankingsRepository{
		db: db,
	}
}

// GetLastRankingsForEncounter retrieves the last rankings for a given encounter
func (r *RankingsRepository) GetLastRankingForEncounter(ctx context.Context, encounterID uint, className string, specName string) (*warcraftlogsBuilds.ClassRanking, error) {
	var ranking warcraftlogsBuilds.ClassRanking

	result := r.db.Where("encounter_id = ? AND class = ? AND spec = ?",
		encounterID,
		className,
		specName,
	).Order("updated_at DESC").First(&ranking)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last ranking: %w", result.Error)
	}

	return &ranking, nil
}

// StoreRankings saves a list of class rankings to the database on a transaction
func (r *RankingsRepository) StoreRankings(ctx context.Context, encounterID uint, newRankings []*warcraftlogsBuilds.ClassRanking) error {
	if len(newRankings) > 150 {
		return fmt.Errorf("too many rankings provided: got %d, maximum allowed is 150", len(newRankings))
	}

	log.Printf("[INFO] Storing rankings for encounter %d: %d new rankings to process", encounterID, len(newRankings))

	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingRankings []*warcraftlogsBuilds.ClassRanking
		if err := tx.WithContext(ctx).
			Where("encounter_id = ?", encounterID).
			Find(&existingRankings).Error; err != nil {
			return fmt.Errorf("failed to fetch existing rankings: %w", err)
		}

		log.Printf("[DEBUG] Found %d existing rankings for encounter %d", len(existingRankings), encounterID)

		existingMap := make(map[string]*warcraftlogsBuilds.ClassRanking, len(existingRankings))
		for _, rank := range existingRankings {
			existingMap[buildRankingKey(rank)] = rank
		}

		// Process new rankings
		insertCount := 0
		for _, newRank := range newRankings {
			key := buildRankingKey(newRank)
			if _, exists := existingMap[key]; !exists {
				if err := tx.Create(newRank).Error; err != nil {
					return fmt.Errorf("failed to insert new ranking: %w", err)
				}
				insertCount++
			}
			delete(existingMap, key)
		}

		log.Printf("[DEBUG] Inserted %d new rankings for encounter %d", insertCount, encounterID)

		// Delete old rankings
		if len(existingMap) > 0 {
			idsToDelete := make([]uint, 0, len(existingMap))
			for _, rank := range existingMap {
				idsToDelete = append(idsToDelete, rank.ID)
			}

			result := tx.Unscoped().
				Where("id IN (?)", idsToDelete).
				Delete(&warcraftlogsBuilds.ClassRanking{})

			if result.Error != nil {
				return fmt.Errorf("failed to delete old rankings: %w", result.Error)
			}

			log.Printf("[DEBUG] Deleted %d old rankings for encounter %d", result.RowsAffected, encounterID)
		}

		// Verify final count
		var finalCount int64
		if err := tx.Model(&warcraftlogsBuilds.ClassRanking{}).
			Where("encounter_id = ?", encounterID).
			Count(&finalCount).Error; err != nil {
			return fmt.Errorf("failed to count rankings: %w", err)
		}

		log.Printf("[INFO] Rankings update completed for encounter %d: %d total rankings stored",
			encounterID, finalCount)

		if finalCount > 150 {
			return fmt.Errorf("too many rankings after update for encounter %d: got %d, maximum allowed is 150",
				encounterID, finalCount)
		}

		return nil
	})
}

func buildRankingKey(r *warcraftlogsBuilds.ClassRanking) string {
	return fmt.Sprintf("%s_%s_%s_%d",
		r.PlayerName,
		r.ServerName,
		r.ServerRegion,
		r.ReportFightID)
}

// GetRankingsForAnalysis retrieves the rankings for a given encounter and dungeon
func (r *RankingsRepository) GetRankingsForAnalysis(ctx context.Context, encounterID uint, limit int) ([]*warcraftlogsBuilds.ClassRanking, error) {
	var rankings []*warcraftlogsBuilds.ClassRanking
	err := r.db.WithContext(ctx).
		Where("encounter_id = ?", encounterID).
		Order("created_at DESC").
		Limit(limit).
		Find(&rankings).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get rankings: %w", err)
	}

	return rankings, nil
}

// GetRankingsForSpec retrieves rankings from the database for a specific class/spec/encounter
func (r *RankingsRepository) GetRankingsForSpec(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
	var rankings []*warcraftlogsBuilds.ClassRanking

	result := r.db.WithContext(ctx).
		Where("class = ? AND spec = ? AND encounter_id = ?",
			className, specName, encounterID).
		Find(&rankings)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get rankings: %w", result.Error)
	}

	return rankings, nil
}
