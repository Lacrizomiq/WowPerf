package email

import (
	"testing"
	"time"
	"wowperf/internal/models"

	"gorm.io/gorm"
)

func TestNewEmailService(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid local config",
			config: &Config{
				Environment: EnvLocal,
				Domain:      "localhost",
				FrontendURL: "https://localhost",
				BackendURL:  "https://localhost/api",
				Dev: DevConfig{
					LogPath: "test.log",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid test config",
			config: &Config{
				Environment: EnvTest,
				Domain:      "test.wowperf.com",
				FrontendURL: "https://test.wowperf.com",
				BackendURL:  "https://test.wowperf.com/api",
				Mailtrap: MailtrapConfig{
					Username: "test",
					Password: "test",
					Host:     "smtp.mailtrap.io",
					Port:     2525,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid production config (missing API key)",
			config: &Config{
				Environment: EnvProduction,
				Domain:      "wowperf.com",
				FrontendURL: "https://wowperf.com",
				BackendURL:  "https://wowperf.com/api",
				Resend: ResendConfig{
					APIKey: "",
					Domain: "wowperf.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewEmailService(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmailService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && service == nil {
				t.Error("NewEmailService() returned nil service without error")
			}
		})
	}
}

func TestEmailService_SendPasswordResetEmail(t *testing.T) {
	// Configuration de test
	config := &Config{
		Environment: EnvLocal,
		Domain:      "localhost",
		FrontendURL: "https://localhost",
		BackendURL:  "https://localhost/api",
		Dev: DevConfig{
			LogPath: "test.log",
		},
	}

	service, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Failed to create email service: %v", err)
	}

	tests := []struct {
		name       string
		user       *models.User
		resetToken string
		wantErr    bool
	}{
		{
			name: "Valid reset email",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			resetToken: "valid-token",
			wantErr:    false,
		},
		{
			name:       "Nil user",
			user:       nil,
			resetToken: "token",
			wantErr:    true,
		},
		{
			name: "Empty email",
			user: &models.User{
				Username: "testuser",
				Email:    "",
			},
			resetToken: "token",
			wantErr:    true,
		},
		{
			name: "Empty token",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			resetToken: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SendPasswordResetEmail(tt.user, tt.resetToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailService.SendPasswordResetEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmailService_Integration(t *testing.T) {
	// Skip in CI/CD if necessary
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a test user
	user := &models.User{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Test config
	config := &Config{
		Environment: EnvLocal,
		Domain:      "localhost",
		FrontendURL: "https://localhost",
		BackendURL:  "https://localhost/api",
		Dev: DevConfig{
			LogPath: "test.log",
		},
	}

	// Create the service
	service, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Failed to create email service: %v", err)
	}

	// Test the complete password reset flow
	t.Run("Complete password reset flow", func(t *testing.T) {
		// 1. Generate the token
		token, err := user.GeneratePasswordResetToken()
		if err != nil {
			t.Fatalf("Failed to generate reset token: %v", err)
		}

		// 2. Send the email
		if err := service.SendPasswordResetEmail(user, token); err != nil {
			t.Errorf("Failed to send password reset email: %v", err)
		}

		// 3. Verify the token
		if !user.ValidatePasswordResetToken(token) {
			t.Error("Token validation failed")
		}

		// 4. Check token expiration
		if user.IsPasswordResetTokenExpired() {
			t.Error("Token expired too soon")
		}
	})
}

func TestEmailService_Close(t *testing.T) {
	config := &Config{
		Environment: EnvLocal,
		Domain:      "localhost",
		FrontendURL: "https://localhost",
		BackendURL:  "https://localhost/api",
	}

	service, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Failed to create email service: %v", err)
	}

	if err := service.Close(); err != nil {
		t.Errorf("EmailService.Close() error = %v", err)
	}
}
