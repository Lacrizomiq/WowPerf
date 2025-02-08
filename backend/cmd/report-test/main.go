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
	defaultTaskQueue = "warcraft-logs-sync"
	targetClass      = "Priest" // Starting with Priest only
)

func main() {
	log.Printf("[INFO] Starting scheduler service")

	// Initialize Temporal client
	temporalClient, err := initTemporalClient()
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// Create logger
	logger := log.New(os.Stdout, "[SCHEDULER] ", log.LstdFlags)

	// Initialize schedule manager
	scheduleManager := scheduler.NewScheduleManager(temporalClient, logger)

	// Load configuration using the new unified config
	cfg, err := workflows.LoadConfig("configs/config_s1_tww.priest.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config: %v", err)
	}

	// Configure schedule
	opts := scheduler.DefaultScheduleOptions()
	if err := scheduleManager.CreateOrGetClassSchedule(context.Background(), targetClass, cfg, opts); err != nil {
		log.Fatalf("[FATAL] Failed to manage schedule: %v", err)
	}

	logger.Printf("Successfully created schedule for class: %s", targetClass)

	// Test immediate execution
	logger.Printf("Triggering immediate execution for testing...")
	if err := scheduleManager.TriggerSchedule(context.Background(), targetClass); err != nil {
		logger.Printf("[ERROR] Failed to trigger schedule: %v", err)
	} else {
		logger.Printf("Successfully triggered schedule execution")
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

	// Cleanup before shutdown
	ctx := context.Background()
	if err := scheduleManager.DeleteSchedule(ctx, targetClass); err != nil {
		logger.Printf("Error deleting schedule for class %s: %v", targetClass, err)
	}

	logger.Printf("Scheduler service shutdown complete")
}
