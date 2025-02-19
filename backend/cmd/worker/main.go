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
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"gorm.io/gorm"
)

const (
	defaultNamespace = "default"
)

// WorkerManager handles multiple workers for different task queues
type WorkerManager struct {
	workers map[string]worker.Worker
	logger  *log.Logger
}

func main() {
	logger := log.New(os.Stdout, "[WORKER] ", log.LstdFlags)
	logger.Printf("Starting Temporal workers")

	// Initialize all required services
	temporalClient, db, warcraftLogsClient, err := initializeServices()
	if err != nil {
		logger.Fatalf("Failed to initialize services: %v", err)
	}
	defer temporalClient.Close()

	// Initialize repositories
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)

	// Initialize activities service
	activitiesService := initializeActivities(warcraftLogsClient, reportsRepo, rankingsRepo, playerBuildsRepo)

	// Create worker manager
	workerMgr := &WorkerManager{
		workers: make(map[string]worker.Worker),
		logger:  logger,
	}

	// Create a worker for each time slot
	for _, slot := range scheduler.ScheduleSlots {
		w := worker.New(temporalClient, slot.TaskQueue, worker.Options{
			MaxConcurrentActivityExecutionSize:     1,
			MaxConcurrentWorkflowTaskExecutionSize: 2,
			WorkerActivitiesPerSecond:              2,
			MaxConcurrentActivityTaskPollers:       1,
			MaxConcurrentWorkflowTaskPollers:       2,
			EnableSessionWorker:                    true,
		})

		registerWorkflowsAndActivities(w, activitiesService)

		workerMgr.workers[slot.TaskQueue] = w
		logger.Printf("Created worker for task queue: %s (Classes: %v)", slot.TaskQueue, slot.Classes)
	}

	// Start all workers
	var wg sync.WaitGroup
	workerErrorChan := make(chan error, len(workerMgr.workers))
	interruptChan := make(chan interface{})

	for taskQueue, w := range workerMgr.workers {
		wg.Add(1)
		go func(tq string, worker worker.Worker) {
			defer wg.Done()
			logger.Printf("Starting worker for task queue: %s", tq)
			if err := worker.Run(interruptChan); err != nil {
				logger.Printf("Worker error in task queue %s: %v", tq, err)
				workerErrorChan <- err
			}
		}(taskQueue, w)
	}

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
		Namespace: defaultNamespace,
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
	w.RegisterWorkflow(workflows.SyncWorkflow)
	w.RegisterWorkflow(workflows.ProcessBuildBatch)

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
		logger.Printf("Shutdown signal received, stopping workers...")
	case err := <-errorChan:
		logger.Printf("Worker error detected, initiating shutdown: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	for taskQueue, w := range mgr.workers {
		logger.Printf("Stopping worker for task queue: %s", taskQueue)
		w.Stop()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Printf("All workers stopped gracefully")
	case <-ctx.Done():
		logger.Printf("Shutdown timeout exceeded, some workers may not have stopped cleanly")
	}
}
