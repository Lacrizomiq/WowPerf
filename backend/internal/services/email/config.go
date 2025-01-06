package email

import (
	"fmt"
	"os"
	"strings"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvProduction  = "production"
)

// Config holds all email service configuration
type Config struct {
	Environment string
	BaseURL     string // Used for building reset password URL

	// Provider specific configuration
	Resend   ResendConfig
	Mailtrap MailtrapConfig
	Dev      DevConfig
}

// ResendConfig holds Resend specific configuration
type ResendConfig struct {
	APIKey string
	Domain string
}

// MailtrapConfig holds Mailtrap specific configuration
type MailtrapConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

// DevConfig holds development specific configuration
type DevConfig struct {
	LogPath string // Where to write email logs in development
}

// NewConfig creates a new email service configuration
func NewConfig() (*Config, error) {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "" {
		env = EnvDevelopment
	}

	baseURL := os.Getenv("APP_URL")
	if baseURL == "" {
		switch env {
		case EnvProduction:
			baseURL = "https://wowperf.com"
		case EnvTest:
			baseURL = "https://test.wowperf.com"
		default:
			baseURL = "http://localhost:8080"
		}
	}

	config := &Config{
		Environment: env,
		BaseURL:     baseURL,
	}

	// Load provider-specific configs based on environment
	switch env {
	case EnvProduction:
		if err := config.loadResendConfig(); err != nil {
			return nil, err
		}
	case EnvTest:
		if err := config.loadMailtrapConfig(); err != nil {
			return nil, err
		}
	case EnvDevelopment:
		config.loadDevConfig()
	}

	return config, nil
}

// loadResendConfig loads Resend configuration
func (c *Config) loadResendConfig() error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not set and is required for production environment")
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return fmt.Errorf("DOMAIN is not set and is required for production environment")
	}

	c.Resend = ResendConfig{
		APIKey: apiKey,
		Domain: domain,
	}

	return nil
}

// loadMailtrapConfig loads Mailtrap configuration
func (c *Config) loadMailtrapConfig() error {
	username := os.Getenv("MAILTRAP_USER")
	if username == "" {
		return fmt.Errorf("MAILTRAP_USER is not set and is required for test environment")
	}

	password := os.Getenv("MAILTRAP_PASS")
	if password == "" {
		return fmt.Errorf("MAILTRAP_PASS is not set and is required for test environment")
	}

	c.Mailtrap = MailtrapConfig{
		Username: username,
		Password: password,
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     2525,
	}

	return nil
}

func (c *Config) loadDevConfig() {
	c.Dev = DevConfig{
		LogPath: os.Getenv("EMAIL_LOG_PATH"), // Optional, will log to stdout if not set
	}
}

// GetResetPasswordURL generates the reset password URL for a given token
func (c *Config) GetResetPasswordURL(token string) string {
	return fmt.Sprintf("%s/reset-password?token=%s", c.BaseURL, token)
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// IsTest returns true if running in test environment
func (c *Config) IsTest() bool {
	return c.Environment == EnvTest
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}
