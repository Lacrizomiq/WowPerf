package warcraftlogsBuildsService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"wowperf/internal/services/warcraftlogs"
	reportsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	"gorm.io/gorm"
)

type ReportInfo struct {
	ReportCode  string
	FightID     int
	EncounterID uint
}

type ReportService struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *reportsRepository.ReportRepository
	db         *gorm.DB
}

func NewReportService(client *warcraftlogs.WarcraftLogsClientService, repo *reportsRepository.ReportRepository, db *gorm.DB) *ReportService {
	return &ReportService{
		client:     client,
		repository: repo,
		db:         db,
	}
}

// getReportsFromRankings retrieves reports from rankings in database
func (s *ReportService) getReportsFromRankings(ctx context.Context) ([]ReportInfo, error) {
	var reports []struct {
		ReportCode  string `gorm:"column:report_code"`
		FightID     int    `gorm:"column:report_fight_id"`
		EncounterID uint   `gorm:"column:encounter_id"`
	}

	err := s.db.WithContext(ctx).
		Table("class_rankings").
		Select("DISTINCT report_code, report_fight_id as fight_id, encounter_id").
		Scan(&reports).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reports from rankings: %w", err)
	}

	// Convert to ReportInfo slice
	result := make([]ReportInfo, len(reports))
	for i, r := range reports {
		result[i] = ReportInfo{
			ReportCode:  r.ReportCode,
			FightID:     r.FightID,
			EncounterID: r.EncounterID,
		}
	}

	return result, nil
}

// StartPeriodicReportProcessing starts the periodic processing of reports
func (s *ReportService) StartPeriodicReportProcessing(ctx context.Context) {
	log.Println("Starting periodic report processing...")

	go func() {
		for {
			log.Println("Starting reports collection cycle...")

			reports, err := s.getReportsFromRankings(ctx)
			if err != nil {
				log.Printf("Error getting reports from rankings: %v", err)
				continue
			}

			log.Printf("Found %d reports to process", len(reports))

			if len(reports) > 0 {
				if err := s.ProcessReports(ctx, reports); err != nil {
					log.Printf("Error processing reports: %v", err)
				}
			}

			select {
			case <-ctx.Done():
				log.Println("Report processing interrupted")
				return
			case <-time.After(7 * 24 * time.Hour):
				continue
			}
		}
	}()
}

// ProcessReports processes a list of reports
func (s *ReportService) ProcessReports(ctx context.Context, reports []ReportInfo) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(reports))

	for _, report := range reports {
		wg.Add(1)
		go func(r ReportInfo) {
			defer wg.Done()

			if err := s.FetchAndStoreReport(ctx, r.ReportCode, r.FightID, r.EncounterID); err != nil {
				errorChan <- fmt.Errorf("failed to process report %s: %w", r.ReportCode, err)
				return
			}
			log.Printf("processed report %s", r.ReportCode)
		}(report)

		// pause to respect WarcraftLogs API rate limit
		time.Sleep(500 * time.Millisecond)
	}

	// wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// collect errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to process some reports: %v", errors)
	}

	return nil
}

// FetchAndStoreReport fetches and stores a report from WarcraftLogs
func (s *ReportService) FetchAndStoreReport(ctx context.Context, code string, fightID int, encounterID uint) error {
	// First request : Get report details and build talents query
	response, err := s.client.MakeRequest(ctx, reportsQueries.GetReportTableQuery, map[string]interface{}{
		"code":        code,
		"fightID":     fightID,
		"encounterID": encounterID,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch report: %w", err)
	}

	// Parse report and get talents query
	report, talentsQuery, err := reportsQueries.ParseReportDetailsResponse(response, code, fightID, encounterID)
	if err != nil {
		return fmt.Errorf("failed to parse report: %w", err)
	}

	// Second request : Get talents data
	talentsResponse, err := s.client.MakeRequest(ctx, talentsQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch talents: %w", err)
	}

	// Parse talents data
	talentCodes, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
	if err != nil {
		return fmt.Errorf("failed to parse talents: %w", err)
	}

	// Convert talents map to JSON and store in report
	talentCodesJSON, err := json.Marshal(talentCodes)
	if err != nil {
		return fmt.Errorf("failed to marshal talents: %w", err)
	}
	report.TalentCodes = talentCodesJSON

	// Store report in database
	if err := s.repository.StoreReport(ctx, report); err != nil {
		return fmt.Errorf("failed to store report: %w", err)
	}

	return nil
}
