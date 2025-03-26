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

// No longer need this constant as we'll use the one from models
// const defaultNamespace = "default"

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

	// Load production configuration
	cfg, err := definitions.LoadConfig("configs/config_s2_tww.dev.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load production config: %v", err)
	}

	// Create production schedule
	opts := scheduler.DefaultScheduleOptions()
	if err := scheduleManager.CreateSchedule(context.Background(), cfg, opts); err != nil {
		logger.Fatalf("[FATAL] Failed to create production schedule: %v", err)
	}
	logger.Printf("[INFO] Successfully created production schedule")

	// Create and trigger test schedule for Priest
	testCfg, err := definitions.LoadConfig("configs/config_s2_tww.priest.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load test config: %v", err)
	}

	if err := scheduleManager.CreateTestSchedule(context.Background(), testCfg, opts); err != nil {
		logger.Printf("[ERROR] Failed to create test schedule: %v", err)
	} else {
		logger.Printf("[TEST] Successfully created test schedule")

		// Trigger test schedule immediately
		if err := scheduleManager.TriggerTestNow(context.Background()); err != nil {
			logger.Printf("[ERROR] Failed to trigger test schedule: %v", err)
		} else {
			logger.Printf("[TEST] Successfully triggered test schedule")
		}
	}

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
		Namespace: models.DefaultNamespace, // Updated to use models
	})
}

func handleGracefulShutdown(scheduleManager *scheduler.ScheduleManager, logger *log.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	logger.Printf("Received signal %v, initiating shutdown", sig)

	// We now have deletion logic implemented but typically don't need to run it on shutdown
	// If you want to clean up on shutdown, you could add:
	// if err := scheduleManager.CleanupAll(context.Background()); err != nil {
	//     logger.Printf("[WARN] Cleanup on shutdown failed: %v", err)
	// }

	logger.Printf("Scheduler service shutdown complete")
}
