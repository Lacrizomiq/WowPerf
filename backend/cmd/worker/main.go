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
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	defaultNamespace = "default"
	defaultTaskQueue = "warcraft-logs-sync"
)

func main() {
	log.Printf("[INFO] Starting Temporal worker")

	// Get Temporal address from environment or use default
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = "localhost:7233"
	}

	// Initialize Temporal client
	temporalClient, err := client.Dial(client.Options{
		HostPort:  temporalAddress,
		Namespace: defaultNamespace,
	})
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}

	// Initialize WarcraftLogs client service
	warcraftLogsClient, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize WarcraftLogs client: %v", err)
	}

	// Initialize repositories
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)

	// Configure worker with optimized options
	w := worker.New(temporalClient, defaultTaskQueue, worker.Options{
		// Limit concurrent executions to match batch sizes
		MaxConcurrentActivityExecutionSize:     5,
		MaxConcurrentWorkflowTaskExecutionSize: 2,

		// Rate limiting for better resource management
		WorkerActivitiesPerSecond: 5,

		// Optimize task polling
		MaxConcurrentActivityTaskPollers: 2,
		MaxConcurrentWorkflowTaskPollers: 2,

		// Enable sessions for better resource management
		EnableSessionWorker: true,
	})

	// Initialize activities
	rankingsActivity := activities.NewRankingsActivity(warcraftLogsClient, rankingsRepo)
	reportsActivity := activities.NewReportsActivity(warcraftLogsClient, reportsRepo)
	playerBuildsActivity := activities.NewPlayerBuildsActivity(playerBuildsRepo)
	rateLimitActivity := activities.NewRateLimitActivity(warcraftLogsClient)
	activitiesService := activities.NewActivities(
		rankingsActivity,
		reportsActivity,
		playerBuildsActivity,
		rateLimitActivity,
	)

	// Register workflow and activities with proper activity names
	w.RegisterWorkflow(workflows.SyncWorkflow)
	w.RegisterWorkflow(workflows.ProcessBuildBatch)

	// Rankings activities
	w.RegisterActivity(activitiesService.Rankings.FetchAndStore)

	// Reports activities
	w.RegisterActivity(activitiesService.Reports.ProcessReports)
	w.RegisterActivity(activitiesService.Reports.GetProcessedReports)
	w.RegisterActivity(activitiesService.Reports.GetReportsForEncounter)
	w.RegisterActivity(activitiesService.Reports.CountReportsForEncounter)
	w.RegisterActivity(activitiesService.Reports.GetReportsForEncounterBatch)

	// Player builds activities
	w.RegisterActivity(activitiesService.PlayerBuilds.ProcessBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountPlayerBuilds)

	// Rate limit activity
	w.RegisterActivity(activitiesService.RateLimit.ReservePoints)
	w.RegisterActivity(activitiesService.RateLimit.ReleasePoints)
	w.RegisterActivity(activitiesService.RateLimit.CheckRemainingPoints)

	// Setup signal handling for graceful shutdown
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	log.Printf("[INFO] Starting worker, waiting for workflows...")

	// Start worker in a goroutine with enhanced error handling
	workerErrorChan := make(chan error, 1)
	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Printf("[ERROR] Worker stopped with error: %v", err)
			workerErrorChan <- err
			interruptChan <- os.Interrupt // Signal main routine to exit
		}
	}()

	// Implement graceful shutdown
	select {
	case <-interruptChan:
		log.Printf("[INFO] Shutdown signal received, initiating graceful shutdown...")

		// Create context with timeout for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer shutdownCancel()

		// Stop the worker with timeout
		done := make(chan struct{})
		go func() {
			w.Stop()
			close(done)
		}()

		select {
		case <-done:
			log.Printf("[INFO] Worker stopped gracefully")
		case <-shutdownCtx.Done():
			log.Printf("[WARN] Worker stop timeout exceeded, forcing shutdown")
		}

	case err := <-workerErrorChan:
		log.Printf("[ERROR] Worker encountered an error, shutting down: %v", err)
	}
}
