// utils.go
package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"fmt"
	"strings"
	"time"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// Constants for workflow IDs and names
const (
	WorkflowIDFormat = "warcraft-logs-%s-workflow-%s"
	DateFormat       = "2006-01-02"
)

// GenerateWorkflowID generates a workflow ID for a given class
func GenerateWorkflowID(className string) string {
	return fmt.Sprintf(WorkflowIDFormat, className, time.Now().UTC().Format(DateFormat))
}

// ExtractClassFromWorkflowID extracts the class name from a workflow ID
// and ensures it has the correct capitalization according to WoW standards.
func ExtractClassFromWorkflowID(workflowID string) string {
	// Split the workflow ID by hyphen
	parts := strings.Split(workflowID, "-")

	fmt.Printf("DEBUG ExtractClassFromWorkflowID: workflowID=%s, parts=%v\n", workflowID, parts)

	// Check if we have enough parts to extract the class name
	if len(parts) >= 3 {
		// Class name is typically the 3rd part (index 2), regardless of how many parts follow
		rawClassName := parts[2]

		// Convert to lowercase for consistent mapping
		lowerClassName := strings.ToLower(rawClassName)

		// Comprehensive mapping of all class names
		classMap := map[string]string{
			"priest":      "Priest",
			"deathknight": "DeathKnight",
			"demonhunter": "DemonHunter",
			"druid":       "Druid",
			"hunter":      "Hunter",
			"mage":        "Mage",
			"monk":        "Monk",
			"paladin":     "Paladin",
			"rogue":       "Rogue",
			"shaman":      "Shaman",
			"warlock":     "Warlock",
			"warrior":     "Warrior",
			"evoker":      "Evoker",
		}

		// Try to find the class in our mapping
		if correctClass, exists := classMap[lowerClassName]; exists {
			fmt.Printf("DEBUG ExtractClassFromWorkflowID: mapped %s to %s\n", lowerClassName, correctClass)
			return correctClass
		}

		// Fallback: Capitalize first letter if no mapping found
		if len(rawClassName) > 0 {
			properCased := strings.ToUpper(rawClassName[0:1]) + strings.ToLower(rawClassName[1:])
			fmt.Printf("DEBUG ExtractClassFromWorkflowID: fallback capitalization %s to %s\n",
				rawClassName, properCased)
			return properCased
		}

		return rawClassName
	}

	fmt.Printf("DEBUG ExtractClassFromWorkflowID: insufficient parts in workflowID, returning empty string\n")
	return ""
}

// GenerateSpecKey generates a unique key for a spec
func GenerateSpecKey(spec models.ClassSpec) string {
	return fmt.Sprintf("%s-%s", spec.ClassName, spec.SpecName)
}

// GenerateDungeonKey generates a unique key for a dungeon
func GenerateDungeonKey(spec models.ClassSpec, dungeon models.Dungeon) string {
	return fmt.Sprintf("%s_%s_%d", spec.ClassName, spec.SpecName, dungeon.ID)
}

// FilterSpecsForClass filters specs for a specific class
// The comparison is case-insensitive to handle various casing formats
// Example: Filtering for "priest" will match specs with ClassName "Priest"
func FilterSpecsForClass(specs []models.ClassSpec, className string) []models.ClassSpec {
	var filteredSpecs []models.ClassSpec

	// Log input for debugging
	fmt.Printf("DEBUG FilterSpecsForClass: input className=%s, total specs=%d\n",
		className, len(specs))

	// Normalize the input class name to lowercase for comparison
	lowerClassName := strings.ToLower(className)

	// Map of lowercase class names to their proper casing
	// This allows us to standardize class names for comparison
	classMap := map[string]string{
		"priest":      "Priest",
		"deathknight": "DeathKnight",
		"demonhunter": "DemonHunter",
		"druid":       "Druid",
		"hunter":      "Hunter",
		"mage":        "Mage",
		"monk":        "Monk",
		"paladin":     "Paladin",
		"rogue":       "Rogue",
		"shaman":      "Shaman",
		"warlock":     "Warlock",
		"warrior":     "Warrior",
		"evoker":      "Evoker",
	}

	// Get the proper casing for the class name if it exists in our map
	properClassName, exists := classMap[lowerClassName]
	if !exists {
		properClassName = className // Use original if not in map
	}

	fmt.Printf("DEBUG FilterSpecsForClass: normalized className to %s\n", properClassName)

	// Filter specs based on the proper class name
	for i, spec := range specs {
		// For debugging: log each spec we're checking
		fmt.Printf("DEBUG FilterSpecsForClass: checking spec[%d] with class=%s\n",
			i, spec.ClassName)

		// Case-insensitive comparison
		if strings.EqualFold(spec.ClassName, properClassName) {
			filteredSpecs = append(filteredSpecs, spec)
			fmt.Printf("DEBUG FilterSpecsForClass: matched spec %s-%s\n",
				spec.ClassName, spec.SpecName)
		}
	}

	fmt.Printf("DEBUG FilterSpecsForClass: found %d matching specs\n", len(filteredSpecs))
	return filteredSpecs
}
