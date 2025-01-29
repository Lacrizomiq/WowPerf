package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

// ScheduleManager handles Temporal schedules for different classes
type ScheduleManager struct {
	client    client.Client
	schedules map[string]client.ScheduleHandle
	logger    *log.Logger
}

// NewScheduleManager creates a new ScheduleManager instance
func NewScheduleManager(temporalClient client.Client, logger *log.Logger) *ScheduleManager {
	return &ScheduleManager{
		client:    temporalClient,
		schedules: make(map[string]client.ScheduleHandle),
		logger:    logger,
	}
}

// CreateClassSchedule creates a schedule for a specific class
func (sm *ScheduleManager) CreateClassSchedule(
	ctx context.Context,
	className string,
	cfg *workflows.Config,
	opts *ScheduleOptions,
) error {

	sm.logger.Printf("Creating schedule for class: %s", className)

	scheduleID := fmt.Sprintf("warcraft-logs-%s-schedule", className)
	workflowID := fmt.Sprintf("warcraft-logs-%s-workflow", className)

	// Prepare simplified workflow parameters
	workflowParams := workflows.WorkflowParams{
		Config:     cfg,
		WorkflowID: workflowID,
	}

	sm.logger.Printf("Creating schedule with params - ConfigSpecs: %d",
		len(workflowParams.Config.Specs))

	// Create schedule options
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{opts.Policy.CronExpression},
			// Add jitter pour éviter les démarrages simultanés
			Jitter: time.Minute,
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  workflows.SyncWorkflow,
			TaskQueue: "warcraft-logs-sync",
			Args:      []interface{}{workflowParams},
			// Timeouts adaptés par classe
			WorkflowRunTimeout:  opts.Timeout,
			WorkflowTaskTimeout: time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    opts.Retry.InitialInterval,
				BackoffCoefficient: opts.Retry.BackoffCoefficient,
				MaximumInterval:    opts.Retry.MaximumInterval,
				MaximumAttempts:    int32(opts.Retry.MaximumAttempts),
			},
		},
	}

	// Create the schedule
	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create schedule for %s: %w", className, err)
	}

	sm.schedules[className] = handle
	sm.logger.Printf("Created schedule for class %s with ID %s", className, scheduleID)

	sm.monitorSchedule(ctx, handle, className)
	return nil
}

// PauseSchedule pauses an existing schedule
func (sm *ScheduleManager) PauseSchedule(ctx context.Context, className string) error {
	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	if err := handle.Pause(ctx, client.SchedulePauseOptions{}); err != nil {
		return fmt.Errorf("failed to pause schedule: %w", err)
	}

	sm.logger.Printf("Paused schedule for class %s", className)
	return nil
}

// ResumeSchedule resumes a paused schedule
func (sm *ScheduleManager) ResumeSchedule(ctx context.Context, className string) error {
	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	if err := handle.Unpause(ctx, client.ScheduleUnpauseOptions{}); err != nil {
		return fmt.Errorf("failed to resume schedule: %w", err)
	}

	sm.logger.Printf("Resumed schedule for class %s", className)
	return nil
}

// DeleteSchedule deletes an existing schedule
func (sm *ScheduleManager) DeleteSchedule(ctx context.Context, className string) error {
	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	if err := handle.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	delete(sm.schedules, className)
	sm.logger.Printf("Deleted schedule for class %s", className)
	return nil
}

// BackfillSchedule triggers missed executions for a specific time range
func (sm *ScheduleManager) BackfillSchedule(ctx context.Context, className string, start, end time.Time) error {
	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	err := handle.Backfill(ctx, client.ScheduleBackfillOptions{
		Backfill: []client.ScheduleBackfill{
			{
				Start:   start,
				End:     end,
				Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to backfill schedule: %w", err)
	}

	sm.logger.Printf("Backfilled schedule for class %s from %v to %v",
		className, start.Format(time.RFC3339), end.Format(time.RFC3339))
	return nil
}

// TriggerSchedule triggers an immediate execution of the schedule
func (sm *ScheduleManager) TriggerSchedule(ctx context.Context, className string) error {
	handle, exists := sm.schedules[className]
	if !exists {
		return fmt.Errorf("no schedule found for class %s", className)
	}

	err := handle.Trigger(ctx, client.ScheduleTriggerOptions{
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL,
	})
	if err != nil {
		return fmt.Errorf("failed to trigger schedule: %w", err)
	}

	sm.logger.Printf("Triggered immediate execution for class %s", className)
	return nil
}

// ListSchedules returns the list of active schedules
func (sm *ScheduleManager) ListSchedules(ctx context.Context) []string {
	var scheduleIDs []string
	for className := range sm.schedules {
		scheduleIDs = append(scheduleIDs, className)
	}
	return scheduleIDs
}

// GetScheduleDescription retrieves the details of a schedule
func (sm *ScheduleManager) GetScheduleDescription(ctx context.Context, className string) (*client.ScheduleDescription, error) {
	handle, exists := sm.schedules[className]
	if !exists {
		return nil, fmt.Errorf("no schedule found for class %s", className)
	}

	return handle.Describe(ctx)
}

// CreateOrGetClassSchedule creates a new schedule if it doesn't exist, or returns existing one
func (sm *ScheduleManager) CreateOrGetClassSchedule(ctx context.Context, className string, cfg *workflows.Config, opts *ScheduleOptions) error {
	sm.logger.Printf("CreateOrGetClassSchedule called with className: %s", className)
	scheduleID := fmt.Sprintf("warcraft-logs-%s-schedule", className)

	// Check if schedule already exists using Temporal's List method
	listView, err := sm.client.ScheduleClient().List(ctx, client.ScheduleListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list schedules: %w", err)
	}

	for listView.HasNext() {
		schedule, err := listView.Next()
		if err != nil {
			sm.logger.Printf("Error listing schedule: %v", err)
			continue
		}

		if schedule.ID == scheduleID {
			sm.logger.Printf("Schedule already exists for class %s", className)
			// Get the handle for the existing schedule
			handle := sm.client.ScheduleClient().GetHandle(ctx, scheduleID)
			sm.schedules[className] = handle
			return nil
		}
	}

	// If we get here, schedule doesn't exist, create it
	return sm.CreateClassSchedule(ctx, className, cfg, opts)
}

// monitorSchedule monitors the schedule and logs basic information
func (sm *ScheduleManager) monitorSchedule(ctx context.Context,
	handle client.ScheduleHandle,
	className string) {

	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				desc, err := handle.Describe(ctx)
				if err != nil {
					sm.logger.Printf("Error monitoring schedule for %s: %v",
						className, err)
					continue
				}

				// Log schedule info using available information
				sm.logger.Printf("Schedule monitoring for class %s: Schedule ID: %s, State: %v, Last checked: %v",
					className,
					handle.GetID(),
					desc.Schedule.State,
					time.Now().Format(time.RFC3339))
			}
		}
	}()
}
