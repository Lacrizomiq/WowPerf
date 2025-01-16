package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// Initialize Temporal client
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = "localhost:7233"
	}

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

	// Initialize services and repositories
	warcraftLogsClient, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize WarcraftLogs client: %v", err)
	}

	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db)
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)

	// Create worker with the temporalClient
	w := worker.New(temporalClient, defaultTaskQueue, worker.Options{})

	// Initialize activities
	rankingsActivity := activities.NewRankingsActivity(warcraftLogsClient, rankingsRepo)
	reportsActivity := activities.NewReportsActivity(warcraftLogsClient, reportsRepo)
	playerBuildsActivity := activities.NewPlayerBuildsActivity(playerBuildsRepo)

	activitiesService := activities.NewActivities(
		rankingsActivity,
		reportsActivity,
		playerBuildsActivity,
	)

	// Register workflow and activities
	w.RegisterWorkflow(workflows.SyncWorkflow)
	w.RegisterActivity(activitiesService.Rankings.FetchAndStore)
	w.RegisterActivity(activitiesService.Reports.ProcessReports)
	w.RegisterActivity(activitiesService.Reports.GetProcessedReports)
	w.RegisterActivity(activitiesService.PlayerBuilds.ProcessBuilds)
	w.RegisterActivity(activitiesService.Reports.GetReportsForEncounter)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountPlayerBuilds)

	// Setup signal handling
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	log.Printf("[INFO] Starting worker, waiting for workflows...")

	// Start worker in a goroutine
	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Printf("[ERROR] Worker stopped with error: %v", err)
			interruptChan <- os.Interrupt // Signal main routine to exit
		}
	}()

	// Wait for interrupt signal
	<-interruptChan
	log.Printf("[INFO] Shutdown signal received, stopping worker...")

	// Stop the worker
	w.Stop()
	log.Printf("[INFO] Worker stopped gracefully")
}
