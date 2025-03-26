package warcraftlogsBuildsTemporalWorkflowsCommon

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestWorkflowError tests the creation and formatting of workflow errors.
// It verifies that:
// - Error messages are formatted correctly
// - RetryIn duration is properly included when present
// - Error types are correctly assigned
func TestWorkflowError(t *testing.T) {
	t.Log("Starting workflow error tests")

	testCases := []struct {
		name     string
		err      *WorkflowError
		expected string
	}{
		{
			name: "Basic error without retry",
			err: &WorkflowError{
				Type:      ErrorTypeAPI,
				Message:   "API request failed",
				Retryable: true,
			},
			expected: "api error: API request failed",
		},
		{
			name: "Error with retry duration",
			err: &WorkflowError{
				Type:      ErrorTypeRateLimit,
				Message:   "Rate limit exceeded",
				Retryable: true,
				RetryIn:   5 * time.Minute,
			},
			expected: "rate_limit error: Rate limit exceeded (retry in 5m0s)",
		},
		{
			name: "Configuration error",
			err: &WorkflowError{
				Type:      ErrorTypeConfiguration,
				Message:   "Invalid config",
				Retryable: false,
			},
			expected: "configuration error: Invalid config",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing error with type: %s", tc.err.Type)
			t.Logf("Error message: %s", tc.err.Message)

			errString := tc.err.Error()
			t.Logf("Generated error string: %s", errString)

			assert.Equal(t, tc.expected, errString)
			t.Log("Error string format validated successfully")
		})
	}
}

// TestNewRateLimitError tests the creation of rate limit specific errors.
// It verifies that:
// - Rate limit errors are created with correct type
// - Retry information is properly set
// - Source is correctly assigned
func TestNewRateLimitError(t *testing.T) {
	t.Log("Starting rate limit error creation tests")

	testCases := []struct {
		name    string
		message string
		retryIn time.Duration
	}{
		{
			name:    "Standard rate limit",
			message: "API rate limit exceeded",
			retryIn: 1 * time.Minute,
		},
		{
			name:    "Zero retry duration",
			message: "Quota exceeded",
			retryIn: 0,
		},
		{
			name:    "Long retry duration",
			message: "Daily limit reached",
			retryIn: 24 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Creating rate limit error with message: %s", tc.message)
			t.Logf("Retry duration: %v", tc.retryIn)

			err := NewRateLimitError(tc.message, tc.retryIn)

			t.Log("Validating error properties...")
			assert.Equal(t, ErrorTypeRateLimit, err.Type)
			assert.Equal(t, tc.message, err.Message)
			assert.Equal(t, tc.retryIn, err.RetryIn)
			assert.True(t, err.Retryable)
			assert.Equal(t, "warcraftlogs_api", err.Source)

			t.Log("Rate limit error validated successfully")
		})
	}
}

// TestIsRateLimitError tests the identification of rate limit errors.
// It verifies that:
// - Rate limit errors are correctly identified
// - Other error types are not misidentified
// - Nil errors are handled properly
func TestIsRateLimitError(t *testing.T) {
	t.Log("Starting rate limit error identification tests")

	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name: "Rate limit error",
			err: &WorkflowError{
				Type:    ErrorTypeRateLimit,
				Message: "Rate limit exceeded",
			},
			expected: true,
		},
		{
			name: "Different workflow error",
			err: &WorkflowError{
				Type:    ErrorTypeAPI,
				Message: "API error",
			},
			expected: false,
		},
		{
			name:     "Standard error",
			err:      errors.New("random error"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing error: %v", tc.err)

			isRateLimit := IsRateLimitError(tc.err)
			t.Logf("Is rate limit error: %v", isRateLimit)

			assert.Equal(t, tc.expected, isRateLimit)
			t.Log("Error identification validated successfully")
		})
	}
}
