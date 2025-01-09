package providers

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPProvider is a provider for sending emails via SMTP
type SMTPProvider struct {
	host     string
	port     int
	username string
	password string
	from     string
	auth     smtp.Auth
}

// NewSMTPProvider creates a new SMTPProvider instance
func NewSMTPProvider(host string, port int, username, password, from string) *SMTPProvider {
	return &SMTPProvider{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Initialize sets up the SMTP authentication if credentials are provided
func (p *SMTPProvider) Initialize() error {
	if p.username != "" && p.password != "" {
		p.auth = smtp.PlainAuth("", p.username, p.password, p.host)
	}
	return nil
}

// SendEmail sends an email using SMTP
func (p *SMTPProvider) SendEmail(data EmailData) error {
	// Building email headers
	headers := make([]string, 0)
	headers = append(headers, fmt.Sprintf("From: %s", p.from))
	headers = append(headers, fmt.Sprintf("To: %s", data.To))
	headers = append(headers, fmt.Sprintf("Subject: %s", data.Subject))
	headers = append(headers, "MIME-version: 1.0")
	headers = append(headers, "Content-Type: text/html; charset=\"UTF-8\"")

	// Building email message
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + data.HTML

	// SMTP server address
	addr := fmt.Sprintf("%s:%d", p.host, p.port)

	// Sending email
	if p.auth != nil {
		return smtp.SendMail(addr, p.auth, p.from, []string{data.To}, []byte(message))
	}

	// For MailHog in local, no authentication is needed
	return smtp.SendMail(addr, nil, p.from, []string{data.To}, []byte(message))
}

// Close implements the Provider interface
func (p *SMTPProvider) Close() error {
	return nil // No resources to clean up for SMTP
}
