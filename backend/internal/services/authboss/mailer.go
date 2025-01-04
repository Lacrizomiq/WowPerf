package authboss

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/volatiletech/authboss/v3"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// SMTPMailer implements authboss.Mailer
type SMTPMailer struct {
	config SMTPConfig
	auth   smtp.Auth
}

// NewSMTPMailer creates a new SMTP mailer
func NewSMTPMailer(config SMTPConfig) *SMTPMailer {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return &SMTPMailer{
		config: config,
		auth:   auth,
	}
}

// Send implements authboss.Mailer
func (m *SMTPMailer) Send(ctx context.Context, email authboss.Email) error {
	// Build email headers
	headers := make([]string, 0)
	headers = append(headers, fmt.Sprintf("From: %s <%s>", m.config.FromName, m.config.From))

	// Set To header
	var toAddresses []string
	for i, addr := range email.To {
		if len(email.ToNames) > i && email.ToNames[i] != "" {
			toAddresses = append(toAddresses, fmt.Sprintf("%s <%s>", email.ToNames[i], addr))
		} else {
			toAddresses = append(toAddresses, addr)
		}
	}
	headers = append(headers, fmt.Sprintf("To: %s", strings.Join(toAddresses, ", ")))

	headers = append(headers, fmt.Sprintf("Subject: %s", email.Subject))

	// Add MIME version and content type headers
	headers = append(headers, "MIME-Version: 1.0")

	// Build message body
	var body string
	if email.HTMLBody != "" {
		headers = append(headers, `Content-Type: multipart/alternative; boundary="boundary"`)
		body = fmt.Sprintf(`
--boundary
Content-Type: text/plain; charset=utf-8

%s

--boundary
Content-Type: text/html; charset=utf-8

%s
--boundary--
`, email.TextBody, email.HTMLBody)
	} else {
		headers = append(headers, `Content-Type: text/plain; charset=utf-8`)
		body = email.TextBody
	}

	// Combine headers and body
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	// Send email
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	return smtp.SendMail(
		addr,
		m.auth,
		m.config.From,
		email.To,
		[]byte(message),
	)
}
