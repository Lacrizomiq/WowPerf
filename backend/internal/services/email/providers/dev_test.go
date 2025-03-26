package providers

import (
	"testing"
)

func TestNewDevProvider(t *testing.T) {
	tests := []struct {
		name    string
		logPath string
		wantErr bool
	}{
		{
			name:    "Initialize with no log path",
			logPath: "",
			wantErr: false,
		},
		{
			name:    "Initialize with log path",
			logPath: "test.log",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDevProvider(tt.logPath)
			if err := p.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("DevProvider.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDevProvider_SendEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   EmailData
		wantErr bool
	}{
		{
			name: "Send valid email",
			email: EmailData{
				To:      "test@example.com",
				Subject: "Test Subject",
				HTML:    "<p>Test content</p>",
			},
			wantErr: false,
		},
		{
			name: "Send email with empty fields",
			email: EmailData{
				To:      "",
				Subject: "",
				HTML:    "",
			},
			wantErr: false, // Dev provider should not fail even with empty fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDevProvider("")
			if err := p.Initialize(); err != nil {
				t.Fatalf("Failed to initialize provider: %v", err)
			}

			if err := p.SendEmail(tt.email); (err != nil) != tt.wantErr {
				t.Errorf("DevProvider.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
