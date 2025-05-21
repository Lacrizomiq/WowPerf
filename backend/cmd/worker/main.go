package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wowperf/internal/database"
	"wowperf/internal/services/warcraftlogs"

	// Scheduler pour la configuration de la queue
	scheduler "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler"
	buildsDefinitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	// Import des packages d'initialisation pour chaque feature
	buildsInit "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal"
	playerRankingsInit "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"gorm.io/gorm"
)

// WorkerManager handles the worker for the single task queue
type WorkerManager struct {
	worker worker.Worker
	logger *log.Logger
}

func main() {
	logger := log.New(os.Stdout, "[WORKER] ", log.LstdFlags)
	logger.Printf("Starting Temporal worker")

	// Initialize services
	temporalClient, db, warcraftLogsClient, err := initializeServices()
	if err != nil {
		logger.Fatalf("Failed to initialize services: %v", err)
	}
	defer temporalClient.Close()

	// Charger la configuration
	config, err := buildsDefinitions.LoadConfig("configs/config_s2_tww.dev.yaml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialiser les features
	logger.Printf("Initializing features")

	// Initialiser la feature builds
	_, _, _, _, _, _, _, buildsActivitiesService := buildsInit.InitBuilds(
		db,
		warcraftLogsClient,
		int(config.Rankings.MaxRankingsPerSpec),
	)

	// Initialiser la feature player rankings
	_, playerRankingsActivitiesService := playerRankingsInit.InitPlayerRankings(
		db,
		warcraftLogsClient,
	)

	// Create a single worker
	taskQueue := scheduler.DefaultScheduleConfig.TaskQueue
	w := worker.New(temporalClient, taskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize:     1,
		MaxConcurrentWorkflowTaskExecutionSize: 2,
		WorkerActivitiesPerSecond:              2,
		MaxConcurrentActivityTaskPollers:       1,
		MaxConcurrentWorkflowTaskPollers:       2,
		EnableSessionWorker:                    true,
	})

	// Enregistrer les workflows et activities pour chaque feature
	logger.Printf("Registering workflows and activities")
	buildsInit.RegisterBuilds(w, buildsActivitiesService)
	playerRankingsInit.RegisterPlayerRankings(w, playerRankingsActivitiesService)

	workerMgr := &WorkerManager{
		worker: w,
		logger: logger,
	}
	logger.Printf("Created worker for task queue: %s", taskQueue)

	// Start worker
	var wg sync.WaitGroup
	workerErrorChan := make(chan error, 1)
	interruptChan := make(chan interface{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Printf("Starting worker for task queue: %s", taskQueue)
		if err := workerMgr.worker.Run(interruptChan); err != nil {
			logger.Printf("Worker error in task queue %s: %v", taskQueue, err)
			workerErrorChan <- err
		}
	}()

	// Handle graceful shutdown
	handleGracefulShutdown(workerMgr, &wg, workerErrorChan, logger)
}

func initializeServices() (client.Client, *gorm.DB, *warcraftlogs.WarcraftLogsClientService, error) {
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = "localhost:7233"
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort:  temporalAddress,
		Namespace: models.DefaultNamespace,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create Temporal client: %w", err)
	}

	db, err := database.InitDB()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	warcraftLogsClient, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize WarcraftLogs client: %w", err)
	}

	return temporalClient, db, warcraftLogsClient, nil
}

func handleGracefulShutdown(mgr *WorkerManager, wg *sync.WaitGroup, errorChan chan error, logger *log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		logger.Printf("Shutdown signal received, stopping worker...")
	case err := <-errorChan:
		logger.Printf("Worker error detected, initiating shutdown: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger.Printf("Stopping worker for task queue: %s", scheduler.DefaultScheduleConfig.TaskQueue)
	mgr.worker.Stop()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Printf("Worker stopped gracefully")
	case <-ctx.Done():
		logger.Printf("Shutdown timeout exceeded, worker may not have stopped cleanly")
	}
}
