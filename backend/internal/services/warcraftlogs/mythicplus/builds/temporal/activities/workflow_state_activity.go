package warcraftlogsBuildsTemporalActivities

import (
	"context"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	repository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	"go.temporal.io/sdk/activity"
)

// WorkflowStateActivity is an activity for managing workflow state
type WorkflowStateActivity struct {
	repository *repository.WorkflowStateRepository
}

// NewWorkflowStateActivity creates a new WorkflowStateActivity
func NewWorkflowStateActivity(repo *repository.WorkflowStateRepository) *WorkflowStateActivity {
	return &WorkflowStateActivity{repository: repo}
}

// CreateWorkflowState creates a new workflow state
func (a *WorkflowStateActivity) CreateWorkflowState(ctx context.Context, state *warcraftlogsBuilds.WorkflowState) (*warcraftlogsBuilds.WorkflowState, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Creating workflow state",
		"id", state.ID,
		"type", state.WorkflowType)

	if err := a.repository.CreateWorkflowState(ctx, state); err != nil {
		logger.Error("Failed to create workflow state", "error", err)
		return nil, err
	}

	return state, nil
}

// UpdateWorkflowState updates an existing workflow state
func (a *WorkflowStateActivity) UpdateWorkflowState(ctx context.Context, state *warcraftlogsBuilds.WorkflowState) error {
	logger := activity.GetLogger(ctx)

	logger.Info("Updating workflow state",
		"id", state.ID,
		"type", state.WorkflowType,
		"status", state.Status)

	if err := a.repository.UpdateWorkflowState(ctx, state); err != nil {
		logger.Error("Failed to update workflow state", "error", err)
		return err
	}

	return nil
}

// GetLastWorkflowRun récupère le dernier run d'un type de workflow
func (a *WorkflowStateActivity) GetLastWorkflowRun(ctx context.Context, workflowType string) (*warcraftlogsBuilds.WorkflowState, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Getting last workflow run", "type", workflowType)

	state, err := a.repository.GetLastWorkflowRun(ctx, workflowType)
	if err != nil {
		logger.Error("Failed to get last workflow run", "error", err)
		return nil, err
	}

	return state, nil
}

// GetWorkflowStatistics récupère des statistiques sur les runs de workflow
func (a *WorkflowStateActivity) GetWorkflowStatistics(ctx context.Context, workflowType string, days int) ([]*warcraftlogsBuilds.WorkflowState, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Getting workflow statistics",
		"type", workflowType,
		"days", days)

	states, err := a.repository.GetWorkflowStatistics(ctx, workflowType, days)
	if err != nil {
		logger.Error("Failed to get workflow statistics", "error", err)
		return nil, err
	}

	return states, nil
}

// GetWorkflowStateByID récupère un état de workflow par son ID
func (a *WorkflowStateActivity) GetWorkflowStateByID(ctx context.Context, id string) (*warcraftlogsBuilds.WorkflowState, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Getting workflow state by ID", "id", id)

	state, err := a.repository.GetWorkflowStateByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get workflow state by ID", "error", err)
		return nil, err
	}

	return state, nil
}

// DeleteOldWorkflowStates supprime les anciens états de workflow
func (a *WorkflowStateActivity) DeleteOldWorkflowStates(ctx context.Context, workflowType string, olderThanDays int) (int64, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Deleting old workflow states",
		"type", workflowType,
		"olderThanDays", olderThanDays)

	count, err := a.repository.DeleteWorkflowStates(ctx, workflowType, olderThanDays)
	if err != nil {
		logger.Error("Failed to delete old workflow states", "error", err)
		return 0, err
	}

	logger.Info("Deleted old workflow states", "count", count)
	return count, nil
}
