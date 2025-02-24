package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	scheduler "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/client"
)

const (
	defaultNamespace = "default"
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

	// Load production configuration
	cfg, err := workflows.LoadConfig("configs/config_s1_tww.dev.yaml")
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
	testCfg, err := workflows.LoadConfig("configs/config_s1_tww.priest.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load test config: %v", err)
	}

	if err := scheduleManager.CreateTestSchedule(context.Background(), testCfg, opts); err != nil {
		logger.Printf("[ERROR] Failed to create test schedule: %v", err)
	} else {
		logger.Printf("[TEST] Successfully created test schedule")

		// Trigger test schedule immediately
		if err := scheduleManager.TriggerSyncNow(context.Background()); err != nil {
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

	return client.Dial(client.Options{
		HostPort:  temporalAddress,
		Namespace: defaultNamespace,
	})
}

func handleGracefulShutdown(scheduleManager *scheduler.ScheduleManager, logger *log.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	logger.Printf("Received signal %v, initiating shutdown", sig)

	// No need to delete schedules explicitly; closing the client is sufficient
	// If i want to add deletion logic later, scheduler.go would need a Delete method

	logger.Printf("Scheduler service shutdown complete")
}
