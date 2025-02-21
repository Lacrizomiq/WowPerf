package warcraftlogsBuildsTemporalWorkflowsReports

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// Test_ProcessReports_Success verifies successful processing of reports
func (s *UnitTestSuite) Test_ProcessReports_Success() {
	s.T().Log("Testing successful reports processing")

	// Test data setup
	rankings := []*warcraftlogsBuilds.ClassRanking{
		{
			PlayerName:  "TestPlayer1",
			Class:       "Priest",
			Spec:        "Shadow",
			ReportCode:  "ABC123",
			EncounterID: uint(12660),
		},
		{
			PlayerName:  "TestPlayer2",
			Class:       "Priest",
			Spec:        "Shadow",
			ReportCode:  "ABC124",
			EncounterID: uint(12660),
		},
	}

	workerConfig := models.WorkerConfig{
		NumWorkers:   2,
		RequestDelay: time.Millisecond * 500,
	}

	// Mock activity implementation
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: 2,
				ClassName:      "Priest",
				SpecName:       "Shadow",
				EncounterID:    12660,
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	// Execute workflow
	var result models.BatchResult
	s.env.ExecuteWorkflow(ProcessReportsWorkflow, rankings, workerConfig)

	// Verify results
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(2), result.ProcessedItems)
	s.Equal("Priest", result.ClassName)
	s.Equal("Shadow", result.SpecName)

	s.T().Logf("Successfully processed %d reports", result.ProcessedItems)
}

// / Test_ProcessReports_EmptyRankings verifies handling of empty rankings list
func (s *UnitTestSuite) Test_ProcessReports_EmptyRankings() {
	s.T().Log("Testing empty rankings handling")

	rankings := []*warcraftlogsBuilds.ClassRanking{}
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	// Register the activity even tho its empty
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: 0,
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	// Execute workflow
	var result models.BatchResult
	s.env.ExecuteWorkflow(ProcessReportsWorkflow, rankings, workerConfig)

	// Verify results
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(0), result.ProcessedItems)

	s.T().Log("Successfully handled empty rankings")
}

// Test_ProcessReports_ActivityError verifies error handling during report processing
func (s *UnitTestSuite) Test_ProcessReports_ActivityError() {
	s.T().Log("Testing activity error handling")

	rankings := []*warcraftlogsBuilds.ClassRanking{
		{
			PlayerName:  "TestPlayer",
			Class:       "Priest",
			Spec:        "Shadow",
			ReportCode:  "ABC123",
			EncounterID: uint(12660),
		},
	}
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	// Mock activity with error
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeAPI,
				Message:   "Failed to process reports",
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	// Execute workflow
	s.env.ExecuteWorkflow(ProcessReportsWorkflow, rankings, workerConfig)

	// Verify error handling
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Failed to process reports")

	s.T().Log("Successfully handled activity error")
}

// ProcessReportsWorkflow helper function to run the processor
func ProcessReportsWorkflow(ctx workflow.Context, rankings []*warcraftlogsBuilds.ClassRanking, workerConfig models.WorkerConfig) (*models.BatchResult, error) {
	processor := NewProcessor()
	return processor.ProcessReports(ctx, rankings, workerConfig)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
