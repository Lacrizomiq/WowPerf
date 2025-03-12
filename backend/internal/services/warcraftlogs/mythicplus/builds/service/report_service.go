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
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	warcraftlogsBuildsMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/service/sync/metrics"

	"gorm.io/gorm"
)

type ReportInfo struct {
	ReportCode  string `gorm:"column:report_code"`
	FightID     int    `gorm:"column:report_fight_id"`
	EncounterID uint   `gorm:"column:encounter_id"`
}

type ReportService struct {
	client              *warcraftlogs.WarcraftLogsClientService
	repository          *reportsRepository.ReportRepository
	db                  *gorm.DB
	workerpool          *warcraftlogs.WorkerPool
	PlayerBuildsService *PlayerBuildsService
	metrics             *warcraftlogsBuildsMetrics.SyncMetrics
}

func NewReportService(client *warcraftlogs.WarcraftLogsClientService, repo *reportsRepository.ReportRepository, db *gorm.DB, metrics *warcraftlogsBuildsMetrics.SyncMetrics) *ReportService {
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)
	playerBuildsService := NewPlayerBuildsService(playerBuildsRepo)

	return &ReportService{
		client:              client,
		repository:          repo,
		db:                  db,
		metrics:             metrics,
		workerpool:          warcraftlogs.NewWorkerPool(client, 3, metrics),
		PlayerBuildsService: playerBuildsService,
	}
}

