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

// WorkflowTestSuite defines a test suite for validating the rankings workflow.
// It uses testify's Suite for assertions and Temporal's WorkflowTestSuite for testing workflows.
type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

// SetupTest prepares the test environment before each workflow test case.
// It initializes a fresh Temporal test environment and registers the workflow.
func (s *WorkflowTestSuite) SetupTest() {
	// Create a new test workflow environment for isolated testing.
	s.env = s.NewTestWorkflowEnvironment()
	// Register WorkflowProcessRankings to ensure the test environment can execute it.
	// This prevents errors like "no specs found for class" due to an unregistered workflow.
	s.env.RegisterWorkflow(WorkflowProcessRankings)
}

// Test_Workflow_ProcessRankings_Success verifies successful execution of the rankings workflow.
// It mocks a successful activity and checks that the workflow completes with expected results.
func (s *WorkflowTestSuite) Test_Workflow_ProcessRankings_Success() {
	// Log the test start for debugging and traceability.
	s.T().Log("Testing successful rankings processing")

	// Define test inputs: a valid class spec, dungeon, and batch configuration.
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

	// Register a mock FetchRankingsActivity that simulates a successful response.
	// Returns a BatchResult with 42 processed items matching the input data.
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

	// Execute the workflow and capture the result.
	var result models.BatchResult
	s.env.ExecuteWorkflow(WorkflowProcessRankings, spec, dungeon, batchConfig)

	// Verify that the workflow completed successfully and returned the expected output.
	s.True(s.env.IsWorkflowCompleted(), "Workflow should complete successfully")
	s.NoError(s.env.GetWorkflowResult(&result), "Workflow should return no error")
	s.Equal(int32(42), result.ProcessedItems, "Processed items should match the mocked value")
	s.Equal("Priest", result.ClassName, "Class name should match the input")
	s.Equal("Shadow", result.SpecName, "Spec name should match the input")

	// Log the successful outcome with the number of processed items.
	s.T().Logf("Successfully processed %d items", result.ProcessedItems)
}

// Test_Workflow_ProcessRankings_RateLimit tests the workflow’s handling of a rate limit error.
// It mocks an activity that fails with a rate limit error and verifies proper error handling.
func (s *WorkflowTestSuite) Test_Workflow_ProcessRankings_RateLimit() {
	// Log the test start for debugging purposes.
	s.T().Log("Testing rate limit error handling")

	// Define test inputs: a valid class spec, dungeon, and batch config with no retries.
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
		RetryDelay:  0, // No delay to ensure immediate failure
		MaxAttempts: 1, // Single attempt to avoid retries
	}

	// Register a mock FetchRankingsActivity that simulates a rate limit error.
	// Returns a WorkflowError with "Rate limit exceeded" message.
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, s models.ClassSpec, d models.Dungeon, b models.BatchConfig) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeRateLimit,
				Message:   "Rate limit exceeded",
				RetryIn:   5 * time.Minute,
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: "FetchAndStore"},
	)

	// Execute the workflow with the test inputs.
	s.env.ExecuteWorkflow(WorkflowProcessRankings, spec, dungeon, batchConfig)

	// Verify that the workflow completed with the expected rate limit error.
	s.True(s.env.IsWorkflowCompleted(), "Workflow should complete (with error)")
	err := s.env.GetWorkflowError()
	s.Error(err, "Workflow should return an error")
	s.Contains(err.Error(), "Rate limit exceeded", "Error should indicate rate limit")

	// Log confirmation of successful rate limit handling.
	s.T().Log("Successfully handled rate limit error")
}

// Test_Workflow_ProcessRankings_InvalidSpec tests the workflow with an invalid spec.
// It ensures the workflow fails appropriately when given an empty class name.
func (s *WorkflowTestSuite) Test_Workflow_ProcessRankings_InvalidSpec() {
	// Log the test start for traceability.
	s.T().Log("Testing invalid spec handling")

	// Define test inputs with an invalid spec (empty class name).
	spec := models.ClassSpec{
		ClassName: "", // Invalid: empty class name triggers validation error
		SpecName:  "Shadow",
	}
	dungeon := models.Dungeon{
		ID:   1,
		Name: "Test Dungeon",
	}
	batchConfig := models.BatchConfig{
		Size: 100,
	}

	// Execute the workflow with the invalid inputs.
	s.env.ExecuteWorkflow(WorkflowProcessRankings, spec, dungeon, batchConfig)

	// Verify that the workflow completed with an error and includes expected messages.
	s.True(s.env.IsWorkflowCompleted(), "Workflow should complete (with error)")
	err := s.env.GetWorkflowError()
	s.Error(err, "Workflow should return an error")
	s.Contains(err.Error(), "class name cannot be empty", "Error should indicate invalid class name")
	s.Contains(err.Error(), "configuration error", "Error should classify as configuration issue")

	// Log the actual error and success message for debugging and confirmation.
	s.T().Logf("Actual error: %v", err)
	s.T().Log("Successfully handled invalid spec error")
}

// WorkflowProcessRankings is the workflow function that orchestrates rankings processing.
// It creates a processor instance and delegates the processing task.
func WorkflowProcessRankings(ctx workflow.Context, spec models.ClassSpec, dungeon models.Dungeon, batchConfig models.BatchConfig) (*models.BatchResult, error) {
	// Create a new processor instance to handle the rankings logic.
	processor := NewProcessor()
	// Execute the processing and return the result or any error.
	return processor.ProcessRankings(ctx, spec, dungeon, batchConfig)
}

// TestWorkflowTestSuite is the entry point for running the workflow test suite.
// It uses testify’s suite.Run to execute all workflow-specific test cases.
func TestWorkflowTestSuite(t *testing.T) {
	// Run the test suite with a new instance of WorkflowTestSuite.
	suite.Run(t, new(WorkflowTestSuite))
}
