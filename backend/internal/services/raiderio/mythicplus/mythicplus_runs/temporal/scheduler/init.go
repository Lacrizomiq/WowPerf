// internal/services/raiderio/mythicplus/mythicplus_runs/temporal/scheduler/init.go
package raiderioMythicPlusRunsTemporalScheduler

import (
	"context"
	"log"

	definitions "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/definitions"
)

// InitMythicPlusRunsSchedule initialise le schedule pour le mythicplus runs workflow
func InitMythicPlusRunsSchedule(ctx context.Context, scheduleManager *MythicPlusRunsScheduleManager, configPath string, opts *ScheduleOptions, logger *log.Logger) error {
	// Charger les paramètres du workflow mythicplus runs
	mythicPlusRunsParams, err := definitions.LoadMythicRunsParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load mythicplus runs params: %v", err)
		return err
	}

	// Créer le schedule (utilise la configuration par défaut si non spécifiée)
	if err := scheduleManager.CreateMythicPlusRunsSchedule(ctx, mythicPlusRunsParams, nil, opts); err != nil {
		logger.Printf("[ERROR] Failed to create mythicplus runs schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created mythicplus runs schedule with batch ID: %s", mythicPlusRunsParams.BatchID)
	return nil
}

// LogMythicPlusRunsManualTriggerInstructions affiche les instructions pour le déclenchement manuel
func LogMythicPlusRunsManualTriggerInstructions(logger *log.Logger) {
	logger.Printf("[INFO] To trigger mythicplus runs workflow manually via code, you can use:")
	logger.Printf("[INFO] - MythicPlusRuns: mythicPlusRunsScheduleManager.TriggerMythicPlusRunsNow(ctx)")
}
