// internal/services/warcraftlogs/mythicplus/player_rankings/temporal/init.go
package temporal

import (
	"wowperf/internal/services/warcraftlogs"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"gorm.io/gorm"

	playerRankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/repository"
	playerRankingsActivities "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/activities"
	playerRankingsDefinitions "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/definitions"
	playerRankingsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/player_rankings"
)

// InitPlayerRankings initialise les repositories et activités pour la feature
func InitPlayerRankings(db *gorm.DB, warcraftLogsClient *warcraftlogs.WarcraftLogsClientService) (*playerRankingsRepository.PlayerRankingsRepository, *playerRankingsActivities.Activities) {
	// Initialiser le repository
	playerRankingsRepo := playerRankingsRepository.NewPlayerRankingsRepository(db)

	// Initialiser les activités
	playerRankingsActivity := playerRankingsActivities.NewPlayerRankingsActivity(
		warcraftLogsClient,
		playerRankingsRepo,
	)

	activitiesService := playerRankingsActivities.NewActivities(playerRankingsActivity)

	return playerRankingsRepo, activitiesService
}

// RegisterPlayerRankings enregistre les workflows et activités avec le worker
func RegisterPlayerRankings(w worker.Worker, activitiesService *playerRankingsActivities.Activities) {
	// Enregistrer le workflow
	playerRankingsWorkflowImpl := playerRankingsWorkflow.NewPlayerRankingsWorkflow()
	w.RegisterWorkflowWithOptions(playerRankingsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: playerRankingsDefinitions.PlayerRankingsWorkflowName,
	})

	// Enregistrer les activités
	w.RegisterActivity(activitiesService.PlayerRankings.FetchAllDungeonRankings)
	w.RegisterActivity(activitiesService.PlayerRankings.StoreRankings)
	w.RegisterActivity(activitiesService.PlayerRankings.CalculateDailyMetrics)
}
