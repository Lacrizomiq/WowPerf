package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wowperf/internal/database"
	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsConfig "wowperf/internal/services/warcraftlogs/mythicplus/builds/config"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	worker "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/worker"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/client"
)

const (
	temporalNamespace = "default"
	temporalTaskQueue = "warcraft-logs-sync"
)

func main() {
	startTime := time.Now()
	log.Printf("[INFO] Starting WoW Performance sync process at %v", startTime.Format(time.RFC3339))

	// Load config
	// Load configuration
	log.Printf("[INFO] Loading configuration...")
	cfg, err := warcraftlogsBuildsConfig.Load("configs/config_s1_tww.priest.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config: %v", err)
	}

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

	// Initialize repositories
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)

	// Initialize temporal worker
	workerConfig := worker.WorkerConfig{
		TemporalAddress: os.Getenv("TEMPORAL_ADDRESS"),
		Namespace:       temporalNamespace,
		TaskQueue:       temporalTaskQueue,
	}

	temporalWorker, err := worker.NewWorker(
		workerConfig,
		warcraftLogsClient,
		rankingsRepo,
		reportsRepo,
		playerBuildsRepo,
	)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Temporal worker: %v", err)
	}

	// Start the worker in a goroutine
	go func() {
		if err := temporalWorker.Start(); err != nil {
			log.Fatalf("[FATAL] Failed to start worker: %v", err)
		}
	}()

	// Create temporal client for workflow execution
	temporalClient, err := client.NewClient(client.Options{
		HostPort:  os.Getenv("TEMPORAL_ADDRESS"),
		Namespace: temporalNamespace,
	})
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// Convert config to workflow params
	workflowParams := workflows.WorkflowParams{
		Specs:    make([]workflows.ClassSpec, len(cfg.Specs)),
		Dungeons: make([]workflows.Dungeon, len(cfg.Dungeons)),
		BatchConfig: workflows.BatchConfig{
			Size:        cfg.Rankings.Batch.Size,
			MaxPages:    cfg.Rankings.Batch.MaxPages,
			RetryDelay:  cfg.Rankings.Batch.RetryDelay,
			MaxAttempts: cfg.Rankings.Batch.MaxAttempts,
		},
		Rankings: struct {
			MaxRankingsPerSpec int           `json:"max_rankings_per_spec"`
			UpdateInterval     time.Duration `json:"update_interval"`
		}{
			MaxRankingsPerSpec: cfg.Rankings.MaxRankingsPerSpec,
			UpdateInterval:     cfg.Rankings.UpdateInterval,
		},
	}

	// Convertir les Specs
	for i, spec := range cfg.Specs {
		workflowParams.Specs[i] = workflows.ClassSpec{
			ClassName: spec.ClassName,
			SpecName:  spec.SpecName,
		}
	}

	// Convertir les Dungeons
	for i, dungeon := range cfg.Dungeons {
		workflowParams.Dungeons[i] = workflows.Dungeon{
			ID:          dungeon.ID,
			EncounterID: dungeon.EncounterID,
			Name:        dungeon.Name,
			Slug:        dungeon.Slug,
		}
	}

	// Start workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "warcraft-logs-sync-workflow",
		TaskQueue: temporalTaskQueue,
	}

	syncWorkflow := workflows.NewSyncWorkflow()
	we, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		syncWorkflow.Execute,
		workflowParams,
	)
	if err != nil {
		log.Fatalf("[FATAL] Failed to start workflow: %v", err)
	}

	log.Printf("[INFO] Started workflow: ID=%s, RunID=%s", we.GetID(), we.GetRunID())

	// Wait for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interruption or workflow completion
	select {
	case <-sigChan:
		log.Println("[INFO] Received shutdown signal")
	default:
		var result workflows.WorkflowResult
		if err := we.Get(context.Background(), &result); err != nil {
			log.Printf("[ERROR] Workflow failed: %v", err)
		} else {
			log.Printf("[INFO] Workflow completed successfully: Rankings=%d, Reports=%d, Builds=%d",
				result.RankingsProcessed,
				result.ReportsProcessed,
				result.BuildsProcessed)
		}
	}

	// Graceful shutdown
	temporalWorker.Stop()
	log.Printf("[INFO] Sync process completed in %v", time.Since(startTime))
}
