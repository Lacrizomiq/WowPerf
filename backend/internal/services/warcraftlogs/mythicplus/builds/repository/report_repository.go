package warcraftlogsBuildsRepository

import (
	"context"
	"errors"
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
	// Deduplicate reports before processing
	uniqueReports := make(map[string]*warcraftlogsBuilds.Report)
	for _, report := range batch {
		key := fmt.Sprintf("%s-%d", report.Code, report.FightID)
		uniqueReports[key] = report
	}

	// Convert map to slice for processing
	deduplicatedBatch := make([]*warcraftlogsBuilds.Report, 0, len(uniqueReports))
	for _, report := range uniqueReports {
		deduplicatedBatch = append(deduplicatedBatch, report)
	}

	// Continue with processing on deduplicated reports
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Set timestamps for all reports in batch
		for _, report := range deduplicatedBatch {
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
		}).Create(&deduplicatedBatch)

		if result.Error != nil {
			return fmt.Errorf("failed to bulk store reports: %w", result.Error)
		}

		return nil
	})
}

// SyncReportsWithRankings synchronizes reports with the provided rankings
func (r *ReportRepository) SyncReportsWithRankings(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) error {
	if len(rankings) == 0 {
		return nil
	}

	log.Printf("[INFO] SyncReportsWithRankings - Processing %d rankings without deleting any reports", len(rankings))

	// Deletion temporarily disabled to test the solution
	return nil
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

// GetAllUniqueReportReferences retrieves all unique report references from rankings
func (r *ReportRepository) GetAllUniqueReportReferences(ctx context.Context) ([]*warcraftlogsBuilds.ClassRanking, error) {
	var rankings []*warcraftlogsBuilds.ClassRanking

	// Get all unique report_code + report_fight_id combinations
	// Even though we're using rankings table, this is a report functionality
	result := r.db.WithContext(ctx).
		Model(&warcraftlogsBuilds.ClassRanking{}).
		Distinct("report_code", "report_fight_id", "encounter_id").
		Find(&rankings)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get unique report references: %w", result.Error)
	}

	log.Printf("[INFO] Retrieved %d unique report references", len(rankings))
	return rankings, nil
}

// MarkReportsForBuildProcessing marks reports as ready for build processing
func (r *ReportRepository) MarkReportsForBuildProcessing(ctx context.Context, identifiers []ReportIdentifier, batchID string) error {
	if len(identifiers) == 0 {
		return nil
	}

	// Use a transaction for this operation
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range identifiers {
			result := tx.Model(&warcraftlogsBuilds.Report{}).
				Where("code = ? AND fight_id = ?", id.Code, id.FightID).
				Updates(map[string]interface{}{
					"build_extraction_status": "pending",
					"processing_batch_id":     batchID,
					"build_extraction_at":     nil,
				})

			if result.Error != nil {
				return fmt.Errorf("failed to mark report %s-#%d: %w", id.Code, id.FightID, result.Error)
			}

			if result.RowsAffected == 0 {
				log.Printf("[WARN] No report found for %s-#%d when marking for build processing", id.Code, id.FightID)
			}
		}
		return nil
	})
}

// GetReportsNeedingBuildExtraction retrieves reports that need build extraction
// This function is used to get reports that need build extraction
// It takes a limit and an offset to paginate through the reports
// It also takes a maxAge to only get reports that are older than the maxAge
func (r *ReportRepository) GetReportsNeedingBuildExtraction(ctx context.Context, limit int, offset int, maxAge time.Duration) ([]*warcraftlogsBuilds.Report, error) {
	var reports []*warcraftlogsBuilds.Report
	minDate := time.Now().Add(-maxAge)

	result := r.db.WithContext(ctx).
		Where("(build_extraction_status = ? OR (build_extraction_status = ? AND build_extraction_at < ?)) AND created_at > ?",
			"pending", "failed", minDate, minDate).
		Order("created_at ASC, code ASC, fight_id ASC").
		Limit(limit).
		Offset(offset).
		Find(&reports)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []*warcraftlogsBuilds.Report{}, nil
		}
		return nil, fmt.Errorf("failed to get reports needing build extraction: %w", result.Error)
	}
	return reports, nil
}

// MarkReportsAsProcessedForBuilds updates the processing status of reports (for builds)
func (r *ReportRepository) MarkReportsAsProcessedForBuilds(ctx context.Context, identifiers []ReportIdentifier, batchID string) error {
	if len(identifiers) == 0 {
		return nil
	}

	// Use a transaction for this operation
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, id := range identifiers {
			result := tx.Model(&warcraftlogsBuilds.Report{}).
				Where("code = ? AND fight_id = ?", id.Code, id.FightID).
				Updates(map[string]interface{}{
					"build_extraction_status": "processed",
					"build_extraction_at":     time.Now(),
					"processing_batch_id":     batchID,
				})

			if result.Error != nil {
				return fmt.Errorf("failed to mark report %s-#%d as processed: %w", id.Code, id.FightID, result.Error)
			}

			if result.RowsAffected == 0 {
				log.Printf("[ERROR] No report found for %s-#%d when marking as processed for builds", id.Code, id.FightID)
				// This situation is probably an error, because the report should exist
			}
		}
		return nil
	})
}

// CountReportsNeedingBuildExtraction counts the total number of reports that need build extraction
func (r *ReportRepository) CountReportsNeedingBuildExtraction(ctx context.Context, maxAge time.Duration) (int64, error) {
	var count int64
	minDate := time.Now().Add(-maxAge)

	result := r.db.WithContext(ctx).
		Model(&warcraftlogsBuilds.Report{}).
		Where("(build_extraction_status = ? OR (build_extraction_status = ? AND build_extraction_at < ?)) AND created_at > ?",
			"pending", "failed", minDate, minDate).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count reports needing build extraction: %w", result.Error)
	}
	return count, nil
}
