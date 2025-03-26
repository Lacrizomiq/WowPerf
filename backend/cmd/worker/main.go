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

	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"
	scheduler "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	syncWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/sync"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
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

	// Initialize repositories and activities (no changes needed)
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)
	activitiesService := initializeActivities(warcraftLogsClient, reportsRepo, rankingsRepo, playerBuildsRepo)

	// Create a single worker that will handle both production and test schedules
	taskQueue := scheduler.DefaultScheduleConfig.TaskQueue
	w := worker.New(temporalClient, taskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize:     1,
		MaxConcurrentWorkflowTaskExecutionSize: 2,
		WorkerActivitiesPerSecond:              2,
		MaxConcurrentActivityTaskPollers:       1,
		MaxConcurrentWorkflowTaskPollers:       2,
		EnableSessionWorker:                    true,
	})

	registerWorkflowsAndActivities(w, activitiesService)

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

func initializeActivities(
	warcraftLogsClient *warcraftlogs.WarcraftLogsClientService,
	reportsRepo *reportsRepository.ReportRepository,
	rankingsRepo *rankingsRepository.RankingsRepository,
	playerBuildsRepo *playerBuildsRepository.PlayerBuildsRepository,
) *activities.Activities {
	return activities.NewActivities(
		activities.NewRankingsActivity(warcraftLogsClient, rankingsRepo),
		activities.NewReportsActivity(warcraftLogsClient, reportsRepo),
		activities.NewPlayerBuildsActivity(playerBuildsRepo),
		activities.NewRateLimitActivity(warcraftLogsClient),
	)
}

func registerWorkflowsAndActivities(w worker.Worker, activitiesService *activities.Activities) {
	// Register workflows
	syncWorkflowImpl := syncWorkflow.NewSyncWorkflow()

	w.RegisterWorkflowWithOptions(syncWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.SyncWorkflowName,
	})

	// Register activities
	w.RegisterActivity(activitiesService.Rankings.FetchAndStore)
	w.RegisterActivity(activitiesService.Rankings.GetStoredRankings)
	w.RegisterActivity(activitiesService.Reports.ProcessReports)
	w.RegisterActivity(activitiesService.Reports.GetReportsBatch)
	w.RegisterActivity(activitiesService.Reports.CountAllReports)
	w.RegisterActivity(activitiesService.PlayerBuilds.ProcessAllBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountPlayerBuilds)
	w.RegisterActivity(activitiesService.RateLimit.ReservePoints)
	w.RegisterActivity(activitiesService.RateLimit.ReleasePoints)
	w.RegisterActivity(activitiesService.RateLimit.CheckRemainingPoints)
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
