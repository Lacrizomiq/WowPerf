package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	scheduler "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	"go.temporal.io/sdk/client"
)

func main() {
	logger := log.New(os.Stdout, "[SCHEDULER] ", log.LstdFlags)
	logger.Printf("[INFO] Starting scheduler service")

	// Initialize Temporal client
	temporalClient, err := initTemporalClient()
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// Initialize schedule manager
	scheduleManager := scheduler.NewScheduleManager(temporalClient, logger)

	// Perform cleanup of existing schedules and workflows before creating new ones
	logger.Printf("[INFO] Starting cleanup of existing schedules and workflows")
	if err := scheduleManager.CleanupAll(context.Background()); err != nil {
		logger.Printf("[WARN] Cleanup encountered some errors: %v", err)
	}
	logger.Printf("[INFO] Cleanup completed successfully")

	// Default options for all schedules
	opts := scheduler.DefaultScheduleOptions()

	// Configuration file path
	configPath := "configs/config_s2_tww.priest.yaml"

	// ===== Creation of schedules for new decoupled workflows =====

	// 1. Schedule for RankingsWorkflow
	rankingsParams, err := definitions.LoadRankingsParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load rankings params: %v", err)
	} else {
		if err := scheduleManager.CreateRankingsSchedule(context.Background(), *rankingsParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create rankings schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created rankings schedule with batch ID: %s", rankingsParams.BatchID)
		}
	}

	// 2. Schedule for ReportsWorkflow
	reportsParams, err := definitions.LoadReportsParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load reports params: %v", err)
	} else {
		if err := scheduleManager.CreateReportsSchedule(context.Background(), reportsParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create reports schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created reports schedule with batch ID: %s", reportsParams.BatchID)
		}
	}

	// 3. Schedule for BuildsWorkflow
	buildsParams, err := definitions.LoadBuildsParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load builds params: %v", err)
	} else {
		if err := scheduleManager.CreateBuildsSchedule(context.Background(), buildsParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create builds schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created builds schedule with batch ID: %s", buildsParams.BatchID)
		}
	}

	// ===== Creation of old schedules (for compatibility) =====

	// Loading the configuration for old workflows
	testCfg, err := definitions.LoadConfig(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load test config: %v", err)
	} else {
		// Create the test schedule (for Priest only)
		/* if err := scheduleManager.CreateTestSchedule(context.Background(), testCfg, opts); err != nil {
			logger.Printf("[ERROR] Failed to create test schedule: %v", err)
		} else {
			logger.Printf("[TEST] Successfully created test schedule")
		} */

		// Create the analysis schedule
		if err := scheduleManager.CreateAnalyzeSchedule(context.Background(), testCfg, opts); err != nil {
			logger.Printf("[ERROR] Failed to create analysis schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created analysis schedule")
		}
	}

	// Information on how to trigger workflows manually
	logger.Printf("[INFO] All schedules created. Workflows can now be triggered manually from Temporal UI.")
	logger.Printf("[INFO] To trigger workflows manually via code, you can use:")
	logger.Printf("[INFO] - Rankings: scheduleManager.TriggerRankingsNow(ctx)")
	logger.Printf("[INFO] - Reports: scheduleManager.TriggerReportsNow(ctx)")
	logger.Printf("[INFO] - Builds: scheduleManager.TriggerBuildsNow(ctx)")

	// Waiting for 50 minutes before triggering the analysis workflow (if necessary)
	// Uncomment to activate the automatic triggering of the analysis after 50 minutes
	/*
		logger.Printf("[INFO] Waiting 50 minutes before triggering analysis workflow...")
		time.Sleep(50 * time.Minute)

		if err := scheduleManager.TriggerAnalyzeNow(context.Background()); err != nil {
			logger.Printf("[ERROR] Failed to trigger analysis schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully triggered analysis schedule")
		}
	*/

	// Handle graceful shutdown
	handleGracefulShutdown(scheduleManager, logger)
}

func initTemporalClient() (client.Client, error) {
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = "localhost:7233"
	}

	// Use the shared namespace constant from models
	return client.Dial(client.Options{
		HostPort:  temporalAddress,
		Namespace: models.DefaultNamespace,
	})
}

func handleGracefulShutdown(scheduleManager *scheduler.ScheduleManager, logger *log.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	logger.Printf("Received signal %v, initiating shutdown", sig)

	// We keep the schedules when shutting down to allow scheduled executions
	logger.Printf("Scheduler service shutdown complete")
}
