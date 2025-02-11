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
	cfg, err := workflows.LoadConfig("configs/config_s1_tww.dev.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config: %v", err)
	}

	// Configure schedule
	opts := scheduler.DefaultScheduleOptions()

	// Create a map to track unique classes
	classMap := make(map[string]bool)
	for _, spec := range cfg.Specs {
		if !classMap[spec.ClassName] {
			classMap[spec.ClassName] = true
			// Create a schedule for each class
			if err := scheduleManager.CreateOrGetClassSchedule(context.Background(), spec.ClassName, cfg, opts); err != nil {
				log.Printf("[ERROR] Failed to manage schedule for class %s: %v", spec.ClassName, err)
				continue
			}
			logger.Printf("Successfully created schedule for class: %s", spec.ClassName)

			// Trigger immediate execution
			if err := scheduleManager.TriggerSchedule(context.Background(), spec.ClassName); err != nil {
				logger.Printf("[ERROR] Failed to trigger schedule for class %s: %v", spec.ClassName, err)
			}
		}
	}

	// Handle graceful shutdown
	handleGracefulShutdown(scheduleManager, logger, classMap)
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

func handleGracefulShutdown(scheduleManager *scheduler.ScheduleManager, logger *log.Logger, classMap map[string]bool) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	logger.Printf("Received signal %v, initiating shutdown", sig)

	// Cleanup before shutdown
	ctx := context.Background()
	for className := range classMap {
		if err := scheduleManager.DeleteSchedule(ctx, className); err != nil {
			logger.Printf("Error deleting schedule for class %s: %v", className, err)
		}
	}

	logger.Printf("Scheduler service shutdown complete")
}
