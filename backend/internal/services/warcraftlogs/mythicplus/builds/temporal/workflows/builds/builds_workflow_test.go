package warcraftlogsBuildsTemporalWorkflowsBuilds

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *WorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.env.RegisterWorkflow(WorkflowProcessBuilds)
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessBuilds_Success() {
	s.T().Log("Testing successful builds processing workflow")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{
			NumWorkers:   2,
			RequestDelay: time.Millisecond * 100,
		},
	}

	// Mock GetReportsBatch activity
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, batchSize, offset int32) ([]*warcraftlogsBuilds.Report, error) {
			return []*warcraftlogsBuilds.Report{
				{
					Code:             "ABC123",
					FightID:          1,
					TalentCodes:      []byte(`{"Priest_Shadow_talents":"test-code"}`),
					PlayerDetailsDps: []byte(`[{"name":"TestPlayer","class":"Priest","spec":"Shadow"}]`),
				},
			}, nil
		},
		activity.RegisterOptions{Name: definitions.GetReportsBatchActivity},
	)

	// Mock ProcessBuilds activity
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: int32(len(reports)),
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildsActivity},
	)

	// Execute workflow
	var result models.WorkflowResult
	s.env.ExecuteWorkflow(WorkflowProcessBuilds, config)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(1), result.BuildsProcessed)
	s.NotZero(result.StartedAt)
	s.NotZero(result.CompletedAt)

	s.T().Logf("Successfully processed %d builds", result.BuildsProcessed)
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessBuilds_NoReports() {
	s.T().Log("Testing workflow with no reports")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, batchSize, offset int32) ([]*warcraftlogsBuilds.Report, error) {
			return []*warcraftlogsBuilds.Report{}, nil
		},
		activity.RegisterOptions{Name: definitions.GetReportsBatchActivity},
	)

	var result models.WorkflowResult
	s.env.ExecuteWorkflow(WorkflowProcessBuilds, config)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Zero(result.BuildsProcessed)

	s.T().Log("Successfully handled no reports case")
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessBuilds_DatabaseError() {
	s.T().Log("Testing database error handling")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
	}

	// Return some reports from database
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, batchSize, offset int32) ([]*warcraftlogsBuilds.Report, error) {
			return []*warcraftlogsBuilds.Report{
				{Code: "TEST123"},
			}, nil
		},
		activity.RegisterOptions{Name: definitions.GetReportsBatchActivity},
	)

	// Simulate database error during build processing
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeDatabase,
				Message:   "Database connection failed",
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildsActivity},
	)

	s.env.ExecuteWorkflow(WorkflowProcessBuilds, config)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Database connection failed")

	s.T().Log("Successfully handled database error")
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessBuilds_Cancellation() {
	s.T().Log("Testing workflow cancellation handling")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, batchSize, offset int32) ([]*warcraftlogsBuilds.Report, error) {
			time.Sleep(time.Millisecond * 100)
			return []*warcraftlogsBuilds.Report{{Code: "TEST123"}}, nil
		},
		activity.RegisterOptions{Name: definitions.GetReportsBatchActivity},
	)

	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Millisecond*50)

	s.env.ExecuteWorkflow(WorkflowProcessBuilds, config)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.True(temporal.IsCanceledError(err))

	s.T().Log("Successfully handled workflow cancellation")
}

// WorkflowProcessBuilds is a test wrapper that sets up workflow with proper activity options
func WorkflowProcessBuilds(ctx workflow.Context, config models.WorkflowConfig) (*models.WorkflowResult, error) {
	wf := &BuildsWorkflow{
		stateManager: state.NewManager(),
		processor:    NewProcessor(),
	}

	// Initialize state
	state := wf.stateManager.GetState()
	if state.StartedAt.IsZero() {
		state.StartedAt = workflow.Now(ctx)
	}

	// Set activity options
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	return wf.Execute(ctx, config)
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
