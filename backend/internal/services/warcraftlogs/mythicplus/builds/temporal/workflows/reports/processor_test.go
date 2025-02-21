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

// Test_ProcessReports_Timeout verifies that the processor handles activity timeouts correctly
func (s *UnitTestSuite) Test_ProcessReports_Timeout() {
	s.T().Log("Testing activity timeout handling")

	// Setup test data with a single ranking
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

	// Mock activity to simulate a timeout
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeTimeout,
				Message:   "Activity timed out after maximum retries",
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	// Execute workflow and verify timeout handling
	s.env.ExecuteWorkflow(ProcessReportsWorkflow, rankings, workerConfig)

	// Verify results
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Activity timed out")

	s.T().Log("Successfully verified timeout handling")
}

// Test_ProcessReports_InvalidReportCode verifies handling of invalid report codes
func (s *UnitTestSuite) Test_ProcessReports_InvalidReportCode() {
	s.T().Log("Testing invalid report code handling")

	// Setup test data with invalid report codes
	rankings := []*warcraftlogsBuilds.ClassRanking{
		{
			PlayerName:  "TestPlayer1",
			Class:       "Priest",
			Spec:        "Shadow",
			ReportCode:  "", // Invalid: empty report code
			EncounterID: uint(12660),
		},
		{
			PlayerName:  "TestPlayer2",
			Class:       "Priest",
			Spec:        "Shadow",
			ReportCode:  "ABC123", // Valid report code
			EncounterID: uint(12660),
		},
	}
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	// Mock activity to process reports with validation
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			validReports := 0
			for _, ranking := range r {
				if ranking.ReportCode != "" {
					validReports++
				}
			}
			return &models.BatchResult{
				ProcessedItems: int32(validReports),
				ClassName:      "Priest",
				SpecName:       "Shadow",
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
	s.Equal(int32(1), result.ProcessedItems) // Only valid reports should be processed

	s.T().Log("Successfully verified invalid report code handling")
}

// Test_ProcessReports_PartialSuccess verifies handling of partial successes
func (s *UnitTestSuite) Test_ProcessReports_PartialSuccess() {
	s.T().Log("Testing partial success handling")

	// Setup test data with multiple rankings
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
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	// Mock activity to simulate partial success
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			// Simulate only first report being processed successfully
			return &models.BatchResult{
				ProcessedItems: 1,
				ClassName:      "Priest",
				SpecName:       "Shadow",
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
	s.Equal(int32(1), result.ProcessedItems)

	s.T().Log("Successfully verified partial success handling")
}

// Test_ProcessReports_RateLimit verifies handling of API rate limits
func (s *UnitTestSuite) Test_ProcessReports_RateLimit() {
	s.T().Log("Testing rate limit handling")

	// Setup test data
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

	// Mock activity to simulate rate limit
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeRateLimit,
				Message:   "API rate limit exceeded",
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	// Execute workflow
	s.env.ExecuteWorkflow(ProcessReportsWorkflow, rankings, workerConfig)

	// Verify rate limit handling
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "rate limit exceeded")

	s.T().Log("Successfully verified rate limit handling")
}

// Test_ProcessReports_WorkerConfig verifies proper handling of worker configuration
func (s *UnitTestSuite) Test_ProcessReports_WorkerConfig() {
	s.T().Log("Testing worker configuration handling")

	// Setup test data
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
		RequestDelay: time.Millisecond * 100,
	}

	// Mock activity to verify worker config
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: int32(len(r)),
				ClassName:      "Priest",
				SpecName:       "Shadow",
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessReportsActivity},
	)

	// Execute workflow
	var result models.BatchResult
	s.env.ExecuteWorkflow(ProcessReportsWorkflow, rankings, workerConfig)

	// Verify results with worker configuration
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(2), result.ProcessedItems)

	s.T().Log("Successfully verified worker configuration handling")
}
