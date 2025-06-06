// internal/services/warcraftlogs/mythicplus/player_rankings/temporal/scheduler/init.go
package warcraftlogsPlayerRankingsTemporalScheduler

import (
	"context"
	"log"

	playerRankingsDefinitions "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/definitions"
)

// InitPlayerRankingsSchedule initialise le schedule pour le player rankings workflow
func InitPlayerRankingsSchedule(ctx context.Context, scheduleManager *PlayerRankingsScheduleManager, configPath string, opts *ScheduleOptions, logger *log.Logger) error {
	// Charger les paramètres du workflow player rankings
	playerRankingsParams, err := playerRankingsDefinitions.LoadPlayerRankingsParams(configPath)
	if err != nil {
		logger.Printf("[ERROR] Failed to load player rankings params: %v", err)
		return err
	}

	// Créer le schedule (utilise la configuration par défaut si non spécifiée)
	if err := scheduleManager.CreatePlayerRankingsSchedule(ctx, playerRankingsParams, nil, opts); err != nil {
		logger.Printf("[ERROR] Failed to create player rankings schedule: %v", err)
		return err
	}

	logger.Printf("[INFO] Successfully created player rankings schedule with batch ID: %s", playerRankingsParams.BatchID)
	return nil
}

// LogPlayerRankingsManualTriggerInstructions affiche les instructions pour le déclenchement manuel
func LogPlayerRankingsManualTriggerInstructions(logger *log.Logger) {
	logger.Printf("[INFO] To trigger player rankings workflow manually via code, you can use:")
	logger.Printf("[INFO] - PlayerRankings: playerRankingsScheduleManager.TriggerPlayerRankingsNow(ctx)")
}
