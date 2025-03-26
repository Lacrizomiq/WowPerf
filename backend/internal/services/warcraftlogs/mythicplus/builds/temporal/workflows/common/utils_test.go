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
// - Special class names like DeathKnight and DemonHunter are handled correctly
// - Mixed case input is normalized properly
func TestExtractClassFromWorkflowID(t *testing.T) {
	t.Log("Starting class name extraction tests")

	testCases := []struct {
		name       string
		workflowID string
		want       string
	}{
		// Basic test cases
		{
			name:       "Valid workflow ID - Simple class",
			workflowID: "warcraft-logs-priest-workflow-2025-02-19",
			want:       "Priest",
		},
		{
			name:       "Valid workflow ID - Test workflow",
			workflowID: "warcraft-logs-priest-test-2025-03-22",
			want:       "Priest",
		},
		// Compound class names
		{
			name:       "DeathKnight class",
			workflowID: "warcraft-logs-deathknight-workflow-2025-02-19",
			want:       "DeathKnight",
		},
		{
			name:       "DemonHunter class",
			workflowID: "warcraft-logs-demonhunter-workflow-2025-02-19",
			want:       "DemonHunter",
		},
		// Mixed case testing
		{
			name:       "Mixed case class name - PRIEST",
			workflowID: "warcraft-logs-PRIEST-workflow-2025-02-19",
			want:       "Priest",
		},
		{
			name:       "Mixed case class name - pRiEsT",
			workflowID: "warcraft-logs-pRiEsT-workflow-2025-02-19",
			want:       "Priest",
		},
		{
			name:       "Mixed case compound class - DeAtHkNiGhT",
			workflowID: "warcraft-logs-DeAtHkNiGhT-workflow-2025-02-19",
			want:       "DeathKnight",
		},
		// Edge cases
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
			name:       "Class with hyphen - not properly handled yet",
			workflowID: "warcraft-logs-Death-Knight-workflow-2025-02-19",
			want:       "Death", // Known limitation: splits at first hyphen
		},
		// All other classes
		{
			name:       "Druid class",
			workflowID: "warcraft-logs-druid-workflow-2025-02-19",
			want:       "Druid",
		},
		{
			name:       "Hunter class",
			workflowID: "warcraft-logs-hunter-workflow-2025-02-19",
			want:       "Hunter",
		},
		{
			name:       "Mage class",
			workflowID: "warcraft-logs-mage-workflow-2025-02-19",
			want:       "Mage",
		},
		{
			name:       "Monk class",
			workflowID: "warcraft-logs-monk-workflow-2025-02-19",
			want:       "Monk",
		},
		{
			name:       "Paladin class",
			workflowID: "warcraft-logs-paladin-workflow-2025-02-19",
			want:       "Paladin",
		},
		{
			name:       "Rogue class",
			workflowID: "warcraft-logs-rogue-workflow-2025-02-19",
			want:       "Rogue",
		},
		{
			name:       "Shaman class",
			workflowID: "warcraft-logs-shaman-workflow-2025-02-19",
			want:       "Shaman",
		},
		{
			name:       "Warlock class",
			workflowID: "warcraft-logs-warlock-workflow-2025-02-19",
			want:       "Warlock",
		},
		{
			name:       "Warrior class",
			workflowID: "warcraft-logs-warrior-workflow-2025-02-19",
			want:       "Warrior",
		},
		{
			name:       "Evoker class",
			workflowID: "warcraft-logs-evoker-workflow-2025-02-19",
			want:       "Evoker",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing workflow ID: '%s'", tc.workflowID)

			got := ExtractClassFromWorkflowID(tc.workflowID)
			t.Logf("Extracted class name: '%s'", got)

			assert.Equal(t, tc.want, got, "Expected extracted class to be '%s', got '%s'", tc.want, got)
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
			name: "Composite class name",
			spec: models.ClassSpec{
				ClassName: "DeathKnight",
				SpecName:  "Blood",
			},
			want: "DeathKnight-Blood",
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

			assert.Equal(t, tc.want, got, "Expected spec key to be '%s', got '%s'", tc.want, got)
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
// - Case insensitive matching works correctly
// - Special class names like DeathKnight and DemonHunter are handled correctly
func TestFilterSpecsForClass(t *testing.T) {
	t.Log("Starting spec filtering tests")

	// Setup comprehensive test data covering all classes
	specs := []models.ClassSpec{
		// Simple class names
		{ClassName: "Priest", SpecName: "Shadow"},
		{ClassName: "Priest", SpecName: "Holy"},
		{ClassName: "Priest", SpecName: "Discipline"},
		{ClassName: "Mage", SpecName: "Frost"},
		{ClassName: "Mage", SpecName: "Fire"},
		{ClassName: "Mage", SpecName: "Arcane"},
		{ClassName: "Warrior", SpecName: "Arms"},
		{ClassName: "Warrior", SpecName: "Fury"},

		// Compound class names
		{ClassName: "DeathKnight", SpecName: "Blood"},
		{ClassName: "DeathKnight", SpecName: "Frost"},
		{ClassName: "DeathKnight", SpecName: "Unholy"},
		{ClassName: "DemonHunter", SpecName: "Havoc"},
		{ClassName: "DemonHunter", SpecName: "Vengeance"},
	}
	t.Logf("Initialized test data with %d specs across multiple classes", len(specs))

	testCases := []struct {
		name      string
		className string
		want      int
		wantClass string
	}{
		// Simple class tests
		{
			name:      "Filter Priest specs - exact case",
			className: "Priest",
			want:      3,
			wantClass: "Priest",
		},
		{
			name:      "Filter Priest specs - lowercase",
			className: "priest",
			want:      3,
			wantClass: "Priest",
		},
		{
			name:      "Filter Priest specs - mixed case",
			className: "pRiEsT",
			want:      3,
			wantClass: "Priest",
		},
		{
			name:      "Filter Mage specs",
			className: "Mage",
			want:      3,
			wantClass: "Mage",
		},

		// Compound class tests
		{
			name:      "Filter DeathKnight specs - exact case",
			className: "DeathKnight",
			want:      3,
			wantClass: "DeathKnight",
		},
		{
			name:      "Filter DeathKnight specs - lowercase",
			className: "deathknight",
			want:      3,
			wantClass: "DeathKnight",
		},
		{
			name:      "Filter DeathKnight specs - mixed case",
			className: "dEaThKnIgHt",
			want:      3,
			wantClass: "DeathKnight",
		},
		{
			name:      "Filter DemonHunter specs",
			className: "demonhunter",
			want:      2,
			wantClass: "DemonHunter",
		},

		// Edge cases
		{
			name:      "Filter non-existent class",
			className: "InvalidClass",
			want:      0,
			wantClass: "",
		},
		{
			name:      "Filter with empty class name",
			className: "",
			want:      0,
			wantClass: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing filtering for class: '%s'", tc.className)

			filtered := FilterSpecsForClass(specs, tc.className)
			t.Logf("Found %d specs for class '%s'", len(filtered), tc.className)

			assert.Len(t, filtered, tc.want, "Expected %d specs, got %d", tc.want, len(filtered))

			// Additional validation for non-empty results
			if len(filtered) > 0 {
				t.Log("Validating filtered specs...")
				for i, spec := range filtered {
					t.Logf("Spec %d: Class='%s', Spec='%s'", i+1, spec.ClassName, spec.SpecName)
					assert.Equal(t, tc.wantClass, spec.ClassName,
						"Expected class name to be '%s', got '%s'", tc.wantClass, spec.ClassName)
				}
				t.Log("All filtered specs validated successfully")
			}
		})
	}
}
