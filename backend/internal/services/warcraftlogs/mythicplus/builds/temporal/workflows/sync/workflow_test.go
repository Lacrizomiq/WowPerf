package warcraftlogsBuildsTemporalWorkflowsSync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
	state "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/state"
)

type SyncWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env       *testsuite.TestWorkflowEnvironment
	rateLimit bool // Flag to control rate limit behavior
}

func (s *SyncWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.env.RegisterWorkflow(ExecuteTestSyncWorkflow)

	// Mock Rankings Workflow with fully qualified path
	s.env.RegisterWorkflowWithOptions(
		func(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
			if s.rateLimit {
				return nil, temporal.NewApplicationError("API rate limit exceeded", "RateLimitError")
			}
			return &models.WorkflowResult{
				RankingsProcessed: 20,
				StartedAt:         workflow.Now(ctx),
				CompletedAt:       workflow.Now(ctx),
			}, nil
		},
		workflow.RegisterOptions{
			Name: "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/rankings.RankingsWorkflow",
		},
	)

	// Mock ProcessBuildBatch Workflow (for both Reports and Builds)
	s.env.RegisterWorkflowWithOptions(
		func(ctx workflow.Context, params models.WorkflowConfig) (*models.WorkflowResult, error) {
			return &models.WorkflowResult{
				ReportsProcessed: 10,
				BuildsProcessed:  5,
				StartedAt:        workflow.Now(ctx),
				CompletedAt:      workflow.Now(ctx),
			}, nil
		},
		workflow.RegisterOptions{Name: definitions.ProcessBuildBatchWorkflowName},
	)
}

func (s *SyncWorkflowTestSuite) Test_SyncWorkflow_FullExecution() {
	s.T().Log("Testing full sync workflow execution")
	s.rateLimit = false

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{
			NumWorkers:   2,
			RequestDelay: time.Millisecond * 100,
		},
		Rankings: models.RankingsConfig{
			MaxRankingsPerSpec: 150,
			Batch: models.BatchConfig{
				Size:        20,
				RetryDelay:  time.Second,
				MaxAttempts: 3,
			},
		},
		Specs: []models.ClassSpec{
			{ClassName: "Priest", SpecName: "Shadow"},
			{ClassName: "Priest", SpecName: "Holy"},
		},
		Dungeons: []models.Dungeon{
			{ID: 12660, Name: "Test Dungeon", EncounterID: 12660},
		},
	}

	var result models.WorkflowResult
	s.env.ExecuteWorkflow(ExecuteTestSyncWorkflow, config)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
	s.NoError(s.env.GetWorkflowResult(&result))

	s.Greater(result.RankingsProcessed, int32(0))
	s.Greater(result.ReportsProcessed, int32(0))
	s.Greater(result.BuildsProcessed, int32(0))

	s.NotZero(result.StartedAt)
	s.NotZero(result.CompletedAt)
	s.True(result.CompletedAt.After(result.StartedAt))

	s.T().Logf("Successfully completed sync workflow - Rankings: %d, Reports: %d, Builds: %d",
		result.RankingsProcessed, result.ReportsProcessed, result.BuildsProcessed)
}

func (s *SyncWorkflowTestSuite) Test_SyncWorkflow_RateLimit() {
	s.T().Log("Testing rate limit handling")
	s.rateLimit = true

	config := models.WorkflowConfig{
		Worker: models.WorkerConfig{NumWorkers: 1},
		Specs: []models.ClassSpec{
			{ClassName: "Priest", SpecName: "Shadow"},
		},
		Dungeons: []models.Dungeon{
			{ID: 12660, Name: "Test Dungeon", EncounterID: 12660},
		},
	}

	s.env.ExecuteWorkflow(ExecuteTestSyncWorkflow, config)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.True(workflow.IsContinueAsNewError(err))
	s.Contains(err.Error(), "API rate limit exceeded")

	s.T().Log("Successfully handled rate limit error")
}

func ExecuteTestSyncWorkflow(ctx workflow.Context, config models.WorkflowConfig) (*models.WorkflowResult, error) {
	wf := &SyncWorkflow{
		stateManager: state.NewManager(),
		orchestrator: NewOrchestrator(),
	}

	workflowState := wf.stateManager.GetState()
	if len(config.Specs) > 0 {
		workflowState.CurrentSpec = &config.Specs[0]
	}
	if len(config.Dungeons) > 0 {
		workflowState.CurrentDungeon = &config.Dungeons[0]
	}
	workflowState.StartedAt = workflow.Now(ctx)

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

func TestSyncWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(SyncWorkflowTestSuite))
}
