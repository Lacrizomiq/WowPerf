package email

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Environment constants
const (
	EnvLocal      = "local"
	EnvTest       = "test"
	EnvProduction = "production"
)

// Config holds all email service configuration
type Config struct {
	Environment string
	Domain      string
	FrontendURL string // Used for building reset password URLs
	BackendURL  string

	// Provider specific configurations
	Resend   ResendConfig
	Mailtrap MailtrapConfig
	Dev      DevConfig
	SMTP     SMTPConfig
}

// ResendConfig holds Resend-specific configuration
type ResendConfig struct {
	APIKey string
	Domain string
}

// MailtrapConfig holds Mailtrap-specific configuration
type MailtrapConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

// SMTPConfig holds SMTP-specific configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// DevConfig holds development-specific configuration
type DevConfig struct {
	LogPath string
}

// NewConfig creates a new email service configuration
func NewConfig() (*Config, error) {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = EnvLocal
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("DOMAIN environment variable is required")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return nil, fmt.Errorf("FRONTEND_URL environment variable is required")
	}

	config := &Config{
		Environment: env,
		Domain:      domain,
		FrontendURL: frontendURL,
		BackendURL:  os.Getenv("BACKEND_URL"),
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     getEnvAsInt("SMTP_PORT", 1025),
			Username: os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},
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
	case EnvLocal:
		config.loadDevConfig()
	}

	return config, nil
}

func (c *Config) loadResendConfig() error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is required in production environment")
	}

	c.Resend = ResendConfig{
		APIKey: apiKey,
		Domain: c.Domain,
	}

	return nil
}

func (c *Config) loadMailtrapConfig() error {
	username := os.Getenv("MAILTRAP_USER")
	if username == "" {
		return fmt.Errorf("MAILTRAP_USER is required in test environment")
	}

	password := os.Getenv("MAILTRAP_PASS")
	if password == "" {
		return fmt.Errorf("MAILTRAP_PASS is required in test environment")
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
		LogPath: os.Getenv("EMAIL_LOG_PATH"),
	}
}

// GetResetPasswordURL generates the reset password URL for a given token
func (c *Config) GetResetPasswordURL(token string) string {
	return fmt.Sprintf("%s/reset-password?token=%s", c.FrontendURL, token)
}

// Environment check methods
func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

func (c *Config) IsTest() bool {
	return c.Environment == EnvTest
}

func (c *Config) IsLocal() bool {
	return c.Environment == EnvLocal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
