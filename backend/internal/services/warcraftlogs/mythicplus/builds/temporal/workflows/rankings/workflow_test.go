package warcraftlogsBuildsTemporalWorkflowsRankings

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *WorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *WorkflowTestSuite) Test_Execute_Success() {
	s.T().Log("Testing successful workflow execution")

	// Prepare test data
	params := models.WorkflowConfig{
		Specs: []models.ClassSpec{
			{ClassName: "Priest", SpecName: "Shadow"},
			{ClassName: "Priest", SpecName: "Holy"},
		},
		Dungeons: []models.Dungeon{
			{ID: 1, Name: "Dungeon1"},
			{ID: 2, Name: "Dungeon2"},
		},
		Rankings: models.RankingsConfig{
			MaxRankingsPerSpec: 100,
			Batch: models.BatchConfig{
				Size:        50,
				RetryDelay:  time.Second,
				MaxAttempts: 3,
			},
		},
	}

	// Mock successful rankings processing
	s.env.OnActivity(definitions.FetchRankingsActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&models.BatchResult{
			ProcessedItems: 42,
			ClassName:      "Priest",
			SpecName:       "Shadow",
		}, nil)

	// Set test identifier
	s.env.SetWorkflowID("warcraft-logs-Priest-workflow-2025-02-19")

	workflow := NewRankingsWorkflow()
	s.env.ExecuteWorkflow(workflow.Execute, params)

	s.True(s.env.IsWorkflowCompleted())

	var result models.WorkflowResult
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(168), result.RankingsProcessed) // 42 * 4
	s.False(result.CompletedAt.IsZero())

	s.T().Log("Successfully completed workflow")
}

func (s *WorkflowTestSuite) Test_Execute_InvalidWorkflowID() {
	s.T().Log("Testing workflow with invalid ID")

	params := models.WorkflowConfig{}

	// Set invalid workflow ID
	s.env.SetWorkflowID("invalid-workflow-id")

	workflow := NewRankingsWorkflow()
	s.env.ExecuteWorkflow(workflow.Execute, params)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "invalid workflow ID format")

	s.T().Log("Successfully handled invalid workflow ID")
}

func (s *WorkflowTestSuite) Test_Execute_RateLimit() {
	s.T().Log("Testing rate limit handling")

	params := models.WorkflowConfig{
		Specs: []models.ClassSpec{
			{ClassName: "Priest", SpecName: "Shadow"},
		},
		Dungeons: []models.Dungeon{
			{ID: 1, Name: "Dungeon1"},
		},
		Rankings: models.RankingsConfig{
			Batch: models.BatchConfig{
				Size:        50,
				MaxAttempts: 1,
			},
		},
	}

	// Set workflow ID
	s.env.SetWorkflowID("warcraft-logs-Priest-workflow-2025-02-19")

	// Mock rate limit error
	rateLimitErr := &common.WorkflowError{
		Type:    common.ErrorTypeRateLimit,
		Message: "Rate limit exceeded",
		RetryIn: 5 * time.Minute,
	}

	// Register activity with rate limit error
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, spec models.ClassSpec, dungeon models.Dungeon, config models.BatchConfig) (*models.BatchResult, error) {
			return nil, rateLimitErr
		},
		activity.RegisterOptions{Name: definitions.FetchRankingsActivity},
	)

	workflow := NewRankingsWorkflow()
	s.env.ExecuteWorkflow(workflow.Execute, params)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Rate limit exceeded")

	s.T().Log("Successfully handled rate limit")
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}
