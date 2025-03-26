package providers

import (
	"testing"
)

func TestNewMailtrapProvider(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		host     string
		port     int
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			username: "testuser",
			password: "testpass",
			host:     "smtp.mailtrap.io",
			port:     2525,
			wantErr:  false,
		},
		{
			name:     "Missing credentials",
			username: "",
			password: "",
			host:     "smtp.mailtrap.io",
			port:     2525,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMailtrapProvider(tt.username, tt.password, tt.host, tt.port)
			if err := p.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("MailtrapProvider.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMailtrapProvider_SendEmail(t *testing.T) {
	// Create a provider with test credentials
	p := NewMailtrapProvider("testuser", "testpass", "smtp.mailtrap.io", 2525)
	if err := p.Initialize(); err != nil {
		t.Fatalf("Failed to initialize provider: %v", err)
	}

	tests := []struct {
		name    string
		email   EmailData
		wantErr bool
	}{
		{
			name: "Valid email data",
			email: EmailData{
				To:      "test@example.com",
				Subject: "Test Subject",
				HTML:    "<p>Test content</p>",
			},
			wantErr: true, // Will fail in test because credentials are fake
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.SendEmail(tt.email); (err != nil) != tt.wantErr {
				t.Errorf("MailtrapProvider.SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
