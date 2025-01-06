package providers

import (
	"log"
	"os"
)

type DevProvider struct {
	logger *log.Logger
}

func NewDevProvider(logPath string) *DevProvider {
	return &DevProvider{}
}

func (p *DevProvider) Initialize() error {
	output := os.Stdout
	p.logger = log.New(output, "[EMAIL-DEV] ", log.LstdFlags)
	return nil
}

func (p *DevProvider) SendEmail(data EmailData) error {
	p.logger.Printf("================== EMAIL LOG ==================")
	p.logger.Printf("To: %s", data.To)
	p.logger.Printf("Subject: %s", data.Subject)
	p.logger.Printf("Content:\n%s", data.HTML)
	p.logger.Printf("=============================================")
	return nil
}

func (p *DevProvider) Close() error {
	return nil // No cleanup needed for dev provider
}
