package raiderioMythicPlusRunsTemporalScheduler

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

	definitions "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/definitions"
	models "wowperf/internal/services/raiderio/mythicplus/mythicplus_runs/temporal/workflows/models"
)

// Constantes pour les IDs des schedules
const (
	mythicPlusRunsScheduleID = "mythicplus-runs-daily"
)

// MythicPlusRunsScheduleManager gère le schedule Temporal pour le workflow MythicPlusRuns
type MythicPlusRunsScheduleManager struct {
	client                 client.Client
	logger                 *log.Logger
	mythicPlusRunsSchedule client.ScheduleHandle
}

// NewMythicPlusRunsScheduleManager crée une nouvelle instance de MythicPlusRunsScheduleManager
func NewMythicPlusRunsScheduleManager(temporalClient client.Client, logger *log.Logger) *MythicPlusRunsScheduleManager {
	return &MythicPlusRunsScheduleManager{
		client: temporalClient,
		logger: logger,
	}
}

// CreateMythicPlusRunsSchedule crée le schedule pour le MythicPlusRunsWorkflow
func (sm *MythicPlusRunsScheduleManager) CreateMythicPlusRunsSchedule(ctx context.Context, params *models.MythicRunsWorkflowParams, config *MythicPlusRunsScheduleConfig, opts *ScheduleOptions) error {
	if config == nil {
		config = &DefaultMythicPlusRunsScheduleConfig
	}

	if opts == nil {
		opts = DefaultScheduleOptions()
	}

	// Valider la configuration
	if err := ValidateScheduleConfig(config); err != nil {
		return fmt.Errorf("invalid schedule configuration: %w", err)
	}

	scheduleID := mythicPlusRunsScheduleID

	// Générer un BatchID unique pour ce schedule
	if params.BatchID == "" {
		params.BatchID = fmt.Sprintf("mythicplus-runs-%s", uuid.New().String())
	}

	// Définir le CRON pour une exécution quotidienne à l'heure spécifiée
	// Format: minute heure * * *
	cronExpression := fmt.Sprintf("%d %d %d * *", config.Minute, config.Hour, config.Day)

	workflowID := fmt.Sprintf("mythicplus-runs-%s", time.Now().UTC().Format("2006-01-02"))

	// Créer le schedule
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{cronExpression},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  definitions.MythicPlusRunsWorkflowName,
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
		return fmt.Errorf("failed to create mythicplus runs schedule: %w", err)
	}

	sm.mythicPlusRunsSchedule = handle
	sm.logger.Printf("[INFO] Created mythicplus runs schedule: %s (cron: %s)", scheduleID, cronExpression)
	return nil
}

// TriggerMythicPlusRunsNow déclenche l'exécution immédiate du schedule
func (sm *MythicPlusRunsScheduleManager) TriggerMythicPlusRunsNow(ctx context.Context) error {
	if sm.mythicPlusRunsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.mythicPlusRunsSchedule = sm.client.ScheduleClient().GetHandle(ctx, mythicPlusRunsScheduleID)
	}

	return sm.mythicPlusRunsSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}

// PauseMythicPlusRunsSchedule met en pause le schedule
func (sm *MythicPlusRunsScheduleManager) PauseMythicPlusRunsSchedule(ctx context.Context) error {
	if sm.mythicPlusRunsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.mythicPlusRunsSchedule = sm.client.ScheduleClient().GetHandle(ctx, mythicPlusRunsScheduleID)
	}

	return sm.mythicPlusRunsSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

// UnpauseMythicPlusRunsSchedule réactive le schedule
func (sm *MythicPlusRunsScheduleManager) UnpauseMythicPlusRunsSchedule(ctx context.Context) error {
	if sm.mythicPlusRunsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.mythicPlusRunsSchedule = sm.client.ScheduleClient().GetHandle(ctx, mythicPlusRunsScheduleID)
	}

	return sm.mythicPlusRunsSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}

// DeleteMythicPlusRunsSchedule supprime le schedule
func (sm *MythicPlusRunsScheduleManager) DeleteMythicPlusRunsSchedule(ctx context.Context) error {
	if sm.mythicPlusRunsSchedule == nil {
		// Essayer de récupérer le handle si déjà créé
		sm.mythicPlusRunsSchedule = sm.client.ScheduleClient().GetHandle(ctx, mythicPlusRunsScheduleID)
	}

	err := sm.mythicPlusRunsSchedule.Delete(ctx)
	if err == nil {
		sm.mythicPlusRunsSchedule = nil
		sm.logger.Printf("[INFO] Deleted mythicplus runs schedule: %s", mythicPlusRunsScheduleID)
	}
	return err
}

// CleanupMythicPlusRunsWorkflows termine tous les workflows MythicPlusRuns en cours d'exécution
func (sm *MythicPlusRunsScheduleManager) CleanupMythicPlusRunsWorkflows(ctx context.Context) error {
	var terminatedCount int

	// Définir le type de workflow à nettoyer
	workflowType := definitions.MythicPlusRunsWorkflowName

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
func (sm *MythicPlusRunsScheduleManager) CleanupAll(ctx context.Context) error {
	// Supprimer d'abord le schedule
	if err := sm.DeleteMythicPlusRunsSchedule(ctx); err != nil {
		sm.logger.Printf("[WARN] Failed to delete schedule: %v", err)
		// Continuer malgré les erreurs
	}

	// Ensuite nettoyer les workflows
	if err := sm.CleanupMythicPlusRunsWorkflows(ctx); err != nil {
		return fmt.Errorf("failed to cleanup workflows: %w", err)
	}

	return nil
}
