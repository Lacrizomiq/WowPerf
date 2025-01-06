package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ResendProvider struct {
	apiKey string
	domain string
	client *http.Client
}

func NewResendProvider(apiKey, domain string) *ResendProvider {
	return &ResendProvider{
		apiKey: apiKey,
		domain: domain,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *ResendProvider) Initialize() error {
	if p.apiKey == "" {
		return fmt.Errorf("resend API key is required")
	}
	if p.domain == "" {
		return fmt.Errorf("domain is required")
	}
	return nil
}

func (p *ResendProvider) SendEmail(data EmailData) error {
	emailData := struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		Subject string   `json:"subject"`
		Html    string   `json:"html"`
	}{
		From:    fmt.Sprintf("noreply@%s", p.domain),
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

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resend API error: status code %d", resp.StatusCode)
	}

	return nil
}

func (p *ResendProvider) Close() error {
	return nil // No cleanup needed for Resend
}
