package email

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Sauvegarde des variables d'environnement originales
	originalEnv := map[string]string{
		"ENVIRONMENT":    os.Getenv("ENVIRONMENT"),
		"DOMAIN":         os.Getenv("DOMAIN"),
		"FRONTEND_URL":   os.Getenv("FRONTEND_URL"),
		"BACKEND_URL":    os.Getenv("BACKEND_URL"),
		"RESEND_API_KEY": os.Getenv("RESEND_API_KEY"),
		"MAILTRAP_USER":  os.Getenv("MAILTRAP_USER"),
		"MAILTRAP_PASS":  os.Getenv("MAILTRAP_PASS"),
		"EMAIL_LOG_PATH": os.Getenv("EMAIL_LOG_PATH"),
	}

	// Restauration des variables d'environnement à la fin
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		check   func(*testing.T, *Config)
	}{
		{
			name: "Local environment config",
			envVars: map[string]string{
				"ENVIRONMENT":  EnvLocal,
				"DOMAIN":       "localhost",
				"FRONTEND_URL": "https://localhost",
				"BACKEND_URL":  "https://localhost/api",
			},
			wantErr: false,
			check: func(t *testing.T, c *Config) {
				if c.Environment != EnvLocal {
					t.Errorf("Expected environment %s, got %s", EnvLocal, c.Environment)
				}
				if c.Domain != "localhost" {
					t.Errorf("Expected domain localhost, got %s", c.Domain)
				}
				if c.FrontendURL != "https://localhost" {
					t.Errorf("Expected frontend URL https://localhost, got %s", c.FrontendURL)
				}
			},
		},
		{
			name: "Test environment config",
			envVars: map[string]string{
				"ENVIRONMENT":   EnvTest,
				"DOMAIN":        "test.wowperf.com",
				"FRONTEND_URL":  "https://test.wowperf.com",
				"BACKEND_URL":   "https://test.wowperf.com/api",
				"MAILTRAP_USER": "testuser",
				"MAILTRAP_PASS": "testpass",
			},
			wantErr: false,
			check: func(t *testing.T, c *Config) {
				if c.Environment != EnvTest {
					t.Errorf("Expected environment %s, got %s", EnvTest, c.Environment)
				}
				if c.Mailtrap.Username != "testuser" {
					t.Errorf("Expected Mailtrap username testuser, got %s", c.Mailtrap.Username)
				}
			},
		},
		{
			name: "Production environment config",
			envVars: map[string]string{
				"ENVIRONMENT":    EnvProduction,
				"DOMAIN":         "wowperf.com",
				"FRONTEND_URL":   "https://wowperf.com",
				"BACKEND_URL":    "https://wowperf.com/api",
				"RESEND_API_KEY": "test-api-key",
			},
			wantErr: false,
			check: func(t *testing.T, c *Config) {
				if c.Environment != EnvProduction {
					t.Errorf("Expected environment %s, got %s", EnvProduction, c.Environment)
				}
				if c.Resend.APIKey != "test-api-key" {
					t.Errorf("Expected Resend API key test-api-key, got %s", c.Resend.APIKey)
				}
			},
		},
		{
			name: "Missing required DOMAIN",
			envVars: map[string]string{
				"ENVIRONMENT":  EnvLocal,
				"FRONTEND_URL": "https://localhost",
			},
			wantErr: true,
		},
		{
			name: "Missing required FRONTEND_URL",
			envVars: map[string]string{
				"ENVIRONMENT": EnvLocal,
				"DOMAIN":      "localhost",
			},
			wantErr: true,
		},
		{
			name: "Missing Resend API key in production",
			envVars: map[string]string{
				"ENVIRONMENT":  EnvProduction,
				"DOMAIN":       "wowperf.com",
				"FRONTEND_URL": "https://wowperf.com",
			},
			wantErr: true,
		},
		{
			name: "Missing Mailtrap credentials in test",
			envVars: map[string]string{
				"ENVIRONMENT":  EnvTest,
				"DOMAIN":       "test.wowperf.com",
				"FRONTEND_URL": "https://test.wowperf.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Nettoyage avant le test
			for key := range originalEnv {
				os.Unsetenv(key)
			}

			// Configuration des variables d'environnement pour le test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config, err := NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				tt.check(t, config)
			}
		})
	}
}

func TestConfig_GetResetPasswordURL(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		token  string
		want   string
	}{
		{
			name: "Local reset URL",
			config: &Config{
				Environment: EnvLocal,
				FrontendURL: "https://localhost",
			},
			token: "test-token",
			want:  "https://localhost/reset-password?token=test-token",
		},
		{
			name: "Test environment reset URL",
			config: &Config{
				Environment: EnvTest,
				FrontendURL: "https://test.wowperf.com",
			},
			token: "test-token",
			want:  "https://test.wowperf.com/reset-password?token=test-token",
		},
		{
			name: "Production reset URL",
			config: &Config{
				Environment: EnvProduction,
				FrontendURL: "https://wowperf.com",
			},
			token: "test-token",
			want:  "https://wowperf.com/reset-password?token=test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetResetPasswordURL(tt.token)
			if got != tt.want {
				t.Errorf("Config.GetResetPasswordURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_EnvironmentChecks(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		checks      map[string]bool
	}{
		{
			name:        "Local environment checks",
			environment: EnvLocal,
			checks: map[string]bool{
				"IsLocal":      true,
				"IsTest":       false,
				"IsProduction": false,
			},
		},
		{
			name:        "Test environment checks",
			environment: EnvTest,
			checks: map[string]bool{
				"IsLocal":      false,
				"IsTest":       true,
				"IsProduction": false,
			},
		},
		{
			name:        "Production environment checks",
			environment: EnvProduction,
			checks: map[string]bool{
				"IsLocal":      false,
				"IsTest":       false,
				"IsProduction": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Environment: tt.environment}

			if got := config.IsLocal(); got != tt.checks["IsLocal"] {
				t.Errorf("Config.IsLocal() = %v, want %v", got, tt.checks["IsLocal"])
			}
			if got := config.IsTest(); got != tt.checks["IsTest"] {
				t.Errorf("Config.IsTest() = %v, want %v", got, tt.checks["IsTest"])
			}
			if got := config.IsProduction(); got != tt.checks["IsProduction"] {
				t.Errorf("Config.IsProduction() = %v, want %v", got, tt.checks["IsProduction"])
			}
		})
	}
}
