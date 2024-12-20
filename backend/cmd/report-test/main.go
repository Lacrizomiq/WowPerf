package main

import (
	"context"
	"log"

	"wowperf/internal/database"
	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsConfig "wowperf/internal/services/warcraftlogs/mythicplus/builds/config"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsService "wowperf/internal/services/warcraftlogs/mythicplus/builds/service"
	warcraftlogsBuildsSync "wowperf/internal/services/warcraftlogs/mythicplus/builds/service/sync"
)

func main() {

	// Load config
	cfg, err := warcraftlogsBuildsConfig.Load("configs/config_s1_tww.dev.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Rankings.UpdateInterval)
	defer cancel()

	// Initialize worker pool
	workerPool := warcraftlogs.NewWorkerPool(warcraftLogsClient, cfg.Worker.NumWorkers) // 3 workers
	workerPool.Start(ctx)                                                               // Start the worker pool
	defer workerPool.Stop()                                                             // Stop the worker pool when the program exits

	// Initialize repository
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)

	// Initialize services
	reportsService := reportsService.NewReportService(warcraftLogsClient, reportsRepo, db)

	// Initialize sync service
	syncConfig := &warcraftlogsBuildsConfig.Config{
		Rankings: cfg.Rankings,
		Worker:   cfg.Worker,
		Specs:    cfg.Specs,
		Dungeons: cfg.Dungeons,
	}

	// Initialize and start sync service
	syncService := warcraftlogsBuildsSync.NewSyncService(
		workerPool,
		rankingsRepo,
		syncConfig,
	)

	log.Println("Starting rankings synchronization")
	if err := syncService.StartSync(ctx); err != nil {
		log.Printf("Rankings sync error: %v", err)
	}

	// Process reports after rankings sync
	log.Println("Processing reports from rankings ...")
	reports, err := reportsService.GetReportsFromRankings(ctx)
	if err != nil {
		log.Fatalf("Failed to get reports from rankings: %v", err)
	}

	log.Printf("Found %d reports", len(reports))
	if len(reports) > 0 {
		if err := reportsService.ProcessReports(ctx, reports); err != nil {
			log.Fatalf("Failed to process reports: %v", err)
		}
	}

	log.Println("Reports processed successfully")
}

// extractDungeonIDs extracts the dungeon IDs from the given dungeons
func extractDungeonIDs(dungeons []warcraftlogsBuildsConfig.Dungeon) []uint {
	ids := make([]uint, len(dungeons))
	for i, dungeon := range dungeons {
		ids[i] = uint(dungeon.ID)
	}
	return ids
}
