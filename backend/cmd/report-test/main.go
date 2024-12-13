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

	// Initialize repository and service
	repo := reportsRepository.NewReportRepository(db)
	service := reportsService.NewReportService(warcraftLogsClient, repo, db)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	// Get reports from rankings
	reports, err := service.GetReportsFromRankings(ctx)
	if err != nil {
		log.Fatalf("Failed to get reports from rankings: %v", err)
	}

	log.Printf("Found %d reports to process", len(reports))

	// Process all reports
	if err := service.ProcessReports(ctx, reports); err != nil {
		log.Fatalf("Error processing reports: %v", err)
	}

	log.Println("All reports processed successfully")
}
