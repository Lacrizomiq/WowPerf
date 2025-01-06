// internal/services/email/types.go
package email

// EmailData represents the structure for emails
type EmailData struct {
	To      string
	Subject string
	HTML    string
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendEmail(data EmailData) error
}
