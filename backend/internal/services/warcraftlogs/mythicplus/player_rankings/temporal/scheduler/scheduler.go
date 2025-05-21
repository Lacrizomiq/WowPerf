package warcraftlogsPlayerRankingsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"

	definitions "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/models"
)

// Constantes pour les IDs des schedules
const (
	playerRankingsScheduleID = "player-rankings-daily"
)

// PlayerRankingsScheduleManager gère le schedule Temporal pour le workflow PlayerRankings
type PlayerRankingsScheduleManager struct {
	client                 client.Client
	logger                 *log.Logger
	playerRankingsSchedule client.ScheduleHandle
}

// NewPlayerRankingsScheduleManager crée une nouvelle instance de PlayerRankingsScheduleManager
func NewPlayerRankingsScheduleManager(temporalClient client.Client, logger *log.Logger) *PlayerRankingsScheduleManager {
	return &PlayerRankingsScheduleManager{
		client: temporalClient,
		logger: logger,
	}
}

// CreatePlayerRankingsSchedule crée le schedule pour le PlayerRankingsWorkflow
func (sm *PlayerRankingsScheduleManager) CreatePlayerRankingsSchedule(ctx context.Context, params *models.PlayerRankingWorkflowParams, config *PlayerRankingsScheduleConfig, opts *ScheduleOptions) error {
	if config == nil {
		config = &DefaultPlayerRankingsScheduleConfig
	}

	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	// Valider la configuration
	if err := ValidateScheduleConfig(config); err != nil {
		return fmt.Errorf("invalid schedule configuration: %w", err)
	}

	scheduleID := playerRankingsScheduleID

	// Générer un BatchID unique pour ce schedule
	if params.BatchID == "" {
		params.BatchID = fmt.Sprintf("player-rankings-%s", uuid.New().String())
	}

	// Définir le CRON pour une exécution quotidienne à l'heure spécifiée
	// Format: minute heure * * *
	cronExpression := fmt.Sprintf("%d %d * * *", config.Minute, config.Hour)

	workflowID := fmt.Sprintf("player-rankings-%s", time.Now().UTC().Format("2006-01-02"))

	// Créer le schedule
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{"0 12 * * *"},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  definitions.PlayerRankingsWorkflowName,
			TaskQueue: config.TaskQueue,
			Args:      []interface{}{params},
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
			WorkflowRunTimeout: opts.Timeout,
		},
		Paused: opts.Paused, // En pause par défaut si spécifié dans les options
	}

	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create player rankings schedule: %w", err)
	}

	sm.playerRankingsSchedule = handle
	sm.logger.Printf("[INFO] Created player rankings schedule: %s (cron: %s)", scheduleID, cronExpression)
	return nil
}

// TriggerPlayerRankingsNow déclenche l'exécution immédiate du schedule
func (sm *PlayerRankingsScheduleManager) TriggerPlayerRankingsNow(ctx context.Context) error {
	if sm.playerRankingsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.playerRankingsSchedule = sm.client.ScheduleClient().GetHandle(ctx, playerRankingsScheduleID)
	}

	return sm.playerRankingsSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// PausePlayerRankingsSchedule met en pause le schedule
func (sm *PlayerRankingsScheduleManager) PausePlayerRankingsSchedule(ctx context.Context) error {
	if sm.playerRankingsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.playerRankingsSchedule = sm.client.ScheduleClient().GetHandle(ctx, playerRankingsScheduleID)
	}

	return sm.playerRankingsSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// UnpausePlayerRankingsSchedule réactive le schedule
func (sm *PlayerRankingsScheduleManager) UnpausePlayerRankingsSchedule(ctx context.Context) error {
	if sm.playerRankingsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.playerRankingsSchedule = sm.client.ScheduleClient().GetHandle(ctx, playerRankingsScheduleID)
	}

	return sm.playerRankingsSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// DeletePlayerRankingsSchedule supprime le schedule
func (sm *PlayerRankingsScheduleManager) DeletePlayerRankingsSchedule(ctx context.Context) error {
	if sm.playerRankingsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.playerRankingsSchedule = sm.client.ScheduleClient().GetHandle(ctx, playerRankingsScheduleID)
	}

	err := sm.playerRankingsSchedule.Delete(ctx)
	if err == nil {
		sm.playerRankingsSchedule = nil
		sm.logger.Printf("[INFO] Deleted player rankings schedule: %s", playerRankingsScheduleID)
	}
	return err
}

// CleanupPlayerRankingsWorkflows termine tous les workflows PlayerRankings en cours d'exécution
func (sm *PlayerRankingsScheduleManager) CleanupPlayerRankingsWorkflows(ctx context.Context) error {
	var terminatedCount int

	// Définir le type de workflow à nettoyer
	workflowType := definitions.PlayerRankingsWorkflowName

	// Construire la requête pour ce type de workflow
	query := fmt.Sprintf("WorkflowType='%s'", workflowType)

	sm.logger.Printf("[INFO] Listing workflows of type: %s", workflowType)

	// Récupérer les workflows de ce type
	resp, err := sm.client.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: "default", // Utilisez votre namespace si différent
		Query:     query,
	})

	if err != nil {
		return fmt.Errorf("failed to list workflows of type %s: %w", workflowType, err)
	}

	// Traiter chaque workflow récupéré
	for _, execution := range resp.Executions {
		workflowID := execution.Execution.WorkflowId
		runID := execution.Execution.RunId

		// Ne terminer que les workflows en cours d'exécution
		if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			sm.logger.Printf("[INFO] Skipping non-running workflow: %s (status: %s)",
				workflowID, execution.Status.String())
			continue
		}

		// Terminer le workflow
		err := sm.client.TerminateWorkflow(ctx, workflowID, runID, "Cleanup of workflows during service restart")
		if err != nil {
			sm.logger.Printf("[WARN] Failed to terminate workflow %s: %v", workflowID, err)
			continue
		}

		sm.logger.Printf("[INFO] Terminated workflow: %s", workflowID)
		terminatedCount++
	}

	sm.logger.Printf("[INFO] Workflow cleanup completed - terminated %d workflows", terminatedCount)
	return nil
}

// CleanupAll effectue un nettoyage complet
func (sm *PlayerRankingsScheduleManager) CleanupAll(ctx context.Context) error {
	// Supprimer d'abord le schedule
	if err := sm.DeletePlayerRankingsSchedule(ctx); err != nil {
		sm.logger.Printf("[WARN] Failed to delete schedule: %v", err)
		// Continuer malgré les erreurs
	}

	// Ensuite nettoyer les workflows
	if err := sm.CleanupPlayerRankingsWorkflows(ctx); err != nil {
		return fmt.Errorf("failed to cleanup workflows: %w", err)
	}

	return nil
}