// getReportsFromRankings retrieves reports from rankings in database
func (s *ReportService) GetReportsFromRankings(ctx context.Context) ([]ReportInfo, error) {
	var reports []struct {
		ReportCode  string `gorm:"column:report_code"`
		FightID     int    `gorm:"column:report_fight_id"`
		EncounterID uint   `gorm:"column:encounter_id"`
	}

	// Get all untreated reports, without time condition
	query := s.db.WithContext(ctx).
		Table("class_rankings as cr").
		Select("DISTINCT cr.report_code, cr.report_fight_id, cr.encounter_id").
		Joins("LEFT JOIN warcraft_logs_reports as wlr ON cr.report_code = wlr.code").
		Where("cr.deleted_at IS NULL").
		Where("wlr.code IS NULL") // Only untreated reports

	err := query.Order("cr.report_code DESC").Scan(&reports).Error

	log.Printf("[DEBUG] SQL Query: %v", query.Statement.SQL.String())
	log.Printf("[INFO] Found %d reports needing processing", len(reports))

	if err != nil {
		return nil, fmt.Errorf("failed to get reports from rankings: %w", err)
	}

	log.Printf("[DEBUG] Found %d untreated reports", len(reports))

	result := make([]ReportInfo, len(reports))
	for i, r := range reports {
		result[i] = ReportInfo{
			ReportCode:  r.ReportCode,
			FightID:     r.FightID,
			EncounterID: r.EncounterID,
		}
		log.Printf("[DEBUG] Report to process: %s (FightID: %d)", r.ReportCode, r.FightID)
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

// ProcessReports processes a list of reports with detailed logging and error handling
func (s *ReportService) ProcessReports(ctx context.Context, reports []ReportInfo) error {
	log.Printf("[INFO] Starting processing of %d reports...", len(reports))

	if len(reports) == 0 {
		log.Println("[DEBUG] No reports to process")
		return nil
	}

	totalJobs := len(reports) * 2 // 2 jobs per report (report details and talents)
	log.Printf("[DEBUG] Starting processing of %d reports (expecting %d total jobs)", len(reports), totalJobs)

	// Buffer the channels to avoid blocking
	errorsChan := make(chan error, totalJobs)
	processedReports := make(chan bool, totalJobs)
	done := make(chan struct{}) // Channel to signal completion

	// Start the worker pool
	log.Printf("[DEBUG] Starting worker pool")
	s.workerpool.Start(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	// Goroutine to manage results
	go func() {
		defer func() {
			log.Printf("[DEBUG] Result processor goroutine finishing")
			wg.Done()
		}()

		processedCount := 0
		for processedCount < totalJobs {
			select {
			case result, ok := <-s.workerpool.Results():
				if !ok {
					log.Printf("[DEBUG] Worker pool results channel closed")
					return
				}

				if result.Error != nil {
					log.Printf("[ERROR] Worker result error: %v", result.Error)
					errorsChan <- result.Error
					processedCount++
					continue
				}

				switch result.Job.JobType {
				case "report":
					code := result.Job.Metadata["code"].(string)
					fightID := result.Job.Metadata["fightID"].(int)
					encounterID := result.Job.Metadata["encounterID"].(uint)

					log.Printf("[DEBUG] Processing report %s (FightID: %d, EncounterID: %d)", code, fightID, encounterID)

					report, talentsQuery, err := reportsQueries.ParseReportDetailsResponse(result.Data, code, fightID, encounterID)
					if err != nil {
						log.Printf("[ERROR] Failed to parse report %s: %v", code, err)
						errorsChan <- fmt.Errorf("failed to parse report %s: %w", code, err)
						processedCount++
						continue
					}

					log.Printf("[DEBUG] Successfully parsed report %s, submitting talents query", code)
					processedReports <- true
					processedCount++

					// Submit talents query
					s.workerpool.Submit(warcraftlogs.Job{
						Query:   talentsQuery,
						JobType: "talents",
						Variables: map[string]interface{}{
							"code":    code,
							"fightID": fightID,
						},
						Metadata: map[string]interface{}{
							"report":      report,
							"code":        code,
							"fightID":     fightID,
							"query":       talentsQuery,
							"encounterID": encounterID,
						},
					})
					log.Printf("[DEBUG] Submitted talents query for report %s", code)

				case "talents":
					report := result.Job.Metadata["report"].(*warcraftlogsBuilds.Report)
					code := result.Job.Metadata["code"].(string)
					fightID := result.Job.Metadata["fightID"].(int)

					log.Printf("[DEBUG] Processing talents for report %s (FightID: %d)", code, fightID)

					talentCodes, err := reportsQueries.ParseReportTalentsResponse(result.Data)
					if err != nil {
						log.Printf("[ERROR] Failed to parse talents for report %s: %v", code, err)
						errorsChan <- fmt.Errorf("failed to parse talents for report %s: %w", code, err)
						processedCount++
						continue
					}

					log.Printf("[DEBUG] Successfully parsed talents for report %s", code)

					// Convert and store talents
					talentCodesJSON, err := json.Marshal(talentCodes)
					if err != nil {
						log.Printf("[ERROR] Failed to marshal talents for report %s: %v", code, err)
						errorsChan <- fmt.Errorf("failed to marshal talents for report %s: %w", code, err)
						processedCount++
						continue
					}
					report.TalentCodes = talentCodesJSON

					log.Printf("[DEBUG] Storing report %s in database", code)
					if err := s.repository.StoreReports(ctx, []*warcraftlogsBuilds.Report{report}); err != nil {
						log.Printf("[ERROR] Failed to store report %s: %v", code, err)
						errorsChan <- fmt.Errorf("failed to store report %s: %w", code, err)
						processedCount++
						continue
					}

					// Process player builds
					log.Printf("[DEBUG] Processing player builds for report %s", code)
					if err := s.PlayerBuildsService.ProcessReportBuilds(ctx, []*warcraftlogsBuilds.Report{report}); err != nil {
						log.Printf("[WARN] Failed to process player builds for report %s: %v", code, err)
					}

					log.Printf("[DEBUG] Successfully completed talents processing for report %s", code)
					processedReports <- true
					processedCount++
				}

				if processedCount >= totalJobs {
					log.Printf("[DEBUG] All jobs processed (%d/%d)", processedCount, totalJobs)
					close(done)
					return
				}

			case <-ctx.Done():
				log.Printf("[DEBUG] Context cancelled in result processor")
				return
			}
		}
	}()

	// Submit initial reports to the worker pool
	log.Printf("[DEBUG] Starting submission of %d reports to worker pool", len(reports))
	for i, report := range reports {
		select {
		case <-ctx.Done():
			log.Printf("[DEBUG] Context cancelled during report submission")
			s.workerpool.Stop()
			return ctx.Err()
		case <-done:
			log.Printf("[DEBUG] Processing completed during submission")
			s.workerpool.Stop()
			return nil
		default:
			log.Printf("[DEBUG] Submitting report %d/%d: %s", i+1, len(reports), report.ReportCode)
			s.workerpool.Submit(warcraftlogs.Job{
				Query:   reportsQueries.GetReportTableQuery,
				JobType: "report",
				Variables: map[string]interface{}{
					"code":        report.ReportCode,
					"fightID":     report.FightID,
					"encounterID": report.EncounterID,
				},
				Metadata: map[string]interface{}{
					"code":        report.ReportCode,
					"fightID":     report.FightID,
					"encounterID": report.EncounterID,
				},
			})
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Wait for completion or cancellation
	select {
	case <-done:
		log.Printf("[DEBUG] All processing completed successfully")
		s.workerpool.Stop()
		wg.Wait()
		return nil
	case <-ctx.Done():
		log.Printf("[DEBUG] Context cancelled while waiting for completion")
		s.workerpool.Stop()
		return ctx.Err()
	}
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
	if err := s.repository.StoreReports(ctx, []*warcraftlogsBuilds.Report{report}); err != nil {
		return fmt.Errorf("failed to store report: %w", err)
	}

	if err := s.PlayerBuildsService.ProcessReportBuilds(ctx, []*warcraftlogsBuilds.Report{report}); err != nil {
		return fmt.Errorf("failed to process player builds: %w", err)
	}

	return nil
}
