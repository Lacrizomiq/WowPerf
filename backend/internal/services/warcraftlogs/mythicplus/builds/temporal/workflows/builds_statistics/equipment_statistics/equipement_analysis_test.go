package warcraftlogsBuildsTemporalWorkflowsBuildsStatisticsEquipmentStatistics_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	equipmentWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds_statistics/equipment_statistics"
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// TestEquipmentAnalysisWorkflow_HappyPath tests the normal flow when everything succeeds
func TestEquipmentAnalysisWorkflow_HappyPath(t *testing.T) {
	// Setup
	testEnv, mockParams := setupTestEnvironment()

	// Configure activity mock behaviors
	configureWorkflowStateMocks(testEnv)

	// Configure ProcessItemStatistics to return success
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything, // context
		mockParams.Spec[0].ClassName,
		mockParams.Spec[0].SpecName,
		uint(mockParams.Dungeon[0].EncounterID),
		int(mockParams.BatchSize),
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   100,
		ItemsAnalyzed: 50,
		CompletedAt:   time.Now(),
	}, nil)

	// Execute workflow
	testEnv.ExecuteWorkflow(equipmentWorkflow.NewEquipmentAnalysisWorkflow().Execute, *mockParams)

	// Verify
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result models.EquipmentAnalysisWorkflowResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))

	// Verify result
	assert.Equal(t, int32(100), result.TotalBuilds)
	assert.Equal(t, int32(50), result.ItemsAnalyzed)
	assert.Equal(t, int32(1), result.SpecsProcessed)
	assert.Equal(t, int32(1), result.DungeonsProcessed)
	assert.False(t, result.CompletedAt.IsZero())

	// Verify all expected mocks were called
	testEnv.AssertExpectations(t)
}

// TestEquipmentAnalysisWorkflow_MultipleSpecsAndDungeons tests processing multiple specs and dungeons
func TestEquipmentAnalysisWorkflow_MultipleSpecsAndDungeons(t *testing.T) {
	// Setup
	testEnv, mockParams := setupTestEnvironment()

	// Add more specs and dungeons
	mockParams.Spec = append(mockParams.Spec, models.ClassSpec{
		ClassName: "Mage",
		SpecName:  "Frost",
	})
	mockParams.Dungeon = append(mockParams.Dungeon, models.Dungeon{
		Name:        "Throne of Tides",
		EncounterID: 1002,
	})

	// Configure workflow state mocks
	configureWorkflowStateMocks(testEnv)

	// Configure ProcessItemStatistics for different combinations
	// Warrior/Fury/Black Rook Hold
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		"Warrior", "Fury", uint(1001), mock.Anything,
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   50,
		ItemsAnalyzed: 25,
	}, nil)

	// Warrior/Fury/Throne of Tides
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		"Warrior", "Fury", uint(1002), mock.Anything,
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   60,
		ItemsAnalyzed: 30,
	}, nil)

	// Mage/Frost/Black Rook Hold
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		"Mage", "Frost", uint(1001), mock.Anything,
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   70,
		ItemsAnalyzed: 35,
	}, nil)

	// Mage/Frost/Throne of Tides
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		"Mage", "Frost", uint(1002), mock.Anything,
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   80,
		ItemsAnalyzed: 40,
	}, nil)

	// Execute workflow
	testEnv.ExecuteWorkflow(equipmentWorkflow.NewEquipmentAnalysisWorkflow().Execute, *mockParams)

	// Verify
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result models.EquipmentAnalysisWorkflowResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))

	// Verify aggregated results
	assert.Equal(t, int32(260), result.TotalBuilds)   // 50+60+70+80
	assert.Equal(t, int32(130), result.ItemsAnalyzed) // 25+30+35+40
	assert.Equal(t, int32(2), result.SpecsProcessed)
	assert.Equal(t, int32(2), result.DungeonsProcessed)

	testEnv.AssertExpectations(t)
}

