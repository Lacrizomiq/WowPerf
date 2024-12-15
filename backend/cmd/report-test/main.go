package main

import (
	"context"
	"log"
	"time"

	"wowperf/internal/database"
	"wowperf/internal/services/warcraftlogs"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsService "wowperf/internal/services/warcraftlogs/mythicplus/builds/service"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize WarcraftLogs client
	warcraftLogsClient, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		log.Fatalf("Failed to initialize WarcraftLogs client: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	// Initialize repository
	repo := reportsRepository.NewReportRepository(db)

	// Test single report first
	log.Println("Testing with a single report first...")
	testReport := reportsService.ReportInfo{
		ReportCode:  "g9Lhy8JmkV1xQ3Gj",
		FightID:     26,
		EncounterID: 12660,
	}
	testReports := []reportsService.ReportInfo{testReport}

	// Create test service
	testService := reportsService.NewReportService(warcraftLogsClient, repo, db)
	if err := testService.ProcessReports(ctx, testReports); err != nil {
		log.Fatalf("Error processing test report: %v", err)
	}
	log.Println("Test report processed successfully")

	// Small delay to ensure clean up
	time.Sleep(2 * time.Second)

	// Create new service for full processing
	log.Println("Proceeding with all reports...")
	mainService := reportsService.NewReportService(warcraftLogsClient, repo, db)

	// Get all reports from rankings
	reports, err := mainService.GetReportsFromRankings(ctx)
	if err != nil {
		log.Fatalf("Failed to get reports from rankings: %v", err)
	}

	log.Printf("Found %d reports to process", len(reports))

	// Process all reports
	if err := mainService.ProcessReports(ctx, reports); err != nil {
		log.Fatalf("Error processing reports: %v", err)
	}

	log.Println("All reports processed successfully")
}
