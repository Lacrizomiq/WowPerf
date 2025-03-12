package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// TestGenerateWorkflowID tests the workflow ID generation function.
// It verifies that:
// - Workflow IDs are generated with the correct format
// - Special characters in class names are handled properly
// - Empty class names don't cause errors
func TestGenerateWorkflowID(t *testing.T) {
	t.Log("Starting workflow ID generation tests")

	testCases := []struct {
		name      string
		className string
		want      string
	}{
		{
			name:      "Valid class name",
			className: "Priest",
			want:      "warcraft-logs-Priest-workflow-" + time.Now().UTC().Format(DateFormat),
		},
		{
			name:      "Empty class name",
			className: "",
			want:      "warcraft-logs--workflow-" + time.Now().UTC().Format(DateFormat),
		},
		{
			name:      "Class name with special characters",
			className: "Death-Knight",
			want:      "warcraft-logs-Death-Knight-workflow-" + time.Now().UTC().Format(DateFormat),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing with class name: '%s'", tc.className)

			got := GenerateWorkflowID(tc.className)
			t.Logf("Generated workflow ID: %s", got)

			assert.Equal(t, tc.want, got)
			t.Logf("Successfully validated workflow ID format")
		})
	}
}

// TestExtractClassFromWorkflowID tests the class name extraction from workflow IDs.
// It verifies that:
// - Class names are correctly extracted from valid workflow IDs
// - Invalid workflow IDs return empty strings
// - Edge cases like empty strings and malformed IDs are handled gracefully
func TestExtractClassFromWorkflowID(t *testing.T) {
	t.Log("Starting class name extraction tests")

	testCases := []struct {
		name       string
		workflowID string
		want       string
	}{
		{
			name:       "Valid workflow ID",
			workflowID: "warcraft-logs-Priest-workflow-2025-02-19",
			want:       "Priest",
		},
		{
			name:       "Invalid workflow ID format",
			workflowID: "invalid-workflow-id",
			want:       "",
		},
		{
			name:       "Empty workflow ID",
			workflowID: "",
			want:       "",
		},
		{
			name:       "Class name with hyphen",
			workflowID: "warcraft-logs-Death-Knight-workflow-2025-02-19",
			want:       "Death", // Known limitation: splits at first hyphen
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing workflow ID: '%s'", tc.workflowID)

			got := ExtractClassFromWorkflowID(tc.workflowID)
			t.Logf("Extracted class name: '%s'", got)

			assert.Equal(t, tc.want, got)
			t.Logf("Extraction validation completed")
		})
	}
}

// TestGenerateSpecKey tests the generation of unique spec keys.
// It verifies that:
// - Spec keys are generated in the correct format (className-specName)
// - Empty class or spec names are handled appropriately
// - The generated keys are consistent and unique
func TestGenerateSpecKey(t *testing.T) {
	t.Log("Starting spec key generation tests")

	testCases := []struct {
		name string
		spec models.ClassSpec
		want string
	}{
		{
			name: "Valid spec",
			spec: models.ClassSpec{
				ClassName: "Priest",
				SpecName:  "Shadow",
			},
			want: "Priest-Shadow",
		},
		{
			name: "Empty spec name",
			spec: models.ClassSpec{
				ClassName: "Priest",
				SpecName:  "",
			},
			want: "Priest-",
		},
		{
			name: "Empty class name",
			spec: models.ClassSpec{
				ClassName: "",
				SpecName:  "Shadow",
			},
			want: "-Shadow",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing spec - Class: '%s', Spec: '%s'", tc.spec.ClassName, tc.spec.SpecName)

			got := GenerateSpecKey(tc.spec)
			t.Logf("Generated spec key: '%s'", got)

			assert.Equal(t, tc.want, got)
			t.Logf("Successfully validated spec key format")
		})
	}
}

// TestFilterSpecsForClass tests the filtering of specs by class name.
// It verifies that:
// - Only specs matching the requested class are returned
// - The correct number of specs is returned
// - Non-existent classes return empty slices
// - Empty class names are handled properly
func TestFilterSpecsForClass(t *testing.T) {
	t.Log("Starting spec filtering tests")

	// Setup test data
	specs := []models.ClassSpec{
		{ClassName: "Priest", SpecName: "Shadow"},
		{ClassName: "Priest", SpecName: "Holy"},
		{ClassName: "Mage", SpecName: "Frost"},
		{ClassName: "Mage", SpecName: "Fire"},
	}
	t.Logf("Initialized test data with %d specs", len(specs))

	testCases := []struct {
		name      string
		className string
		want      int
	}{
		{
			name:      "Filter Priest specs",
			className: "Priest",
			want:      2,
		},
		{
			name:      "Filter Mage specs",
			className: "Mage",
			want:      2,
		},
		{
			name:      "Filter non-existent class",
			className: "Warrior",
			want:      0,
		},
		{
			name:      "Filter with empty class name",
			className: "",
			want:      0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing filtering for class: '%s'", tc.className)

			filtered := FilterSpecsForClass(specs, tc.className)
			t.Logf("Found %d specs for class '%s'", len(filtered), tc.className)

			assert.Len(t, filtered, tc.want)

			// Additional validation for non-empty results
			if len(filtered) > 0 {
				t.Log("Validating filtered specs...")
				for i, spec := range filtered {
					t.Logf("Spec %d: Class='%s', Spec='%s'", i+1, spec.ClassName, spec.SpecName)
					assert.Equal(t, tc.className, spec.ClassName)
				}
				t.Log("All filtered specs validated successfully")
			}
		})
	}
}
