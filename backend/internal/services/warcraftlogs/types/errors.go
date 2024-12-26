package warcraftlogs

import (
	"fmt"
	"time"
)

// ErrorType represents different types of errors that can occur
type ErrorType int

const (
	ErrorTypeRateLimit ErrorType = iota
	ErrorTypeAPI
	ErrorTypeNetwork
	ErrorTypeValidation
	ErrorTypeDatabase
	ErrorTypeQuotaExceeded
)

// RateLimitInfo contains detailed rate limit information
type RateLimitInfo struct {
	RemainingPoints float64
	PointsPerHour   int
	ResetIn         time.Duration
	NextRefresh     time.Time
}

// WarcraftLogsError represents a custom error with additional context
type WarcraftLogsError struct {
	Type          ErrorType
	Message       string
	Cause         error
	Retryable     bool
	RetryIn       time.Duration
	RateLimitInfo *RateLimitInfo // optional, only for rate limit errors
}

// Error implements the error interface
func (e *WarcraftLogsError) Error() string {
	msg := e.Message
	if e.RateLimitInfo != nil {
		msg += fmt.Sprintf(" (Points remaining: %.2f, Reset in: %v)",
			e.RateLimitInfo.RemainingPoints,
			e.RateLimitInfo.ResetIn)
	}
	if e.Cause != nil {
		msg += fmt.Sprintf(": %v", e.Cause)
	}
	return msg
}

// NewRateLimitError creates an error with rate limit information
func NewRateLimitError(info *RateLimitInfo, cause error) error {
	var retryIn time.Duration
	if info != nil {
		// If we have few points remaining, wait for the reset
		if info.RemainingPoints < 1.0 {
			retryIn = info.ResetIn
		} else {
			// Otherwise, wait a little for the points to regenerate
			retryIn = time.Second * 5
		}
	}

	return &WarcraftLogsError{
		Type:          ErrorTypeRateLimit,
		Message:       "API rate limit reached",
		Cause:         cause,
		Retryable:     true,
		RetryIn:       retryIn,
		RateLimitInfo: info,
	}
}

// NewQuotaExceededError creates an error when the hourly quota is exceeded
func NewQuotaExceededError(info *RateLimitInfo) error {
	return &WarcraftLogsError{
		Type:          ErrorTypeQuotaExceeded,
		Message:       "API quota exceeded for this hour",
		Retryable:     true,
		RetryIn:       info.ResetIn,
		RateLimitInfo: info,
	}
}

// NewAPIError creates a new WarcraftLogsError for API-related issues
func NewAPIError(statusCode int, cause error) error {
	message := fmt.Sprintf("API request failed with status code %d", statusCode)
	retryable := statusCode >= 500 || statusCode == 429 // Retryable if server error or rate limit exceeded

	return &WarcraftLogsError{
		Type:      ErrorTypeAPI,
		Message:   message,
		Cause:     cause,
		Retryable: retryable,
		RetryIn:   time.Second * 10, // Default retry delay
	}
}

// IsRateLimit checks if an error is a rate limit error
func IsRateLimit(err error) bool {
	if wlErr, ok := err.(*WarcraftLogsError); ok {
		return wlErr.Type == ErrorTypeRateLimit || wlErr.Type == ErrorTypeQuotaExceeded
	}
	return false
}

// GetRateLimitInfo extracts rate limit info from an error if available
func GetRateLimitInfo(err error) *RateLimitInfo {
	if wlErr, ok := err.(*WarcraftLogsError); ok {
		return wlErr.RateLimitInfo
	}
	return nil
}

// IsRetryable checks if an error should be retried
func IsRetryable(err error) bool {
	if wlErr, ok := err.(*WarcraftLogsError); ok {
		return wlErr.Retryable
	}
	return false
}

// GetRetryDelay returns the suggested retry delay for an error
func GetRetryDelay(err error, attempt int) time.Duration {
	if wlErr, ok := err.(*WarcraftLogsError); ok {
		if wlErr.RetryIn > 0 {
			return wlErr.RetryIn
		}
	}
	// Exponential backoff avec jitter
	baseDelay := time.Second * time.Duration(1<<uint(attempt))
	if baseDelay > time.Minute {
		baseDelay = time.Minute
	}
	return baseDelay
}
