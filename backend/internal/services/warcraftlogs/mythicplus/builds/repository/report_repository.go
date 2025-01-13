package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
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
func (r *ReportRepository) StoreReport(ctx context.Context, report *warcraftlogsBuilds.Report) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		log.Printf("Checking for existing report: %s (FightID: %d)", report.Code, report.FightID)

		// Check if report exists
		var existingReport warcraftlogsBuilds.Report
		result := tx.WithContext(ctx).
			Where("code = ? AND fight_id = ?", report.Code, report.FightID).
			First(&existingReport)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check existing report: %w", result.Error)
		}

		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("Report not found, creating new report: %s (FightID: %d)", report.Code, report.FightID)
			// Create new report
			if err := tx.WithContext(ctx).Create(report).Error; err != nil {
				return fmt.Errorf("failed to create report: %w", err)
			}
		} else {
			log.Printf("Report found, updating existing report: %s (FightID: %d)", report.Code, report.FightID)
			// Update existing report
			if err := tx.WithContext(ctx).Model(&existingReport).Updates(report).Error; err != nil {
				return fmt.Errorf("failed to update report: %w", err)
			}
			log.Printf("Report updated successfully: %s (FightID: %d)", report.Code, report.FightID)
		}

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

// GetReportsforEncounter retrieves all reports for a given encounter
func (r *ReportRepository) GetReportsForEncounter(ctx context.Context, encounterID uint) ([]warcraftlogsBuilds.Report, error) {
	var reports []warcraftlogsBuilds.Report
	err := r.db.WithContext(ctx).
		Where("encounter_id = ? AND deleted_at IS NULL", encounterID).
		Find(&reports).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
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
