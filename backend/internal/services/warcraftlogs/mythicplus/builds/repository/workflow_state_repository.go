package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// WorkflowStateRepository handles database operations for workflow states
type WorkflowStateRepository struct {
	db *gorm.DB
}

// NewWorkflowStateRepository creates a new instance of WorkflowStateRepository
func NewWorkflowStateRepository(db *gorm.DB) *WorkflowStateRepository {
	return &WorkflowStateRepository{
		db: db,
	}
}

// CreateWorkflowState persists a new workflow state to the database
func (r *WorkflowStateRepository) CreateWorkflowState(ctx context.Context, state *warcraftlogsBuilds.WorkflowState) error {
	if state.CreatedAt.IsZero() {
		state.CreatedAt = time.Now()
	}
	if state.UpdatedAt.IsZero() {
		state.UpdatedAt = time.Now()
	}

	result := r.db.WithContext(ctx).Create(state)
	if result.Error != nil {
		return fmt.Errorf("failed to create workflow state: %w", result.Error)
	}
	return nil
}

// UpdateWorkflowState updates an existing workflow state
func (r *WorkflowStateRepository) UpdateWorkflowState(ctx context.Context, state *warcraftlogsBuilds.WorkflowState) error {
	state.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(state)
	if result.Error != nil {
		return fmt.Errorf("failed to update workflow state: %w", result.Error)
	}
	return nil
}

// GetLastWorkflowRun retrieves the most recent workflow run for a specific type
func (r *WorkflowStateRepository) GetLastWorkflowRun(ctx context.Context, workflowType string) (*warcraftlogsBuilds.WorkflowState, error) {
	var state warcraftlogsBuilds.WorkflowState

	result := r.db.WithContext(ctx).
		Where("workflow_type = ?", workflowType).
		Order("created_at DESC").
		First(&state)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last workflow run: %w", result.Error)
	}
	return &state, nil
}

// GetWorkflowStatistics retrieves statistics for a specific workflow type within a time period
func (r *WorkflowStateRepository) GetWorkflowStatistics(ctx context.Context, workflowType string, days int) ([]*warcraftlogsBuilds.WorkflowState, error) {
	var states []*warcraftlogsBuilds.WorkflowState

	result := r.db.WithContext(ctx).
		Where("workflow_type = ? AND created_at > ?",
			workflowType,
			time.Now().Add(-time.Hour*24*time.Duration(days))).
		Order("created_at DESC").
		Find(&states)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get workflow statistics: %w", result.Error)
	}

	return states, nil
}

// GetWorkflowStateByID retrieves a specific workflow state by its ID
func (r *WorkflowStateRepository) GetWorkflowStateByID(ctx context.Context, id string) (*warcraftlogsBuilds.WorkflowState, error) {
	var state warcraftlogsBuilds.WorkflowState

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&state)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get workflow state by ID: %w", result.Error)
	}

	return &state, nil
}

// DeleteWorkflowStates deletes workflow states older than a certain number of days
func (r *WorkflowStateRepository) DeleteWorkflowStates(ctx context.Context, workflowType string, olderThanDays int) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("workflow_type = ? AND created_at < ?",
			workflowType,
			time.Now().Add(-time.Hour*24*time.Duration(olderThanDays))).
		Delete(&warcraftlogsBuilds.WorkflowState{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old workflow states: %w", result.Error)
	}

	return result.RowsAffected, nil
}
