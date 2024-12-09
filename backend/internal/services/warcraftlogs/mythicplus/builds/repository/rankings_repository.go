package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
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
func (r *RankingsRepository) GetLastRankingForEncounter(ctx context.Context, encounterID uint) (*warcraftlogsBuilds.ClassRanking, error) {
	var ranking warcraftlogsBuilds.ClassRanking
	result := r.db.WithContext(ctx).
		Where("encounter_id = ?", encounterID).
		Order("created_at DESC").
		First(&ranking)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get last ranking: %w", result.Error)
	}

	return &ranking, nil
}

// StoreRankings saves a list of class rankings to the database on a transaction
func (r *RankingsRepository) StoreRankings(ctx context.Context, encounterID uint, rankings []*warcraftlogsBuilds.ClassRanking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Soft delete existing rankings for the encounter
		if err := tx.WithContext(ctx).
			Model(&warcraftlogsBuilds.ClassRanking{}).
			Where("encounter_id = ? AND deleted_at IS NULL", encounterID).
			Update("deleted_at", time.Now()).
			Error; err != nil {
			return fmt.Errorf("failed to soft delete old rankings: %w", err)
		}

		// Ensure that all the affixes are pq.Int64Array
		for _, ranking := range rankings {
			if ranking.Affixes == nil {
				ranking.Affixes = pq.Int64Array{}
			}
		}

		// Insertion des nouveaux classements
		if err := tx.WithContext(ctx).CreateInBatches(rankings, 100).Error; err != nil {
			return fmt.Errorf("failed to insert new rankings: %w", err)
		}

		return nil
	})
}

// GetRankingsForAnalysis retrieves the rankings for a given encounter and dungeon
func (r *RankingsRepository) GetRankingsForAnalysis(ctx context.Context, encounterID uint, limit int) ([]*warcraftlogsBuilds.ClassRanking, error) {
	var rankings []*warcraftlogsBuilds.ClassRanking
	err := r.db.WithContext(ctx).
		Where("encounter_id = ? AND deleted_at IS NULL", encounterID).
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
