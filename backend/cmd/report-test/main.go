package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	configPath := "configs/config_s2_tww.dev.yaml"

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

	// 2. Schedules per class for ReportsWorkflow
	// List of WoW classes
	classes := []string{"Priest", "Warrior", "Mage", "Rogue", "Paladin", "Hunter",
		"Druid", "Shaman", "Warlock", "Monk", "DeathKnight",
		"DemonHunter", "Evoker"}

	// For each class
	for _, className := range classes {
		// Path to the specific class config file
		classConfigPath := fmt.Sprintf("configs/class_config/%s.yaml", strings.ToLower(className))

		// Check if the file exists
		if _, err := os.Stat(classConfigPath); os.IsNotExist(err) {
			logger.Printf("[WARN] Config file for class %s does not exist at %s", className, classConfigPath)
			continue
		}

		// Load the parameters for this class
		reportsParams, err := definitions.LoadReportsParamsForClass(classConfigPath)
		if err != nil {
			logger.Printf("[ERROR] Failed to load reports params for %s: %v", className, err)
			continue
		}

		// Create a specific schedule ID for this class
		scheduleID := fmt.Sprintf("reports-%s", strings.ToLower(className))

		// Create the schedule for this class
		if err := scheduleManager.CreateReportsScheduleForClass(context.Background(), scheduleID, reportsParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create reports schedule for %s: %v", className, err)
			continue
		}

		logger.Printf("[INFO] Successfully created reports schedule for %s with batch ID: %s",
			className, reportsParams.BatchID)
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

	// 4. Schedule for EquipmentAnalysisWorkflow
	equipmentAnalysisParams, err := definitions.LoadEquipmentAnalysisParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load equipment analysis params: %v", err)
	} else {
		if err := scheduleManager.CreateEquipmentAnalysisSchedule(context.Background(), equipmentAnalysisParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create equipment analysis schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created equipment analysis schedule with batch ID: %s", equipmentAnalysisParams.BatchID)
		}
	}

	// 5. Schedule for TalentAnalysisWorkflow
	talentAnalysisParams, err := definitions.LoadTalentAnalysisParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load talent analysis params: %v", err)
	} else {
		if err := scheduleManager.CreateTalentAnalysisSchedule(context.Background(), talentAnalysisParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create talent analysis schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created talent analysis schedule with batch ID: %s", talentAnalysisParams.BatchID)
		}
	}

	// 6. Schedule for StatAnalysisWorkflow
	statAnalysisParams, err := definitions.LoadStatAnalysisParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load stat analysis params: %v", err)
	} else {
		if err := scheduleManager.CreateStatAnalysisSchedule(context.Background(), statAnalysisParams, opts); err != nil {
			logger.Printf("[ERROR] Failed to create stat analysis schedule: %v", err)
		} else {
			logger.Printf("[INFO] Successfully created stat analysis schedule with batch ID: %s", statAnalysisParams.BatchID)
		}
	}

	// Information on how to trigger workflows manually
	logger.Printf("[INFO] All schedules created. Workflows can now be triggered manually from Temporal UI.")
	logger.Printf("[INFO] To trigger workflows manually via code, you can use:")
	logger.Printf("[INFO] - Rankings: scheduleManager.TriggerRankingsNow(ctx)")
	logger.Printf("[INFO] - Reports for specific class: scheduleManager.TriggerReportsForClassNow(ctx, \"ClassName\")")
	logger.Printf("[INFO] - Builds: scheduleManager.TriggerBuildsNow(ctx)")
	logger.Printf("[INFO] - Equipment Analysis: scheduleManager.TriggerEquipmentAnalysisNow(ctx)")
	logger.Printf("[INFO] - Talent Analysis: scheduleManager.TriggerTalentAnalysisNow(ctx)")
	logger.Printf("[INFO] - Stat Analysis: scheduleManager.TriggerStatAnalysisNow(ctx)")

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

	// Context with timeout for cleaning
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Clean all the workflow before shutdown
	logger.Printf("Cleaning up all workflows before shutdown...")
	if err := scheduleManager.CleanupAllWorkflows(ctx); err != nil {
		logger.Printf("Warning: Error during workflows cleanup: %v", err)
	}

	//Clean up all the schedules
	logger.Printf("Cleaning up all schedules before shutdown...")
	if err := scheduleManager.CleanupDecoupledSchedules(ctx); err != nil {
		logger.Printf("Warning: Error during schedules cleanup: %v", err)
	}

	logger.Printf("Cleanup completed, scheduler service shutting down")
}
