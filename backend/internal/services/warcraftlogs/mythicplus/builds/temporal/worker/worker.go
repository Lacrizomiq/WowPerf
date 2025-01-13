package warcraftlogsBuildsTemporalWorker

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

type WorkerConfig struct {
	TemporalAddress string
	Namespace       string
	TaskQueue       string
}

type Worker struct {
	config     WorkerConfig
	client     client.Client
	worker     worker.Worker
	activities *activities.Activities
	workflows  *workflows.SyncWorkflow
}

func NewWorker(
	config WorkerConfig,
	warcraftlogsClient *warcraftlogs.WarcraftLogsClientService,
	rankingsRepository *warcraftlogsBuildsRepository.RankingsRepository,
	reportsRepository *warcraftlogsBuildsRepository.ReportRepository,
	playerBuildsRepository *warcraftlogsBuildsRepository.PlayerBuildsRepository,
) (*Worker, error) {
	// Create Temporal client
	c, err := client.NewClient(client.Options{
		HostPort:  config.TemporalAddress,
		Namespace: config.Namespace,
	})
	if err != nil {
		return nil, err
	}

	// Initialize activities
	activities := &activities.Activities{
		Rankings:     activities.NewRankingsActivity(warcraftlogsClient, rankingsRepository),
		Reports:      activities.NewReportsActivity(warcraftlogsClient, reportsRepository),
		PlayerBuilds: activities.NewPlayerBuildsActivity(playerBuildsRepository),
	}

	// Initialize workflows
	workflow := workflows.NewSyncWorkflow()

	return &Worker{
		config:     config,
		client:     c,
		activities: activities,
		workflows:  workflow,
	}, nil
}

func (w *Worker) Start() error {
	log.Printf("Starting Temporal worker with namespace: %s, task queue: %s",
		w.config.Namespace, w.config.TaskQueue)

	// Create worker instance
	wrk := worker.New(w.client, w.config.TaskQueue, worker.Options{})

	// Register workflow
	wrk.RegisterWorkflow(w.workflows.Execute)

	// Register activities
	wrk.RegisterActivity(w.activities.Rankings.FetchAndStore)
	wrk.RegisterActivity(w.activities.Reports.ProcessReports)
	wrk.RegisterActivity(w.activities.Reports.GetProcessedReports)
	wrk.RegisterActivity(w.activities.PlayerBuilds.ProcessBuilds)

	// Start listening for the Task Queue
	err := wrk.Run(worker.InterruptCh())
	if err != nil {
		return err
	}

	log.Printf("Worker started successfully")
	return nil
}

func (w *Worker) Stop() {
	if w.client != nil {
		w.client.Close()
	}
	log.Printf("Worker stopped")
}
