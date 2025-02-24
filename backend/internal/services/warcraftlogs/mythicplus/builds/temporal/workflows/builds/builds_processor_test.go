package warcraftlogsBuildsTemporalWorkflowsBuilds

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

// Test_ProcessBuilds_Success verifies successful processing of builds
func (s *UnitTestSuite) Test_ProcessBuilds_Success() {
	s.T().Log("Testing successful builds processing")

	// Test data setup with reports containing player details
	reports := []*warcraftlogsBuilds.Report{
		{
			Code:             "ABC123",
			FightID:          1,
			EncounterID:      12660,
			KeystoneLevel:    20,
			PlayerDetailsDps: []byte(`[{"name":"TestPlayer1","class":"Priest","spec":"Shadow","itemLevel":420}]`),
			TalentCodes:      []byte(`{"Priest_Shadow_talents":"test-talent-code"}`),
		},
		{
			Code:             "ABC124",
			FightID:          2,
			EncounterID:      12660,
			KeystoneLevel:    20,
			PlayerDetailsDps: []byte(`[{"name":"TestPlayer2","class":"Priest","spec":"Shadow","itemLevel":420}]`),
			TalentCodes:      []byte(`{"Priest_Shadow_talents":"test-talent-code"}`),
		},
	}

	workerConfig := models.WorkerConfig{
		NumWorkers:   2,
		RequestDelay: time.Millisecond * 500,
	}

	// Mock ProcessBuilds activity
	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.Report) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: int32(len(r)),
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildsActivity},
	)

	// Execute workflow
	var result models.BatchResult
	s.env.ExecuteWorkflow(ProcessBuildsWorkflow, reports, workerConfig)

	// Verify results
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(2), result.ProcessedItems)

	s.T().Logf("Successfully processed %d builds", result.ProcessedItems)
}

// Test_ProcessBuilds_EmptyReports verifies handling of empty reports list
func (s *UnitTestSuite) Test_ProcessBuilds_EmptyReports() {
	s.T().Log("Testing empty reports handling")

	reports := []*warcraftlogsBuilds.Report{}
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.Report) (*models.BatchResult, error) {
			return &models.BatchResult{
				ProcessedItems: 0,
				ProcessedAt:    time.Now(),
			}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildsActivity},
	)

	var result models.BatchResult
	s.env.ExecuteWorkflow(ProcessBuildsWorkflow, reports, workerConfig)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal(int32(0), result.ProcessedItems)

	s.T().Log("Successfully handled empty reports")
}

// Test_ProcessBuilds_DatabaseError verifies handling of database errors
func (s *UnitTestSuite) Test_ProcessBuilds_DatabaseError() {
	s.T().Log("Testing database error handling")

	reports := []*warcraftlogsBuilds.Report{
		{
			Code:    "ABC123",
			FightID: 1,
		},
	}
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.Report) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeDatabase,
				Message:   "Failed to store builds in database",
				Retryable: true,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildsActivity},
	)

	s.env.ExecuteWorkflow(ProcessBuildsWorkflow, reports, workerConfig)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "Failed to store builds in database")

	s.T().Log("Successfully handled database error")
}

// Test_ProcessBuilds_InvalidReportData verifies handling of invalid report data
func (s *UnitTestSuite) Test_ProcessBuilds_InvalidReportData() {
	s.T().Log("Testing invalid report data handling")

	reports := []*warcraftlogsBuilds.Report{
		{
			Code:             "ABC123",
			FightID:          1,
			PlayerDetailsDps: []byte(`[{"name":"","class":"","spec":""}]`), // Invalid: empty required fields
			TalentCodes:      []byte(`{}`),                                 // Empty talents
		},
	}
	workerConfig := models.WorkerConfig{NumWorkers: 1}

	s.env.RegisterActivityWithOptions(
		func(ctx context.Context, r []*warcraftlogsBuilds.Report) (*models.BatchResult, error) {
			return nil, &common.WorkflowError{
				Type:      common.ErrorTypeConfiguration,
				Message:   "invalid player details: missing required fields",
				Retryable: false,
			}
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildsActivity},
	)

	s.env.ExecuteWorkflow(ProcessBuildsWorkflow, reports, workerConfig)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)
	s.Contains(err.Error(), "invalid player details")

	s.T().Log("Successfully handled invalid report data")
}

// ProcessBuildsWorkflow helper function to run the processor
func ProcessBuildsWorkflow(ctx workflow.Context, reports []*warcraftlogsBuilds.Report, workerConfig models.WorkerConfig) (*models.BatchResult, error) {
	processor := NewProcessor()
	return processor.ProcessBuilds(ctx, reports, workerConfig)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
