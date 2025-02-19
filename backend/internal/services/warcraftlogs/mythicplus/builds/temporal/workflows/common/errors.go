// errors.go
package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"fmt"
	"time"
)

// ErrorType represents different types of workflow errors
type ErrorType string

const (
	ErrorTypeRateLimit     ErrorType = "rate_limit"
	ErrorTypeConfiguration ErrorType = "configuration"
	ErrorTypeAPI           ErrorType = "api"
	ErrorTypeDatabase      ErrorType = "database"
)

// WorkflowError represents a custom error with additional context
type WorkflowError struct {
	Type      ErrorType
	Message   string
	Retryable bool
	RetryIn   time.Duration
	Source    string
}

func (e *WorkflowError) Error() string {
	if e.RetryIn > 0 {
		return fmt.Sprintf("%s error: %s (retry in %v)", e.Type, e.Message, e.RetryIn)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string, retryIn time.Duration) *WorkflowError {
	return &WorkflowError{
		Type:      ErrorTypeRateLimit,
		Message:   message,
		Retryable: true,
		RetryIn:   retryIn,
		Source:    "warcraftlogs_api",
	}
}

// IsRateLimitError checks if an error is a rate limit error
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	if wfErr, ok := err.(*WorkflowError); ok {
		return wfErr.Type == ErrorTypeRateLimit
	}
	return false
}
