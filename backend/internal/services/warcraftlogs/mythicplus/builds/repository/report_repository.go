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

// ReportIdentifier represents the unique identifier for a report
type ReportIdentifier struct {
	Code    string
	FightID int
}

// ReportRepository handles all database operations for reports
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new instance of ReportRepository
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

// StoreReports optimizes batch processing of reports
func (r *ReportRepository) StoreReports(ctx context.Context, reports []*warcraftlogsBuilds.Report) error {
	if len(reports) == 0 {
		log.Printf("[DEBUG] No reports to store")
		return nil
	}

	// Increased batch size for better performance
	const batchSize = 20

	// Process reports in larger batches
	for i := 0; i < len(reports); i += batchSize {
		end := i + batchSize
		if end > len(reports) {
			end = len(reports)
		}

		batch := reports[i:end]
		if err := r.processBatchBulk(ctx, batch); err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i, end, err)
		}

		log.Printf("[DEBUG] Processed reports batch %d-%d of %d", i, end, len(reports))

		// Reduced delay between batches
		time.Sleep(time.Millisecond * 50)
	}

	log.Printf("[INFO] Successfully processed all %d reports", len(reports))
	return nil
}

// processBatchBulk handles bulk insertion of reports
func (r *ReportRepository) processBatchBulk(ctx context.Context, batch []*warcraftlogsBuilds.Report) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Set timestamps for all reports in batch
		for _, report := range batch {
			report.CreatedAt = now
			report.UpdatedAt = now
		}

		// Bulk insert/update using a single query
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
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("failed to bulk store reports: %w", result.Error)
		}

		return nil
	})
}

// SyncReportsWithRankings synchronizes reports with the provided rankings
// It deletes reports that no longer have associated rankings
func (r *ReportRepository) SyncReportsWithRankings(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) error {
	if len(rankings) == 0 {
		return nil
	}

	// Create a map of active report identifiers from rankings
	activeReports := make(map[string]bool)
	for _, ranking := range rankings {
		key := fmt.Sprintf("%s-%d", ranking.ReportCode, ranking.ReportFightID)
		activeReports[key] = true
	}

	// Delete reports that are no longer referenced by any ranking
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reportsToDelete []ReportIdentifier

		// Find all reports for these rankings
		rows, err := tx.Model(&warcraftlogsBuilds.Report{}).
			Select("code, fight_id").
			Where("encounter_id = ?", rankings[0].EncounterID).
			Rows()

		if err != nil {
			return fmt.Errorf("failed to fetch reports: %w", err)
		}
		defer rows.Close()

		// Check each report against active rankings
		for rows.Next() {
			var report ReportIdentifier
			if err := rows.Scan(&report.Code, &report.FightID); err != nil {
				return fmt.Errorf("failed to scan report: %w", err)
			}

			key := fmt.Sprintf("%s-%d", report.Code, report.FightID)
			if !activeReports[key] {
				reportsToDelete = append(reportsToDelete, report)
			}
		}

		// Delete obsolete reports
		if len(reportsToDelete) > 0 {
			for _, report := range reportsToDelete {
				if err := tx.Unscoped().Where("code = ? AND fight_id = ?",
					report.Code, report.FightID).Delete(&warcraftlogsBuilds.Report{}).Error; err != nil {
					return fmt.Errorf("failed to delete report %s-%d: %w",
						report.Code, report.FightID, err)
				}
			}
			log.Printf("[INFO] Deleted %d obsolete reports", len(reportsToDelete))
		}

		return nil
	})
}

// GetReportsByRankings retrieves reports corresponding to the provided rankings
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

// GetReportsBatch retrieves a batch of reports with pagination
// Returns an empty slice if offset exceeds total count
func (r *ReportRepository) GetReportsBatch(ctx context.Context, limit int, offset int) ([]*warcraftlogsBuilds.Report, error) {
	// First check total count to avoid unnecessary queries
	var count int64
	if err := r.db.Model(&warcraftlogsBuilds.Report{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count reports: %w", err)
	}

	// Return empty slice if offset is beyond available data
	if int64(offset) >= count {
		return []*warcraftlogsBuilds.Report{}, nil
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
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count reports: %w", result.Error)
	}

	return count, nil
}