// TestEquipmentAnalysisWorkflow_ActivityError tests handling of activity errors
func TestEquipmentAnalysisWorkflow_ActivityError(t *testing.T) {
	// Setup
	testEnv, mockParams := setupTestEnvironment()

	// Configure workflow state mocks
	configureWorkflowStateMocks(testEnv)

	// Make the activity fail for the first combination
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		"Warrior", "Fury", uint(1001), mock.Anything,
	).Return(nil, errors.New("processing failed"))

	// Add a second dungeon that will succeed
	mockParams.Dungeon = append(mockParams.Dungeon, models.Dungeon{
		Name:        "Throne of Tides",
		EncounterID: 1002,
	})

	// The second combination should succeed
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		"Warrior", "Fury", uint(1002), mock.Anything,
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   100,
		ItemsAnalyzed: 50,
	}, nil)

	// Execute workflow
	testEnv.ExecuteWorkflow(equipmentWorkflow.NewEquipmentAnalysisWorkflow().Execute, *mockParams)

	// Verify
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError()) // Workflow should complete despite activity error

	var result models.EquipmentAnalysisWorkflowResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))

	// Should only contain results from the successful activity
	assert.Equal(t, int32(100), result.TotalBuilds)
	assert.Equal(t, int32(50), result.ItemsAnalyzed)
	assert.Equal(t, int32(1), result.SpecsProcessed)
	assert.Equal(t, int32(1), result.DungeonsProcessed)

	testEnv.AssertExpectations(t)
}

// TestEquipmentAnalysisWorkflow_EmptyParams tests handling of empty parameters
func TestEquipmentAnalysisWorkflow_EmptyParams(t *testing.T) {
	// Setup
	testEnv := testsuite.WorkflowTestSuite{}
	env := testEnv.NewTestWorkflowEnvironment()

	// Empty specs
	emptySpecsParams := &models.EquipmentAnalysisWorkflowParams{
		Spec:          []models.ClassSpec{},
		Dungeon:       []models.Dungeon{{Name: "Black Rook Hold", EncounterID: 1001}},
		BatchSize:     50,
		Concurrency:   4,
		RetryAttempts: 3,
		RetryDelay:    time.Second * 5,
		BatchID:       "test-batch",
	}

	// Execute workflow with empty specs
	env.ExecuteWorkflow(equipmentWorkflow.NewEquipmentAnalysisWorkflow().Execute, *emptySpecsParams)

	// Verify
	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError())

	// Verify specific error
	err := env.GetWorkflowError()
	assert.Contains(t, err.Error(), "no specs found")

	// Reset
	env = testEnv.NewTestWorkflowEnvironment()

	// Empty dungeons
	emptyDungeonsParams := &models.EquipmentAnalysisWorkflowParams{
		Spec:          []models.ClassSpec{{ClassName: "Warrior", SpecName: "Fury"}},
		Dungeon:       []models.Dungeon{},
		BatchSize:     50,
		Concurrency:   4,
		RetryAttempts: 3,
		RetryDelay:    time.Second * 5,
		BatchID:       "test-batch",
	}

	// Execute workflow with empty dungeons
	env.ExecuteWorkflow(equipmentWorkflow.NewEquipmentAnalysisWorkflow().Execute, *emptyDungeonsParams)

	// Verify
	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError())

	// Verify specific error
	err = env.GetWorkflowError()
	assert.Contains(t, err.Error(), "no dungeons found")
}

