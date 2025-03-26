package providers

import (
	"sync"
)

// MockProvider is a test provider that records sent emails
type MockProvider struct {
	initialized bool
	sentEmails  []EmailData
	mu          sync.RWMutex
	shouldFail  bool
}

func NewMockProvider(shouldFail bool) *MockProvider {
	return &MockProvider{
		sentEmails: make([]EmailData, 0),
		shouldFail: shouldFail,
	}
}

func (m *MockProvider) Initialize() error {
	m.initialized = true
	return nil
}

func (m *MockProvider) SendEmail(data EmailData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return &MockError{message: "mock provider failed to send email"}
	}

	m.sentEmails = append(m.sentEmails, data)
	return nil
}

func (m *MockProvider) Close() error {
	return nil
}

func (m *MockProvider) GetSentEmails() []EmailData {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid external modifications
	emails := make([]EmailData, len(m.sentEmails))
	copy(emails, m.sentEmails)
	return emails
}

// MockError is a custom error for testing
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}
