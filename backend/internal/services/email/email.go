package email

import (
	"fmt"
	"log"
	"sync"
	"wowperf/internal/models"
	"wowperf/internal/services/email/providers"
	"wowperf/internal/services/email/templates"
)

// EmailService manages email sending operations
type EmailService struct {
	config          *Config
	provider        providers.Provider
	templateManager *templates.TemplateManager
	mu              sync.RWMutex
}

// NewEmailService creates a new email service with the provided configuration
func NewEmailService(config *Config) (*EmailService, error) {
	service := &EmailService{
		config: config,
	}

	// Initialize template manager
	templateManager, err := templates.NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template manager: %w", err)
	}
	service.templateManager = templateManager

	// Initialize appropriate provider based on environment
	provider, err := service.initializeProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email provider: %w", err)
	}

	service.provider = provider
	return service, nil
}

// initializeProvider creates and initializes the appropriate email provider
func (s *EmailService) initializeProvider() (providers.Provider, error) {
	var provider providers.Provider

	switch {
	case s.config.IsProduction():
		provider = providers.NewResendProvider(
			s.config.Resend.APIKey,
			s.config.Resend.Domain,
		)
	case s.config.IsTest():
		provider = providers.NewMailtrapProvider(
			s.config.Mailtrap.Username,
			s.config.Mailtrap.Password,
			s.config.Mailtrap.Host,
			s.config.Mailtrap.Port,
		)
	default: // environnement local
		provider = providers.NewSMTPProvider(
			s.config.SMTP.Host,
			s.config.SMTP.Port,
			s.config.SMTP.Username,
			s.config.SMTP.Password,
			s.config.SMTP.From,
		)
	}

	if err := provider.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	return provider, nil
}

// SendPasswordResetEmail sends a password reset email to the user
func (s *EmailService) SendPasswordResetEmail(user *models.User, resetToken string) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	resetURL := s.config.GetResetPasswordURL(resetToken)

	// Prepare template data
	data := struct {
		Username string
		ResetURL string
	}{
		Username: user.Username,
		ResetURL: resetURL,
	}

	// Get email subject
	subject, err := s.templateManager.GetTemplateSubject(templates.PasswordReset)
	if err != nil {
		return fmt.Errorf("failed to get email subject: %w", err)
	}

	// Render the password reset template
	htmlContent, err := s.templateManager.RenderTemplate(templates.PasswordReset, data)
	if err != nil {
		return fmt.Errorf("failed to render password reset template: %w", err)
	}

	// Prepare email data
	emailData := providers.EmailData{
		To:      user.Email,
		Subject: subject,
		HTML:    htmlContent,
	}

	// Log email sending attempt
	log.Printf("Sending password reset email to %s (%s)", user.Username, user.Email)

	// Send the email
	s.mu.RLock()
	err = s.provider.SendEmail(emailData)
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Printf("Successfully sent password reset email to %s", user.Email)
	return nil
}

// Close cleans up any resources used by the email service
func (s *EmailService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.provider != nil {
		return s.provider.Close()
	}
	return nil
}

// Helper method to send a generic email
func (s *EmailService) SendEmail(data providers.EmailData) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.provider == nil {
		return fmt.Errorf("email provider not initialized")
	}

	return s.provider.SendEmail(data)
}
