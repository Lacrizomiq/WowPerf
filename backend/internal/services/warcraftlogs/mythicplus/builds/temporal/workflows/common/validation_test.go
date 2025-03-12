package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// TestValidateConfig tests the configuration validation functionality.
// It verifies that:
// - Valid configurations pass validation
// - Invalid configurations return appropriate errors
// - All required fields are properly checked
func TestValidateConfig(t *testing.T) {
	t.Log("Starting configuration validation tests")

	testCases := []struct {
		name        string
		config      *models.WorkflowConfig
		expectError bool
		errorType   ErrorType
	}{
		{
			name: "Valid configuration",
			config: &models.WorkflowConfig{
				Rankings: models.RankingsConfig{
					MaxRankingsPerSpec: 100,
				},
				Worker: models.WorkerConfig{
					NumWorkers: 3,
				},
				Specs: []models.ClassSpec{
					{ClassName: "Priest", SpecName: "Shadow"},
				},
				Dungeons: []models.Dungeon{
					{ID: 1, Name: "Test Dungeon"},
				},
			},
			expectError: false,
		},
		{
			name:        "Nil configuration",
			config:      nil,
			expectError: true,
			errorType:   ErrorTypeConfiguration,
		},
		{
			name: "Invalid max rankings",
			config: &models.WorkflowConfig{
				Rankings: models.RankingsConfig{
					MaxRankingsPerSpec: 0,
				},
				Worker:   models.WorkerConfig{NumWorkers: 3},
				Specs:    []models.ClassSpec{{ClassName: "Priest"}},
				Dungeons: []models.Dungeon{{ID: 1}},
			},
			expectError: true,
			errorType:   ErrorTypeConfiguration,
		},
		{
			name: "No specs configured",
			config: &models.WorkflowConfig{
				Rankings: models.RankingsConfig{MaxRankingsPerSpec: 100},
				Worker:   models.WorkerConfig{NumWorkers: 3},
				Specs:    []models.ClassSpec{},
				Dungeons: []models.Dungeon{{ID: 1}},
			},
			expectError: true,
			errorType:   ErrorTypeConfiguration,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing configuration: %+v", tc.config)

			err := ValidateConfig(tc.config)

			if tc.expectError {
				t.Log("Expecting an error...")
				require.Error(t, err)
				if wfErr, ok := err.(*WorkflowError); ok {
					assert.Equal(t, tc.errorType, wfErr.Type)
					t.Logf("Received expected error: %v", wfErr)
				} else {
					t.Errorf("Expected WorkflowError but got: %T", err)
				}
			} else {
				t.Log("Expecting successful validation")
				require.NoError(t, err)
				t.Log("Configuration validated successfully")
			}
		})
	}
}

// TestValidateSpec tests the spec validation functionality.
// It verifies that:
// - Valid specs pass validation
// - Invalid specs return appropriate errors
// - Both class name and spec name are properly validated
func TestValidateSpec(t *testing.T) {
	t.Log("Starting spec validation tests")

	testCases := []struct {
		name        string
		spec        models.ClassSpec
		expectError bool
		errorType   ErrorType
	}{
		{
			name: "Valid spec",
			spec: models.ClassSpec{
				ClassName: "Priest",
				SpecName:  "Shadow",
			},
			expectError: false,
		},
		{
			name: "Empty class name",
			spec: models.ClassSpec{
				ClassName: "",
				SpecName:  "Shadow",
			},
			expectError: true,
			errorType:   ErrorTypeConfiguration,
		},
		{
			name: "Empty spec name",
			spec: models.ClassSpec{
				ClassName: "Priest",
				SpecName:  "",
			},
			expectError: true,
			errorType:   ErrorTypeConfiguration,
		},
		{
			name: "Both names empty",
			spec: models.ClassSpec{
				ClassName: "",
				SpecName:  "",
			},
			expectError: true,
			errorType:   ErrorTypeConfiguration,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing spec - Class: '%s', Spec: '%s'", tc.spec.ClassName, tc.spec.SpecName)

			err := ValidateSpec(tc.spec)

			if tc.expectError {
				t.Log("Expecting an error...")
				require.Error(t, err)
				if wfErr, ok := err.(*WorkflowError); ok {
					assert.Equal(t, tc.errorType, wfErr.Type)
					t.Logf("Received expected error: %v", wfErr)
				} else {
					t.Errorf("Expected WorkflowError but got: %T", err)
				}
			} else {
				t.Log("Expecting successful validation")
				require.NoError(t, err)
				t.Log("Spec validated successfully")
			}
		})
	}
}
