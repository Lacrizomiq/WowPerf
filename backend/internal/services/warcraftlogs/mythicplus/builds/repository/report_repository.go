package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

// StoreReport saves a report to the database
func (r *ReportRepository) StoreReports(ctx context.Context, newReports []*warcraftlogsBuilds.Report) error {
	if len(newReports) == 0 {
		log.Printf("[DEBUG] No reports to store")
		return nil
	}

	encounterID := newReports[0].EncounterID

	return r.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("[DEBUG] Starting transaction for encounter %d with %d reports", encounterID, len(newReports))

		// Fetch existing reports
		var existingReports []*warcraftlogsBuilds.Report
		if err := tx.Where("encounter_id = ?", encounterID).Find(&existingReports).Error; err != nil {
			return fmt.Errorf("failed to fetch existing reports: %w", err)
		}

		// Create a map of existing reports
		existingMap := make(map[string]*warcraftlogsBuilds.Report)
		for _, report := range existingReports {
			key := fmt.Sprintf("%s_%d", report.Code, report.FightID)
			existingMap[key] = report
		}

		// Process new reports
		for _, newReport := range newReports {
			key := fmt.Sprintf("%s_%d", newReport.Code, newReport.FightID)
			delete(existingMap, key) // Remove from map as we keep it

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

		// Delete reports that no longer exist in the rankings
		for _, oldReport := range existingMap {
			if err := tx.Delete(oldReport).Error; err != nil {
				return fmt.Errorf("failed to delete old report: %w", err)
			}
			log.Printf("[DEBUG] Deleted obsolete report: %s (FightID: %d)", oldReport.Code, oldReport.FightID)
		}

		log.Printf("[INFO] Successfully processed reports for encounter %d: stored %d, deleted %d",
			encounterID, len(newReports), len(existingMap))
		return nil
	})
}

// GetReportByCodeAndFightID retrieves a report by its code and fight ID
func (r *ReportRepository) GetReportByCodeAndFightID(ctx context.Context, code string, fightID int) (*warcraftlogsBuilds.Report, error) {
	var report warcraftlogsBuilds.Report
	result := r.db.WithContext(ctx).Where("code = ? AND fight_id = ?", code, fightID).First(&report)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get report: %w", result.Error)
	}

	return &report, nil
}

func (r *ReportRepository) GetReportsForEncounter(ctx context.Context, encounterID uint, limit int, offset int) ([]warcraftlogsBuilds.Report, error) {
	var reports []warcraftlogsBuilds.Report

	result := r.db.WithContext(ctx).
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
