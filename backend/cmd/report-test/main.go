package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Builds scheduler
	buildsScheduler "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"

	// Player Rankings scheduler
	playerRankingsScheduler "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/scheduler"

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

	// Initialize schedule managers
	buildsScheduleManager := buildsScheduler.NewScheduleManager(temporalClient, logger)
	playerRankingsScheduleManager := playerRankingsScheduler.NewPlayerRankingsScheduleManager(temporalClient, logger)

	// Perform cleanup of existing schedules and workflows before creating new ones
	logger.Printf("[INFO] Starting cleanup of existing schedules and workflows")
	// Cleanup builds schedules
	if err := buildsScheduleManager.CleanupAll(context.Background()); err != nil {
		logger.Printf("[WARN] Builds cleanup encountered some errors: %v", err)
	}
	// Cleanup player rankings schedules
	if err := playerRankingsScheduleManager.CleanupAll(context.Background()); err != nil {
		logger.Printf("[WARN] Player rankings cleanup encountered some errors: %v", err)
	}
	logger.Printf("[INFO] Cleanup completed successfully")

	// Default options for all schedules
	buildsOpts := buildsScheduler.DefaultScheduleOptions()
	playerRankingsOpts := playerRankingsScheduler.DefaultScheduleOptions()

	// Configuration file path
	configPath := "configs/config_s2_tww.dev.yaml"

	// Initialize builds schedules
	if err := buildsScheduler.InitBuildsSchedules(context.Background(), buildsScheduleManager, configPath, buildsOpts, logger); err != nil {
		logger.Printf("[ERROR] Failed to initialize builds schedules: %v", err)
	}

	// Initialize player rankings schedule
	if err := playerRankingsScheduler.InitPlayerRankingsSchedule(context.Background(), playerRankingsScheduleManager, configPath, playerRankingsOpts, logger); err != nil {
		logger.Printf("[ERROR] Failed to initialize player rankings schedule: %v", err)
	}

	// Log manual trigger instructions
	logger.Printf("[INFO] All schedules created. Workflows can now be triggered manually from Temporal UI.")
	buildsScheduler.LogBuildsManualTriggerInstructions(logger)
	playerRankingsScheduler.LogPlayerRankingsManualTriggerInstructions(logger)

	// Handle graceful shutdown
	handleGracefulShutdown(buildsScheduleManager, playerRankingsScheduleManager, logger)
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

func handleGracefulShutdown(buildsScheduleManager *buildsScheduler.ScheduleManager, playerRankingsScheduleManager *playerRankingsScheduler.PlayerRankingsScheduleManager, logger *log.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	logger.Printf("Received signal %v, initiating shutdown", sig)

	// Context with timeout for cleaning
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Clean all the builds workflows before shutdown
	logger.Printf("Cleaning up all builds workflows before shutdown...")
	if err := buildsScheduleManager.CleanupAllWorkflows(ctx); err != nil {
		logger.Printf("Warning: Error during builds workflows cleanup: %v", err)
	}

	// Clean all the player rankings workflows before shutdown
	logger.Printf("Cleaning up all player rankings workflows before shutdown...")
	if err := playerRankingsScheduleManager.CleanupPlayerRankingsWorkflows(ctx); err != nil {
		logger.Printf("Warning: Error during player rankings workflows cleanup: %v", err)
	}

	// Clean up all the builds schedules
	logger.Printf("Cleaning up all builds schedules before shutdown...")
	if err := buildsScheduleManager.CleanupDecoupledSchedules(ctx); err != nil {
		logger.Printf("Warning: Error during builds schedules cleanup: %v", err)
	}

	// Clean up all the player rankings schedules
	logger.Printf("Cleaning up all player rankings schedules before shutdown...")
	if err := playerRankingsScheduleManager.DeletePlayerRankingsSchedule(ctx); err != nil {
		logger.Printf("Warning: Error during player rankings schedules cleanup: %v", err)
	}

	logger.Printf("Cleanup completed, scheduler service shutting down")
}
