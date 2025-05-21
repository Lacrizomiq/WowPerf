// internal/services/warcraftlogs/mythicplus/builds/temporal/scheduler/init.go
package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
)

// InitBuildsSchedules initialise tous les schedules pour la feature builds
func InitBuildsSchedules(ctx context.Context, scheduleManager *ScheduleManager, configPath string, opts *ScheduleOptions, logger *log.Logger) error {
	// 1. Schedule for RankingsWorkflow
	rankingsParams, err := definitions.LoadRankingsParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load rankings params: %v", err)
		return err
	}

	if err := scheduleManager.CreateRankingsSchedule(ctx, *rankingsParams, opts); err != nil {
		logger.Printf("[ERROR] Failed to create rankings schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created rankings schedule with batch ID: %s", rankingsParams.BatchID)

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
		if err := scheduleManager.CreateReportsScheduleForClass(ctx, scheduleID, reportsParams, opts); err != nil {
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
		return err
	}

	if err := scheduleManager.CreateBuildsSchedule(ctx, buildsParams, opts); err != nil {
		logger.Printf("[ERROR] Failed to create builds schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created builds schedule with batch ID: %s", buildsParams.BatchID)

	// 4. Schedule for EquipmentAnalysisWorkflow
	equipmentAnalysisParams, err := definitions.LoadEquipmentAnalysisParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load equipment analysis params: %v", err)
		return err
	}

	if err := scheduleManager.CreateEquipmentAnalysisSchedule(ctx, equipmentAnalysisParams, opts); err != nil {
		logger.Printf("[ERROR] Failed to create equipment analysis schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created equipment analysis schedule with batch ID: %s", equipmentAnalysisParams.BatchID)

	// 5. Schedule for TalentAnalysisWorkflow
	talentAnalysisParams, err := definitions.LoadTalentAnalysisParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load talent analysis params: %v", err)
		return err
	}

	if err := scheduleManager.CreateTalentAnalysisSchedule(ctx, talentAnalysisParams, opts); err != nil {
		logger.Printf("[ERROR] Failed to create talent analysis schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created talent analysis schedule with batch ID: %s", talentAnalysisParams.BatchID)

	// 6. Schedule for StatAnalysisWorkflow
	statAnalysisParams, err := definitions.LoadStatAnalysisParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load stat analysis params: %v", err)
		return err
	}

	if err := scheduleManager.CreateStatAnalysisSchedule(ctx, statAnalysisParams, opts); err != nil {
		logger.Printf("[ERROR] Failed to create stat analysis schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created stat analysis schedule with batch ID: %s", statAnalysisParams.BatchID)

	return nil
}

// LogBuildsManualTriggerInstructions affiche les instructions pour le déclenchement manuel
func LogBuildsManualTriggerInstructions(logger *log.Logger) {
	logger.Printf("[INFO] To trigger builds workflows manually via code, you can use:")
	logger.Printf("[INFO] - Rankings: scheduleManager.TriggerRankingsNow(ctx)")
	logger.Printf("[INFO] - Reports for specific class: scheduleManager.TriggerReportsForClassNow(ctx, \"ClassName\")")
	logger.Printf("[INFO] - Builds: scheduleManager.TriggerBuildsNow(ctx)")
	logger.Printf("[INFO] - Equipment Analysis: scheduleManager.TriggerEquipmentAnalysisNow(ctx)")
	logger.Printf("[INFO] - Talent Analysis: scheduleManager.TriggerTalentAnalysisNow(ctx)")
	logger.Printf("[INFO] - Stat Analysis: scheduleManager.TriggerStatAnalysisNow(ctx)")
}
