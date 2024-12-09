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

	// Test avec un seul report spécifique
	testReport := reportsService.ReportInfo{
		ReportCode:  "g9Lhy8JmkV1xQ3Gj", // Le code que vous avez utilisé précédemment
		FightID:     26,
		EncounterID: 12660,
	}

	log.Printf("Testing with report: %+v", testReport)

	ctx := context.Background()
	if err := service.FetchAndStoreReport(ctx, testReport.ReportCode, testReport.FightID, testReport.EncounterID); err != nil {
		log.Fatalf("Error processing report: %v", err)
	}

	log.Println("Report processing completed successfully")
	time.Sleep(2 * time.Second) // Petit délai pour s'assurer que les logs sont affichés
}
