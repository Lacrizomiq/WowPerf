package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Provider implementations
type (
	ResendEmailService struct {
		apiKey string
		domain string
	}

	MailtrapEmailService struct {
		username string
		password string
		host     string
		port     int
	}

	DevEmailService struct {
		logger *log.Logger
	}
)

// SendEmail implementation for ResendEmailService
func (s *ResendEmailService) SendEmail(data EmailData) error {
	emailData := struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
		Html    string   `json:"html"`
	}{
		From:    fmt.Sprintf("noreply@%s", s.domain),
		To:      []string{data.To},
		Subject: data.Subject,
		Html:    data.HTML,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resend API error: %d", resp.StatusCode)
	}

	return nil
}

// SendEmail implementation for MailtrapEmailService
func (s *MailtrapEmailService) SendEmail(data EmailData) error {
	log.Printf("[MAILTRAP] Sending email to: %s", data.To)
	log.Printf("[MAILTRAP] Subject: %s", data.Subject)
	log.Printf("[MAILTRAP] Content: %s", data.HTML)
	return nil
}

// SendEmail implementation for DevEmailService
func (s *DevEmailService) SendEmail(data EmailData) error {
	s.logger.Printf("================== EMAIL LOG ==================")
	s.logger.Printf("To: %s", data.To)
	s.logger.Printf("Subject: %s", data.Subject)
	s.logger.Printf("Content:\n%s", data.HTML)
	s.logger.Printf("=============================================")
	return nil
}

func NewEmailService(env string) EmailService {
	switch env {
	case "production":
		return &ResendEmailService{
			apiKey: os.Getenv("RESEND_API_KEY"),
			domain: os.Getenv("DOMAIN"),
		}
	case "test":
		return &MailtrapEmailService{
			username: os.Getenv("MAILTRAP_USER"),
			password: os.Getenv("MAILTRAP_PASS"),
			host:     "sandbox.smtp.mailtrap.io",
			port:     2525,
		}
	default:
		return &DevEmailService{
			logger: log.New(os.Stdout, "[EMAIL-DEV] ", log.LstdFlags),
		}
	}
}
