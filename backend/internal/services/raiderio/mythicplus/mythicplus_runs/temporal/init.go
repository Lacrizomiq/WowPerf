// internal/services/raiderio/mythicplus/mythicplus_runs/temporal/init.go
package temporal

import (
	"wowperf/internal/services/raiderio"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"gorm.io/gorm"

	mythicPlusRunsRepository "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/repository"
	mythicPlusRunsActivities "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/activities"
	mythicPlusRunsDefinitions "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/definitions"
	mythicPlusRunsWorkflow "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/mythicplus_runs"
)

// InitMythicPlusRuns initialise les repositories et activités pour la feature
func InitMythicPlusRuns(db *gorm.DB, raiderIOClient *raiderio.RaiderIOService) (*mythicPlusRunsRepository.MythicPlusRunsRepository, *mythicPlusRunsActivities.MythicPlusRunsActivity) {
	// Initialiser le repository
	mythicPlusRunsRepo := mythicPlusRunsRepository.NewMythicPlusRunsRepository(db)

	// Initialiser les activités
	mythicPlusRunsActivity := mythicPlusRunsActivities.NewMythicPlusRunsActivity(
		raiderIOClient,
		mythicPlusRunsRepo,
	)

	return mythicPlusRunsRepo, mythicPlusRunsActivity
}

// RegisterMythicPlusRuns enregistre les workflows et activités avec le worker
func RegisterMythicPlusRuns(w worker.Worker, activitiesService *mythicPlusRunsActivities.MythicPlusRunsActivity) {
	// Enregistrer le workflow
	mythicPlusRunsWorkflowImpl := mythicPlusRunsWorkflow.NewMythicPlusRunsWorkflow()
	w.RegisterWorkflowWithOptions(mythicPlusRunsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: mythicPlusRunsDefinitions.MythicPlusRunsWorkflowName,
	})

	// Enregistrer les activités
	w.RegisterActivity(activitiesService.FetchAndProcessMythicPlusRunsActivity)
}