// TestEquipmentAnalysisWorkflow_WorkflowState tests the workflow state management
func TestEquipmentAnalysisWorkflow_WorkflowState(t *testing.T) {
	// Setup
	testEnv, mockParams := setupTestEnvironment()

	// We'll track the workflow state updates to verify progression
	var lastState *warcraftlogsBuilds.WorkflowState

	// Create workflow state - should happen once
	testEnv.OnActivity(
		definitions.CreateWorkflowStateActivity,
		mock.Anything,
		mock.MatchedBy(func(state *warcraftlogsBuilds.WorkflowState) bool {
			return state.WorkflowType == "equipment-analysis" && state.Status == "running"
		}),
	).Return(&warcraftlogsBuilds.WorkflowState{}, nil)

	// Process activity state updates - will happen multiple times
	testEnv.OnActivity(
		definitions.UpdateWorkflowStateActivity,
		mock.Anything,
		mock.MatchedBy(func(state *warcraftlogsBuilds.WorkflowState) bool {
			lastState = state
			return true // Accept any state update
		}),
	).Return(nil).Times(1) // Initial + processing + final

	// Configure ProcessItemStatistics to return success
	testEnv.OnActivity(
		definitions.ProcessBuildStatisticsActivity,
		mock.Anything,
		mockParams.Spec[0].ClassName,
		mockParams.Spec[0].SpecName,
		uint(mockParams.Dungeon[0].EncounterID),
		int(mockParams.BatchSize),
	).Return(&models.EquipmentAnalysisWorkflowResult{
		TotalBuilds:   100,
		ItemsAnalyzed: 50,
	}, nil)

	// Execute workflow
	testEnv.ExecuteWorkflow(equipmentWorkflow.NewEquipmentAnalysisWorkflow().Execute, *mockParams)

	// Verify
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	// Verify final state
	assert.NotNil(t, lastState)
	assert.Equal(t, "completed", lastState.Status)
	assert.False(t, lastState.CompletedAt.IsZero())
	assert.Equal(t, 50, lastState.ItemsProcessed) // Should match ItemsAnalyzed

	testEnv.AssertExpectations(t)
}

// Helper functions

func setupTestEnvironment() (env *testsuite.TestWorkflowEnvironment, params *models.EquipmentAnalysisWorkflowParams) {
	testSuite := testsuite.WorkflowTestSuite{}
	env = testSuite.NewTestWorkflowEnvironment()

	// Enregistrer les activités avec les bonnes signatures
	env.RegisterActivityWithOptions(
		func(ctx context.Context, state *warcraftlogsBuilds.WorkflowState) (*warcraftlogsBuilds.WorkflowState, error) {
			return state, nil
		},
		activity.RegisterOptions{Name: definitions.CreateWorkflowStateActivity},
	)

	env.RegisterActivityWithOptions(
		func(ctx context.Context, state *warcraftlogsBuilds.WorkflowState) error {
			return nil
		},
		activity.RegisterOptions{Name: definitions.UpdateWorkflowStateActivity},
	)

	env.RegisterActivityWithOptions(
		func(ctx context.Context, className, specName string, encounterID uint, batchSize int) (*models.EquipmentAnalysisWorkflowResult, error) {
			return &models.EquipmentAnalysisWorkflowResult{}, nil
		},
		activity.RegisterOptions{Name: definitions.ProcessBuildStatisticsActivity},
	)

	// Create mock parameters
	params = &models.EquipmentAnalysisWorkflowParams{
		Spec: []models.ClassSpec{
			{
				ClassName: "Warrior",
				SpecName:  "Fury",
			},
		},
		Dungeon: []models.Dungeon{
			{
				Name:        "Black Rook Hold",
				EncounterID: 1001,
			},
		},
		BatchSize:     50,
		Concurrency:   4,
		RetryAttempts: 3,
		RetryDelay:    time.Second * 5,
		BatchID:       "test-batch",
	}

	return env, params
}

func configureWorkflowStateMocks(testEnv *testsuite.TestWorkflowEnvironment) {
	// Mock CreateWorkflowState avec la bonne signature: retourne (state *WorkflowState, error)
	testEnv.OnActivity(
		definitions.CreateWorkflowStateActivity,
		mock.Anything, // context
		mock.Anything, // state
	).Return(&warcraftlogsBuilds.WorkflowState{}, nil)

	// Mock UpdateWorkflowState avec la bonne signature: retourne error
	testEnv.OnActivity(
		definitions.UpdateWorkflowStateActivity,
		mock.Anything, // context
		mock.Anything, // state
	).Return(nil).Times(1) // Allow multiple calls
}
