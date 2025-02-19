package warcraftlogsBuildsTemporalWorkflowsRankings

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"

	common "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/common"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// UnitTestSuite tests the rankings processor functionality
type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

// SetupTest prepares the test environment before each test
func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

// Test_ProcessRankings_Success verifies successful processing of rankings
func (s *UnitTestSuite) Test_ProcessRankings_Success() {
	s.T().Log("Testing successful rankings processing")

	// Test data setup
	spec := models.ClassSpec{
		ClassName: "Priest",
		SpecName:  "Shadow",
	}
	dungeon := models.Dungeon{
		ID:   1,
		Name: "Test Dungeon",
	}
	batchConfig := models.BatchConfig{
		Size:        100,
		RetryDelay:  5 * time.Second,
		MaxAttempts: 3,
	}

	// Mock activity implementation
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, s models.ClassSpec, d models.Dungeon, b models.BatchConfig) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: 42,
				ClassName:      s.ClassName,
				SpecName:       s.SpecName,
				EncounterID:    d.ID,
			}, nil
		},
		activity.RegisterOptions{Name: definitions.FetchRankingsActivity},
	)

	// Execute workflow
	var result models.BatchResult
	s.env.ExecuteWorkflow(ProcessRankingsWorkflow, spec, dungeon, batchConfig)

	// Verify results
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(42), result.ProcessedItems)
	s.Equal("Priest", result.ClassName)
	s.Equal("Shadow", result.SpecName)

	s.T().Logf("Successfully processed %d items", result.ProcessedItems)
}

// Test_ProcessRankings_RateLimit verifies rate limit error handling
func (s *UnitTestSuite) Test_ProcessRankings_RateLimit() {
	s.T().Log("Testing rate limit error handling")

	// Test data setup
	spec := models.ClassSpec{
		ClassName: "Priest",
		SpecName:  "Shadow",
	}
	dungeon := models.Dungeon{
		ID:   1,
		Name: "Test Dungeon",
	}
	batchConfig := models.BatchConfig{
		Size:        100,
		RetryDelay:  0, // No retry delay for test
		MaxAttempts: 1, // Single attempt for test
	}

	// Mock activity to return rate limit error
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, s models.ClassSpec, d models.Dungeon, b models.BatchConfig) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeRateLimit,
				Message:   "Rate limit exceeded",
				RetryIn:   5 * time.Minute,
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.FetchRankingsActivity},
	)

	// Execute workflow
	s.env.ExecuteWorkflow(ProcessRankingsWorkflow, spec, dungeon, batchConfig)

	// Verify error handling
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Rate limit exceeded")

	s.T().Log("Successfully handled rate limit error")
}

// Test_ProcessRankings_InvalidSpec verifies invalid spec handling
func (s *UnitTestSuite) Test_ProcessRankings_InvalidSpec() {
	s.T().Log("Testing invalid spec handling")

	spec := models.ClassSpec{
		ClassName: "", // Invalid: empty class name
		SpecName:  "Shadow",
	}
	dungeon := models.Dungeon{
		ID:   1,
		Name: "Test Dungeon",
	}
	batchConfig := models.BatchConfig{
		Size: 100,
	}

	// Execute workflow
	s.env.ExecuteWorkflow(ProcessRankingsWorkflow, spec, dungeon, batchConfig)

	// Verify error handling
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	// Verify error message and type
	s.Contains(err.Error(), "class name cannot be empty")
	s.Contains(err.Error(), "configuration error")

	// For debug
	s.T().Logf("Actual error: %v", err)
	s.T().Log("Successfully handled invalid spec error")
}

// ProcessRankingsWorkflow helper function to run the processor
func ProcessRankingsWorkflow(ctx workflow.Context, spec models.ClassSpec, dungeon models.Dungeon, batchConfig models.BatchConfig) (*models.BatchResult, error) {
	processor := NewProcessor()
	return processor.ProcessRankings(ctx, spec, dungeon, batchConfig)
}

// TestUnitTestSuite runs the test suite
func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
