package providers

import (
	"fmt"
	"net/smtp"
)

// MailtrapProvider is an implementation of the Provider interface for Mailtrap
type MailtrapProvider struct {
	username string
	password string
	host     string
	port     int
	auth     smtp.Auth
}

// NewMailtrapProvider creates a new MailtrapProvider
func NewMailtrapProvider(username, password, host string, port int) *MailtrapProvider {
	return &MailtrapProvider{
		username: username,
		password: password,
		host:     host,
		port:     port,
	}
}

// Initialize sets up the MailtrapProvider
func (p *MailtrapProvider) Initialize() error {
	if p.username == "" || p.password == "" {
		return fmt.Errorf("mailtrap credentials are required")
	}

	p.auth = smtp.PlainAuth("", p.username, p.password, p.host)
	return nil
}

func (p *MailtrapProvider) SendEmail(data EmailData) error {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", data.To, data.Subject, data.HTML))

	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	return smtp.SendMail(addr, p.auth, "noreply@wowperf.com", []string{data.To}, msg)
}

func (p *MailtrapProvider) Close() error {
	return nil // No cleanup needed for Mailtrap
}
