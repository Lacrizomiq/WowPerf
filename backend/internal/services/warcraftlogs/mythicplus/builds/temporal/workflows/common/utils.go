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
func ExtractClassFromWorkflowID(workflowID string) string {
	parts := strings.Split(workflowID, "-")
	if len(parts) >= 4 {
		return parts[2]
	}
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
func FilterSpecsForClass(specs []models.ClassSpec, className string) []models.ClassSpec {
	filtered := make([]models.ClassSpec, 0)
	for _, spec := range specs {
		if spec.ClassName == className {
			filtered = append(filtered, spec)
		}
	}
	return filtered
}
