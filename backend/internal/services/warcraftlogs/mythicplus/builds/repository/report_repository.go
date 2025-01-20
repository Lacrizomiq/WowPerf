package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

// StoreReports saves multiple reports to the database in a single transaction
func (r *ReportRepository) StoreReports(ctx context.Context, newReports []*warcraftlogsBuilds.Report) error {
	if len(newReports) == 0 {
		log.Printf("[DEBUG] No reports to store")
		return nil
	}

	encounterID := newReports[0].EncounterID

	return r.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("[DEBUG] Starting transaction for encounter %d with %d reports", encounterID, len(newReports))

		// Fetch existing reports for the encounter
		var existingReports []*warcraftlogsBuilds.Report
		if err := tx.WithContext(ctx).
			Where("encounter_id = ?", encounterID).
			Find(&existingReports).Error; err != nil {
			return fmt.Errorf("failed to fetch existing reports: %w", err)
		}

		// Create a map of existing reports for efficient lookup
		existingMap := make(map[string]*warcraftlogsBuilds.Report)
		for _, report := range existingReports {
			key := fmt.Sprintf("%s_%d", report.Code, report.FightID)
			existingMap[key] = report
		}

		// Process new reports
		for _, newReport := range newReports {
			key := fmt.Sprintf("%s_%d", newReport.Code, newReport.FightID)
			delete(existingMap, key) // Remove from map as we keep it

			// Create or update the report
			result := tx.WithContext(ctx).
				Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "code"},
						{Name: "fight_id"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"encounter_id",
						"total_time",
						"item_level",
						"composition",
						"damage_done",
						"healing_done",
						"damage_taken",
						"death_events",
						"player_details_dps",
						"player_details_healers",
						"player_details_tanks",
						"log_version",
						"game_version",
						"keystonelevel",
						"keystonetime",
						"affixes",
						"friendly_players",
						"talent_codes",
						"raw_data",
						"updated_at",
					}),
				}).Create(newReport)

			if result.Error != nil {
				return fmt.Errorf("failed to store report: %w", result.Error)
			}

			log.Printf("[DEBUG] Stored report: %s (FightID: %d)", newReport.Code, newReport.FightID)
		}

		// Delete obsolete reports
		for _, oldReport := range existingMap {
			if err := tx.WithContext(ctx).Delete(oldReport).Error; err != nil {
				return fmt.Errorf("failed to delete old report: %w", err)
			}
			log.Printf("[DEBUG] Deleted obsolete report: %s (FightID: %d)", oldReport.Code, oldReport.FightID)
		}

		log.Printf("[INFO] Successfully processed reports for encounter %d: stored %d, deleted %d",
			encounterID, len(newReports), len(existingMap))
		return nil
	})
}

// GetReportByCodeAndFightID retrieves a single report by its code and fight ID
func (r *ReportRepository) GetReportByCodeAndFightID(ctx context.Context, code string, fightID int) (*warcraftlogsBuilds.Report, error) {
	var report warcraftlogsBuilds.Report

	result := r.db.WithContext(ctx).
		Where("code = ? AND fight_id = ?", code, fightID).
		First(&report)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get report: %w", result.Error)
	}

	return &report, nil
}

// GetReportsForEncounter retrieves reports for an encounter with pagination and field selection
func (r *ReportRepository) GetReportsForEncounter(ctx context.Context, encounterID uint, limit int, offset int) ([]warcraftlogsBuilds.Report, error) {
	var reports []warcraftlogsBuilds.Report

	// Select only necessary fields to optimize data transfer
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
		Where("encounter_id = ? AND deleted_at IS NULL", encounterID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get reports for encounter: %w", result.Error)
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

// GetReportsByCode retrieves multiple reports by their codes
func (r *ReportRepository) GetReportsByCode(ctx context.Context, codes []string) ([]*warcraftlogsBuilds.Report, error) {
	var reports []*warcraftlogsBuilds.Report

	result := r.db.WithContext(ctx).
		Where("code IN (?)", codes).
		Find(&reports)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get reports by codes: %w", result.Error)
	}

	return reports, nil
}
