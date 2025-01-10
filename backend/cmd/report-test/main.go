package main

import (
	"context"
	"log"
	"time"

	"wowperf/internal/database"
	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsConfig "wowperf/internal/services/warcraftlogs/mythicplus/builds/config"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsService "wowperf/internal/services/warcraftlogs/mythicplus/builds/service"
	warcraftlogsBuildsSync "wowperf/internal/services/warcraftlogs/mythicplus/builds/service/sync"
	warcraftlogsBuildsMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/service/sync/metrics"
)

func main() {
	startTime := time.Now()
	log.Printf("[INFO] Starting full sync process for WoW Performance at %v", startTime.Format(time.RFC3339))

	// Load config
	log.Printf("[INFO] Loading configuration...")
	cfg, err := warcraftlogsBuildsConfig.Load("configs/config_s1_tww.priest.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config: %v", err)
	}
	log.Printf("[INFO] Configuration loaded successfully for %d specs and %d dungeons",
		len(cfg.Specs), len(cfg.Dungeons))

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}

	// Initialize WarcraftLogs client
	warcraftLogsClient, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize WarcraftLogs client: %v", err)
	}

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Rankings.UpdateInterval)
	defer cancel()

	// Initialize metrics
	metrics := warcraftlogsBuildsMetrics.NewSyncMetrics()

	// Initialize worker pool
	workerPool := warcraftlogs.NewWorkerPool(warcraftLogsClient, cfg.Worker.NumWorkers, metrics)
	if err := workerPool.Start(ctx); err != nil {
		log.Fatalf("[FATAL] Failed to start worker pool: %v", err)
	}
	defer workerPool.Stop()

	// Initialize repositories
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)

	// Initialize services
	reportsService := reportsService.NewReportService(warcraftLogsClient, reportsRepo, db, metrics)

	// Initialize sync service
	syncConfig := &warcraftlogsBuildsConfig.Config{
		Rankings: cfg.Rankings,
		Worker:   cfg.Worker,
		Specs:    cfg.Specs,
		Dungeons: cfg.Dungeons,
	}

	syncService := warcraftlogsBuildsSync.NewSyncService(
		workerPool,
		rankingsRepo,
		syncConfig,
	)

	// Phase 1: Rankings Synchronization
	log.Printf("[INFO] Starting rankings sync phase")
	rankingsStart := time.Now()
	if err := syncService.StartSync(ctx); err != nil {
		log.Fatalf("[FATAL] Rankings sync failed: %v", err)
	}
	log.Printf("[INFO] Rankings sync phase completed in %v", time.Since(rankingsStart))

	// Phase 2: Reports Processing
	log.Printf("[INFO] Starting reports fetch phase")
	reportsStart := time.Now()
	reports, err := reportsService.GetReportsFromRankings(ctx)
	if err != nil {
		log.Fatalf("[FATAL] Failed to get reports from rankings: %v", err)
	}
	log.Printf("[INFO] Found %d reports to process", len(reports))

	if len(reports) > 0 {
		log.Printf("[INFO] Starting reports processing phase")
		if err := reportsService.ProcessReports(ctx, reports); err != nil {
			log.Fatalf("[FATAL] Failed to process reports: %v", err)
		}
		log.Printf("[INFO] Reports processing completed in %v", time.Since(reportsStart))
	} else {
		log.Printf("[INFO] No new reports to process")
	}

	// Final Summary
	totalDuration := time.Since(startTime)
	log.Printf("[INFO] Full sync process completed in %v", totalDuration)

	// Log final metrics
	summary := metrics.GetSummary()
	log.Printf("[INFO] Final metrics:")
	log.Printf("- Total Rankings: %d", summary["rankings"].(map[string]int)["total"])
	log.Printf("- Total Reports: %d", summary["reports"].(map[string]int)["total"])
	log.Printf("- Rate Limits Hit: %d", summary["rate_limits"].(map[string]interface{})["total_hits"])
}

// extractDungeonIDs extracts the dungeon IDs from the given dungeons
func extractDungeonIDs(dungeons []warcraftlogsBuildsConfig.Dungeon) []uint {
	ids := make([]uint, len(dungeons))
	for i, dungeon := range dungeons {
		ids[i] = uint(dungeon.ID)
	}
	return ids
}
