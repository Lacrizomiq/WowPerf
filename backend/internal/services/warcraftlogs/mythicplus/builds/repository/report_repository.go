package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReportRepository handles all database operations for reports
type ReportRepository struct {
	db    *gorm.DB
	cache sync.Map
}

// NewReportRepository creates a new instance of ReportRepository
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

// StoreReports saves multiple reports to the database in batches
// It uses UPSERT to handle conflicts when a report already exists
func (r *ReportRepository) StoreReports(ctx context.Context, reports []*warcraftlogsBuilds.Report) error {
	if len(reports) == 0 {
		log.Printf("[DEBUG] No reports to store")
		return nil
	}

	// Process reports in batches to avoid memory issues
	const batchSize = 5
	for i := 0; i < len(reports); i += batchSize {
		end := i + batchSize
		if end > len(reports) {
			end = len(reports)
		}

		batch := reports[i:end]
		if err := r.processBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i, end, err)
		}

		log.Printf("[DEBUG] Processed reports batch %d-%d of %d", i, end, len(reports))

		// Small delay between batches to prevent database overload
		time.Sleep(time.Millisecond * 100)
	}

	log.Printf("[INFO] Successfully processed all %d reports", len(reports))
	return nil
}

// processBatch handles a single batch of reports within a transaction
func (r *ReportRepository) processBatch(ctx context.Context, batch []*warcraftlogsBuilds.Report) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, report := range batch {
			// Set timestamps
			now := time.Now()
			report.CreatedAt = now
			report.UpdatedAt = now

			// Create or update the report using UPSERT
			result := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "code"},
					{Name: "fight_id"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"encounter_id",
					"total_time",
					"item_level",
					"composition",
					"player_details_dps",
					"player_details_healers",
					"player_details_tanks",
					"talent_codes",
					"keystonelevel",
					"affixes",
					"updated_at",
				}),
			}).Create(report)

			if result.Error != nil {
				return fmt.Errorf("failed to store report %s (FightID: %d): %w",
					report.Code, report.FightID, result.Error)
			}

			log.Printf("[TRACE] Stored report %s (FightID: %d)",
				report.Code, report.FightID)
		}
		return nil
	})
}

// GetReportsByRankings retrieves reports corresponding to the provided rankings
// This method is used by the workflow to get existing reports
func (r *ReportRepository) GetReportsByRankings(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) ([]*warcraftlogsBuilds.Report, error) {
	if len(rankings) == 0 {
		return nil, nil
	}

	// Extract unique report codes
	codes := make([]string, 0, len(rankings))
	codeMap := make(map[string]bool)
	for _, ranking := range rankings {
		if !codeMap[ranking.ReportCode] {
			codes = append(codes, ranking.ReportCode)
			codeMap[ranking.ReportCode] = true
		}
	}

	var reports []*warcraftlogsBuilds.Report
	result := r.db.WithContext(ctx).
		Select(
			"code",
			"fight_id",
			"encounter_id",
			"player_details_dps",
			"player_details_healers",
			"player_details_tanks",
			"talent_codes",
			"keystonelevel",
			"affixes",
		).
		Where("code IN (?)", codes).
		Find(&reports)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get reports by rankings: %w", result.Error)
	}

	return reports, nil
}

// CountReportsForEncounter counts the total number of reports for an encounter
func (r *ReportRepository) CountReportsForEncounter(ctx context.Context, encounterID uint) (int64, error) {
	var count int64

	if err := r.db.Model(&warcraftlogsBuilds.Report{}).
		Where("encounter_id = ? AND deleted_at IS NULL", encounterID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count reports for encounter: %w", err)
	}

	return count, nil
}

// GetReportByCodeAndFightID retrieves a single report by its code and fight ID
func (r *ReportRepository) GetReportByCodeAndFightID(ctx context.Context, code string, fightID int) (*warcraftlogsBuilds.Report, error) {
	// Check cache first
	key := fmt.Sprintf("%s-%d", code, fightID)
	if report, ok := r.cache.Load(key); ok {
		return report.(*warcraftlogsBuilds.Report), nil
	}

	// Query database if not in cache
	var report warcraftlogsBuilds.Report
	result := r.db.WithContext(ctx).
		Where("code = ? AND fight_id = ?", code, fightID).
		First(&report)

	if result.Error == nil {
		// Store in cache
		r.cache.Store(key, &report)
	}

	return &report, nil
}

// GetReportsBatch retrieves a batch of reports with pagination
func (r *ReportRepository) GetReportsBatch(ctx context.Context, limit int, offset int) ([]*warcraftlogsBuilds.Report, error) {
	var reports []*warcraftlogsBuilds.Report

	result := r.db.WithContext(ctx).
		Select(
			"code",
			"fight_id",
			"encounter_id",
			"player_details_dps",
			"player_details_healers",
			"player_details_tanks",
			"talent_codes",
			"keystonelevel",
			"affixes",
		).
		Where("deleted_at IS NULL").
		Order("code ASC, fight_id ASC").
		Limit(limit).
		Offset(offset).
		Find(&reports)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get reports batch: %w", result.Error)
	}

	return reports, nil
}

// CountAllReports returns the total number of reports available
func (r *ReportRepository) CountAllReports(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&warcraftlogsBuilds.Report{}).
		Where("deleted_at IS NULL").
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count reports: %w", result.Error)
	}

	return count, nil
}
