package warcraftlogsBuildsRepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

type RankingsRepository struct {
	db                 *gorm.DB
	maxRankingsPerSpec int
}

func NewRankingsRepository(db *gorm.DB, maxRankingsPerSpec int) *RankingsRepository {
	return &RankingsRepository{
		db:                 db,
		maxRankingsPerSpec: maxRankingsPerSpec,
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
	if len(newRankings) == 0 || newRankings[0] == nil {
		return nil
	}

	// Extract class and spec information from the first ranking
	className := newRankings[0].Class
	specName := newRankings[0].Spec

	if len(newRankings) > r.maxRankingsPerSpec {
		return fmt.Errorf("too many rankings provided: got %d, maximum allowed is %d", len(newRankings), r.maxRankingsPerSpec)
	}

	log.Printf("[INFO] Storing rankings for encounter %d, class %s, spec %s: %d new rankings to process",
		encounterID, className, specName, len(newRankings))

	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingRankings []*warcraftlogsBuilds.ClassRanking
		if err := tx.WithContext(ctx).
			Where("encounter_id = ? AND class = ? AND spec = ?", encounterID, className, specName).
			Find(&existingRankings).Error; err != nil {
			return fmt.Errorf("failed to fetch existing rankings: %w", err)
		}

		log.Printf("[DEBUG] Found %d existing rankings for encounter %d, class %s, spec %s",
			len(existingRankings), encounterID, className, specName)

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

		if finalCount > int64(r.maxRankingsPerSpec) {
			return fmt.Errorf("too many rankings after update for encounter %d: got %d, maximum allowed is %d",
				encounterID, finalCount, r.maxRankingsPerSpec)
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

// MarkRankingsAsProcessedForReports marks the rankings as ready for report processing
func (r *RankingsRepository) MarkRankingsAsProcessedForReports(ctx context.Context, ids []uint, batchID string) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&warcraftlogsBuilds.ClassRanking{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"report_processing_status": "processed",
			"processing_batch_id":      batchID,
			"report_processing_at":     time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to mark rankings as pending for reports: %w", result.Error)
	}

	return nil
}

// GetRankingsNeedingReportProcessing retrieves the rankings that need report processing
// It retrieves the rankings that need report processing and are older than the max age
// It also filters by class if a class is provided
func (r *RankingsRepository) GetRankingsNeedingReportProcessing(ctx context.Context, className string, limit int, maxAge time.Duration) ([]*warcraftlogsBuilds.ClassRanking, error) {
	var rankings []*warcraftlogsBuilds.ClassRanking
	minDate := time.Now().Add(-maxAge)

	query := r.db.WithContext(ctx).
		Where("(report_processing_status = ? OR (report_processing_status = ? AND report_processing_at < ?)) AND created_at > ?",
			"pending", "failed", minDate, minDate)

	// Add the class filter
	if className != "" {
		query = query.Where("class = ?", className)
	}

	result := query.Order("created_at DESC").
		Limit(limit).
		Find(&rankings)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []*warcraftlogsBuilds.ClassRanking{}, nil
		}
		return nil, fmt.Errorf("failed to get rankings needing report processing: %w", result.Error)
	}

	return rankings, nil
}
