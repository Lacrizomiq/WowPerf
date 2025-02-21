package warcraftlogsBuildsTemporalWorkflowsReports

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
	s.env.RegisterWorkflow(WorkflowProcessReports)
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessReports_Success() {
	s.T().Log("Testing successful reports processing")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{
			NumWorkers:   2,
			RequestDelay: time.Millisecond * 100,
		},
		Specs: []models.ClassSpec{
			{ClassName: "Priest", SpecName: "Shadow"},
		},
		Dungeons: []models.Dungeon{
			{Name: "TestDungeon", EncounterID: 12660},
		},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
			return []*warcraftlogsBuilds.ClassRanking{
				{
					PlayerName:    "TestPlayer1",
					Class:         "Priest",
					Spec:          "Shadow",
					ReportCode:    "ABC123",
					EncounterID:   12660,
					ReportFightID: 1,
				},
			}, nil
		},
		activity.RegisterOptions{Name: definitions.GetStoredRankingsActivity},
	)

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking, worker models.WorkerConfig) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: 1,
				ClassName:      "Priest",
				SpecName:       "Shadow",
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	var result models.WorkflowResult
	s.env.ExecuteWorkflow(WorkflowProcessReports, config)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(1), result.ReportsProcessed)
	s.NotZero(result.StartedAt)
	s.NotZero(result.CompletedAt)
	s.T().Logf("Successfully processed %d reports", result.ReportsProcessed)
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessReports_NoRankings() {
	s.T().Log("Testing workflow with no rankings")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
		Specs:  []models.ClassSpec{{ClassName: "Priest", SpecName: "Shadow"}},
		Dungeons: []models.Dungeon{
			{Name: "TestDungeon", EncounterID: 12660},
		},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
			return []*warcraftlogsBuilds.ClassRanking{}, nil
		},
		activity.RegisterOptions{Name: definitions.GetStoredRankingsActivity},
	)

	var result models.WorkflowResult
	s.env.ExecuteWorkflow(WorkflowProcessReports, config)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Zero(result.ReportsProcessed)
	s.T().Log("Successfully handled no rankings case")
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessReports_RateLimit() {
	s.T().Log("Testing rate limit handling")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
		Specs:  []models.ClassSpec{{ClassName: "Priest", SpecName: "Shadow"}},
		Dungeons: []models.Dungeon{
			{Name: "TestDungeon", EncounterID: 12660},
		},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
			return []*warcraftlogsBuilds.ClassRanking{
				{ReportCode: "TEST123"},
			}, nil
		},
		activity.RegisterOptions{Name: definitions.GetStoredRankingsActivity},
	)

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking, worker models.WorkerConfig) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeRateLimit,
				Message:   "Rate limit exceeded",
				RetryIn:   5 * time.Minute,
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	s.env.ExecuteWorkflow(WorkflowProcessReports, config)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Rate limit exceeded")
	s.T().Log("Successfully handled rate limit error")
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessReports_CancellationHandling() {
	s.T().Log("Testing workflow cancellation handling")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
		Specs:  []models.ClassSpec{{ClassName: "Priest", SpecName: "Shadow"}},
		Dungeons: []models.Dungeon{
			{Name: "TestDungeon", EncounterID: 12660},
		},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
			time.Sleep(time.Millisecond * 100)
			return []*warcraftlogsBuilds.ClassRanking{{ReportCode: "TEST123"}}, nil
		},
		activity.RegisterOptions{Name: definitions.GetStoredRankingsActivity},
	)

	s.env.RegisterDelayedCallback(func() {
		s.env.CancelWorkflow()
	}, time.Millisecond*50)

	s.env.ExecuteWorkflow(WorkflowProcessReports, config)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.True(temporal.IsCanceledError(err))
	s.T().Log("Successfully handled workflow cancellation")
}

func (s *WorkflowTestSuite) Test_Workflow_ProcessReports_ErrorRecovery() {
	s.T().Log("Testing error recovery handling")

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{
			NumWorkers:   1,
			RequestDelay: time.Millisecond,
		},
		Specs: []models.ClassSpec{{ClassName: "Priest", SpecName: "Shadow"}},
		Dungeons: []models.Dungeon{
			{Name: "TestDungeon", EncounterID: 12660},
		},
	}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, className, specName string, encounterID uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
			return []*warcraftlogsBuilds.ClassRanking{{ReportCode: "TEST123"}}, nil
		},
		activity.RegisterOptions{Name: definitions.GetStoredRankingsActivity},
	)

	callCount := 0
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking, worker models.WorkerConfig) (*models.BatchResult, error) {
			callCount++
			if callCount == 1 {
				return nil, temporal.NewApplicationError("temporary error", "TEST_ERROR")
			}
			return &models.BatchResult{ProcessedItems: 1}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	var result models.WorkflowResult
	s.env.ExecuteWorkflow(WorkflowProcessReports, config)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(1), result.ReportsProcessed)
	s.Equal(2, callCount)
	s.T().Log("Successfully recovered from temporary error")
}

// WorkflowProcessReports is a test wrapper that sets activity options and state
func WorkflowProcessReports(ctx workflow.Context, config models.WorkflowConfig) (*models.WorkflowResult, error) {
	wf := &ReportsWorkflow{
		stateManager: state.NewManager(),
		processor:    NewProcessor(),
	}

	// Set initial state from config
	state := wf.stateManager.GetState()
	if len(config.Specs) > 0 && state.CurrentSpec == nil {
		state.CurrentSpec = &config.Specs[0] // Use first spec from config
	}
	if len(config.Dungeons) > 0 && state.CurrentDungeon == nil {
		state.CurrentDungeon = &config.Dungeons[0] // Use first dungeon from config
	}
	if state.StartedAt.IsZero() {
		state.StartedAt = workflow.Now(ctx)
	}

	// Set activity options with required timeout
	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute, // Required timeout for activities
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
