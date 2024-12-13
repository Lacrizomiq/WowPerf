package warcraftlogsBuildsService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	reportsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	"gorm.io/gorm"
)

type ReportInfo struct {
	ReportCode  string `gorm:"column:report_code"`
	FightID     int    `gorm:"column:report_fight_id"`
	EncounterID uint   `gorm:"column:encounter_id"`
}

type ReportService struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *reportsRepository.ReportRepository
	db         *gorm.DB
	workerpool *warcraftlogs.WorkerPool
}

func NewReportService(client *warcraftlogs.WarcraftLogsClientService, repo *reportsRepository.ReportRepository, db *gorm.DB) *ReportService {
	return &ReportService{
		client:     client,
		repository: repo,
		db:         db,
		workerpool: warcraftlogs.NewWorkerPool(client, 3), // 3 concurrent requests by default
	}
}

// getReportsFromRankings retrieves reports from rankings in database
func (s *ReportService) GetReportsFromRankings(ctx context.Context) ([]ReportInfo, error) {
	var reports []struct {
		ReportCode  string `gorm:"column:report_code"`
		FightID     int    `gorm:"column:report_fight_id"`
		EncounterID uint   `gorm:"column:encounter_id"`
	}

	err := s.db.WithContext(ctx).
		Table("class_rankings").
		Select("DISTINCT report_code, report_fight_id, encounter_id").
		Where("deleted_at IS NULL AND encounter_id = ?", 12660). // 12660 = Ara-Kara
		Order("report_code DESC").
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

			reports, err := s.GetReportsFromRankings(ctx)
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
	log.Printf("Processing %d reports with worker pool", len(reports))

	s.workerpool.Start(ctx)
	defer s.workerpool.Stop()

	// Channels for the results management
	errorsChan := make(chan error, len(reports))
	processedReports := make(chan bool, len(reports))

	// Create a WaitGroup to ensure we wait for the result processor goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine to manage results
	go func() {
		defer wg.Done()
		for result := range s.workerpool.Results() {
			if result.Error != nil {
				log.Printf("Error processing report: %v", result.Error)
				errorsChan <- result.Error
				continue
			}

			switch result.Job.JobType {
			case "report":
				// Extract the job metadata
				code := result.Job.Metadata["code"].(string)
				fightID := result.Job.Metadata["fightID"].(int)
				encounterID := result.Job.Metadata["encounterID"].(uint)

				log.Printf("Processing report result: %s (FightID: %d, EncounterID: %d)", code, fightID, encounterID)

				// Parse the report
				report, talentsQuery, err := reportsQueries.ParseReportDetailsResponse(result.Data, code, fightID, encounterID)
				if err != nil {
					log.Printf("Failed to parse report %s: %v", code, err)
					errorsChan <- fmt.Errorf("failed to parse report %s: %w", code, err)
					continue
				}

				// Submit the talents query
				log.Printf("Submitting talents query for report: %s", code)
				s.workerpool.Submit(warcraftlogs.Job{
					Query:    talentsQuery,
					JobType:  "talents",
					Metadata: map[string]interface{}{"report": report},
				})

			case "talents":
				report := result.Job.Metadata["report"].(*warcraftlogsBuilds.Report)
				log.Printf("Processing talents for report: %s", report.Code)

				talentCodes, err := reportsQueries.ParseReportTalentsResponse(result.Data)
				if err != nil {
					log.Printf("Failed to parse talents for report %s: %v", report.Code, err)
					errorsChan <- fmt.Errorf("failed to parse talents for report %s: %w", report.Code, err)
					continue
				}

				talentCodesJSON, err := json.Marshal(talentCodes)
				if err != nil {
					log.Printf("Failed to marshal talents for report %s: %v", report.Code, err)
					errorsChan <- fmt.Errorf("failed to marshal talents for report %s: %w", report.Code, err)
					continue
				}
				report.TalentCodes = talentCodesJSON

				if err := s.repository.StoreReport(ctx, report); err != nil {
					log.Printf("Failed to store report %s: %v", report.Code, err)
					errorsChan <- fmt.Errorf("failed to store report %s: %w", report.Code, err)
					continue
				}

				processedReports <- true
			}
		}
	}()

	// Submit reports to the worker pool
	for _, report := range reports {
		log.Printf("Submitting report: %s (FightID: %d, EncounterID: %d)", report.ReportCode, report.FightID, report.EncounterID)
		s.workerpool.Submit(warcraftlogs.Job{
			Query:   reportsQueries.GetReportTableQuery,
			JobType: "report",
			Metadata: map[string]interface{}{
				"code":        report.ReportCode,
				"fightID":     report.FightID,
				"encounterID": report.EncounterID,
			},
		})
	}

	// Wait for all reports to be processed
	var errors []error
	completedCount := 0
	totalReports := len(reports)

	for completedCount < totalReports {
		select {
		case err := <-errorsChan:
			errors = append(errors, err)
		case <-processedReports:
			completedCount++
			log.Printf("Progress: %d/%d reports processed", completedCount, totalReports)
		case <-ctx.Done():
			log.Println("Report processing interrupted")
			return ctx.Err()
		}
	}

	// Wait for the results processing goroutine to finish
	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("failed to process some reports: %v", errors)
	}

	log.Println("All reports processed successfully")
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
