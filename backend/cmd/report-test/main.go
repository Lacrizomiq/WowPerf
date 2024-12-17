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

	// Create a single service instance that will be used for all operations
	service := reportsService.NewReportService(warcraftLogsClient, repo, db)

	// Process all reports first
	log.Println("Processing all reports...")
	reports, err := service.GetReportsFromRankings(ctx)
	if err != nil {
		log.Fatalf("Failed to get reports from rankings: %v", err)
	}

	log.Printf("Found %d reports to process", len(reports))

	if len(reports) > 0 {
		if err := service.ProcessReports(ctx, reports); err != nil {
			log.Fatalf("Error processing reports: %v", err)
		}
		log.Println("All reports processed successfully")
	}

	// Small delay to ensure clean up
	time.Sleep(2 * time.Second)

	// Force update a specific report for testing
	log.Println("Testing force update of a specific report...")
	testReport := reportsService.ReportInfo{
		ReportCode:  "g9Lhy8JmkV1xQ3Gj",
		FightID:     26,
		EncounterID: 12660,
	}

	// Force update the test report
	if err := service.FetchAndStoreReport(ctx, testReport.ReportCode, testReport.FightID, testReport.EncounterID); err != nil {
		log.Fatalf("Error processing test report: %v", err)
	}
	log.Println("Test report processed successfully")

	// Optional: Start periodic processing
	// Uncomment the following lines if you want to test the periodic processing
	/*
		log.Println("Starting periodic processing...")
		service.StartPeriodicReportProcessing(ctx)

		// Keep the program running
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, shutting down...")
		}
	*/
}
