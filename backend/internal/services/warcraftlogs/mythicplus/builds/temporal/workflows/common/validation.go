// validation.go
package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// ValidateConfig validates the workflow configuration
func ValidateConfig(config *models.WorkflowConfig) error {
	if config == nil {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "configuration cannot be nil",
		}
	}

	if config.Rankings.MaxRankingsPerSpec <= 0 {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "max rankings per spec must be greater than 0",
		}
	}

	if config.Worker.NumWorkers <= 0 {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "number of workers must be greater than 0",
		}
	}

	if len(config.Specs) == 0 {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "at least one spec must be configured",
		}
	}

	if len(config.Dungeons) == 0 {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "at least one dungeon must be configured",
		}
	}

	return nil
}

// ValidateSpec validates a class specialization
func ValidateSpec(spec models.ClassSpec) error {
	if spec.ClassName == "" {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "class name cannot be empty",
		}
	}

	if spec.SpecName == "" {
		return &WorkflowError{
			Type:    ErrorTypeConfiguration,
			Message: "spec name cannot be empty",
		}
	}

	return nil
}
