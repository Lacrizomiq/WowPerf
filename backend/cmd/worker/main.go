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

	// Ranking, report and player builds repository
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflowStatesRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	// package
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"
	scheduler "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	// workflow
	analyzeWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/analyze"
	buildsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds"
	rankingsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/rankings"
	reportsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/reports"
	syncWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/sync"

	// build analysis repository
	buildsStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	statStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	talentStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

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
	buildsStatsRepo := buildsStatisticsRepository.NewBuildsStatisticsRepository(db)
	talentStatsRepo := talentStatisticsRepository.NewTalentStatisticsRepository(db)
	statStatsRepo := statStatisticsRepository.NewStatStatisticsRepository(db)
	workflowStatesRepo := workflowStatesRepository.NewWorkflowStateRepository(db)
	// Initialize activities with all services
	activitiesService := initializeActivities(
		warcraftLogsClient,
		reportsRepo,
		rankingsRepo,
		playerBuildsRepo,
		buildsStatsRepo,
		talentStatsRepo,
		statStatsRepo,
		workflowStatesRepo,
	)
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

// Initialize activities with all services
func initializeActivities(
	warcraftLogsClient *warcraftlogs.WarcraftLogsClientService,
	reportsRepo *reportsRepository.ReportRepository,
	rankingsRepo *rankingsRepository.RankingsRepository,
	playerBuildsRepo *playerBuildsRepository.PlayerBuildsRepository,
	buildsStatisticsRepo *buildsStatisticsRepository.BuildsStatisticsRepository,
	talentStatisticsRepo *talentStatisticsRepository.TalentStatisticsRepository,
	statStatisticsRepo *statStatisticsRepository.StatStatisticsRepository,
	workflowStatesRepo *workflowStatesRepository.WorkflowStateRepository,
) *activities.Activities {
	// Create the existing activities
	rankingsActivity := activities.NewRankingsActivity(warcraftLogsClient, rankingsRepo)
	reportsActivity := activities.NewReportsActivity(warcraftLogsClient, reportsRepo, rankingsRepo)
	playerBuildsActivity := activities.NewPlayerBuildsActivity(playerBuildsRepo, reportsRepo)
	rateLimitActivity := activities.NewRateLimitActivity(warcraftLogsClient)
	workflowStatesActivity := activities.NewWorkflowStateActivity(workflowStatesRepo)

	// Create the new analysis activities
	buildsStatisticsActivity := activities.NewBuildsStatisticsActivity(
		playerBuildsRepo,
		buildsStatisticsRepo,
	)
	talentStatisticActivity := activities.NewTalentStatisticActivity(
		playerBuildsRepo,
		talentStatisticsRepo,
	)
	statStatisticsActivity := activities.NewStatStatisticsActivity(
		playerBuildsRepo,
		statStatisticsRepo,
	)

	// Return the complete Activities structure
	return &activities.Activities{
		Rankings:         rankingsActivity,
		Reports:          reportsActivity,
		PlayerBuilds:     playerBuildsActivity,
		RateLimit:        rateLimitActivity,
		BuildStatistics:  buildsStatisticsActivity,
		TalentStatistics: talentStatisticActivity,
		StatStatistics:   statStatisticsActivity,
		WorkflowState:    workflowStatesActivity,
	}
}

func registerWorkflowsAndActivities(w worker.Worker, activitiesService *activities.Activities) {
	// Register workflows (New workflows, will be used soon)
	rankingsWorkflowImpl := rankingsWorkflow.NewRankingsWorkflow()
	reportsWorkflowImpl := reportsWorkflow.NewReportsWorkflow()
	buildsWorkflowImpl := buildsWorkflow.NewBuildsWorkflow()

	w.RegisterWorkflowWithOptions(rankingsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: "RankingsWorkflow", // Doit correspondre exactement au nom utilisé dans le scheduler
	})

	w.RegisterWorkflowWithOptions(reportsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: "ReportsWorkflow",
	})

	w.RegisterWorkflowWithOptions(buildsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: "BuildsWorkflow",
	})

	// Register workflows (Old workflows, will be removed soon)
	syncWorkflowImpl := syncWorkflow.NewSyncWorkflow()
	analyzeWorkflowImpl := analyzeWorkflow.NewAnalyzeWorkflow()

	w.RegisterWorkflowWithOptions(syncWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.SyncWorkflowName,
	})

	w.RegisterWorkflowWithOptions(analyzeWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.AnalyzeBuildsWorkflowName,
	})

	// Register rankingsactivities
	w.RegisterActivity(activitiesService.Rankings.FetchAndStore)
	w.RegisterActivity(activitiesService.Rankings.GetStoredRankings)
	w.RegisterActivity(activitiesService.Rankings.MarkRankingsForReportProcessing)

	// Register reports activities
	w.RegisterActivity(activitiesService.Reports.ProcessReports)
	w.RegisterActivity(activitiesService.Reports.GetReportsBatch)
	w.RegisterActivity(activitiesService.Reports.CountAllReports)
	w.RegisterActivity(activitiesService.Reports.GetUniqueReportReferences)
	w.RegisterActivity(activitiesService.Reports.GetRankingsNeedingReportProcessing)
	w.RegisterActivity(activitiesService.Reports.MarkReportsForBuildProcessing)

	// Register player builds activities
	w.RegisterActivity(activitiesService.PlayerBuilds.ProcessAllBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountPlayerBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.GetReportsNeedingBuildExtraction)
	w.RegisterActivity(activitiesService.PlayerBuilds.MarkReportsAsProcessedForBuilds)

	// Register rate limit activities
	w.RegisterActivity(activitiesService.RateLimit.ReservePoints)
	w.RegisterActivity(activitiesService.RateLimit.ReleasePoints)
	w.RegisterActivity(activitiesService.RateLimit.CheckRemainingPoints)

	// Register new analysis activities
	w.RegisterActivity(activitiesService.BuildStatistics.ProcessItemStatistics)
	w.RegisterActivity(activitiesService.TalentStatistics.ProcessTalentStatistics)
	w.RegisterActivity(activitiesService.StatStatistics.ProcessStatStatistics)

	// Workflow state activities
	w.RegisterActivity(activitiesService.WorkflowState.CreateWorkflowState)
	w.RegisterActivity(activitiesService.WorkflowState.UpdateWorkflowState)
	w.RegisterActivity(activitiesService.WorkflowState.GetLastWorkflowRun)
	w.RegisterActivity(activitiesService.WorkflowState.GetWorkflowStatistics)
	w.RegisterActivity(activitiesService.WorkflowState.GetWorkflowStateByID)
	w.RegisterActivity(activitiesService.WorkflowState.DeleteOldWorkflowStates)
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
