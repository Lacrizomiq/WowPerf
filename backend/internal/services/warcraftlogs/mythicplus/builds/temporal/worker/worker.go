package warcraftlogsBuildsTemporalWorker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

const defaultTaskQueue = "warcraft-logs-sync"

type Worker struct {
	client client.Client
	worker worker.Worker
}

func NewWorker(
	warcraftlogsClient *warcraftlogs.WarcraftLogsClientService,
	rankingsRepository *warcraftlogsBuildsRepository.RankingsRepository,
	reportsRepository *warcraftlogsBuildsRepository.ReportRepository,
	playerBuildsRepository *warcraftlogsBuildsRepository.PlayerBuildsRepository,
) (*Worker, error) {
	// Get Temporal configuration from environment
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = "localhost:7233"
	}

	temporalNamespace := os.Getenv("TEMPORAL_NAMESPACE")
	if temporalNamespace == "" {
		temporalNamespace = "default"
	}

	taskQueue := os.Getenv("TEMPORAL_TASKQUEUE")
	if taskQueue == "" {
		taskQueue = defaultTaskQueue
	}

	// Create the Temporal client
	temporalClient, err := client.Dial(client.Options{
		HostPort:  temporalAddress,
		Namespace: temporalNamespace,
	})
	if err != nil {
		return nil, err
	}

	// Create worker
	w := worker.New(temporalClient, taskQueue, worker.Options{})

	// Initialize activities
	rankingsActivity := activities.NewRankingsActivity(warcraftlogsClient, rankingsRepository)
	reportsActivity := activities.NewReportsActivity(warcraftlogsClient, reportsRepository)
	playerBuildsActivity := activities.NewPlayerBuildsActivity(playerBuildsRepository)

	activitiesService := activities.NewActivities(
		rankingsActivity,
		reportsActivity,
		playerBuildsActivity,
	)

	// Register workflow
	w.RegisterWorkflow(workflows.SyncWorkflow)

	// Register activities
	w.RegisterActivity(activitiesService.Rankings.FetchAndStore)
	w.RegisterActivity(activitiesService.Reports.ProcessReports)
	w.RegisterActivity(activitiesService.Reports.GetProcessedReports)
	w.RegisterActivity(activitiesService.Reports.GetReportsForEncounter)
	w.RegisterActivity(activitiesService.PlayerBuilds.ProcessBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountPlayerBuilds)

	return &Worker{
		client: temporalClient,
		worker: w,
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	// Start worker
	log.Printf("[INFO] Starting Temporal worker...")
	if err := w.worker.Start(); err != nil {
		return err
	}

	// Setup graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal or context cancellation
	select {
	case <-signalChan:
		log.Println("[INFO] Shutdown signal received")
	case <-ctx.Done():
		log.Println("[INFO] Context cancelled")
	}

	w.Stop()
	return nil
}

func (w *Worker) Stop() {
	log.Println("[INFO] Stopping worker...")
	w.worker.Stop()
	w.client.Close()
	log.Println("[INFO] Worker stopped")
}
