package warcraftlogsBuildsTemporalScheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
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
	scheduleID := fmt.Sprintf("warcraft-logs-%s-schedule", className)
	workflowID := fmt.Sprintf("warcraft-logs-%s-workflow", className)

	// Filter specs for this class
	var classSpecs []workflows.ClassSpec
	for _, spec := range cfg.Specs {
		if spec.ClassName == className {
			classSpecs = append(classSpecs, spec)
		}
	}

	if len(classSpecs) == 0 {
		return fmt.Errorf("no specs found for class %s", className)
	}

	// Prepare workflow parameters for this class
	configCopy := new(workflows.Config)
	*configCopy = *cfg
	configCopy.Specs = classSpecs

	workflowParams := workflows.WorkflowParams{
		Specs:       classSpecs,
		Dungeons:    configCopy.Dungeons,
		BatchConfig: configCopy.Rankings.Batch,
		Rankings: struct {
			MaxRankingsPerSpec int           `json:"max_rankings_per_spec"`
			UpdateInterval     time.Duration `json:"update_interval"`
		}{
			MaxRankingsPerSpec: configCopy.Rankings.MaxRankingsPerSpec,
			UpdateInterval:     configCopy.Rankings.UpdateInterval,
		},
		Config: configCopy,
	}

	// Create schedule options according to the documentation
	scheduleOptions := client.ScheduleOptions{
		ID: scheduleID,

		// Schedule Spec - Define when to run (Tuesday 7am)
		Spec: client.ScheduleSpec{
			CronExpressions: []string{opts.Policy.CronExpression},
		},

		// Workflow action to execute
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  workflows.SyncWorkflow,
			TaskQueue: "warcraft-logs-sync",
			Args:      []interface{}{workflowParams},
			// The timeouts must be defined here
			WorkflowRunTimeout:  opts.Timeout,
			WorkflowTaskTimeout: 10 * time.Second, // Default recommended value
		},
	}

	// Create the schedule
	handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	sm.schedules[className] = handle
	sm.logger.Printf("Created schedule for class %s with ID %s", className, scheduleID)
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
