package captcha

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// CaptchaService gère la validations des captchas hCaptcha
type CaptchaService struct {
	secretKey string
	siteKey   string
	enabled   bool
	client    *http.Client
}

// hCaptchaResponse représente la réponse du captcha
type hCaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// NewCaptchaService crée une nouvelle instance de CaptchaService
func NewCaptchaService() *CaptchaService {
	enabled := strings.ToLower(os.Getenv("HCAPTCHA_ENABLED")) == "true"
	secretKey := os.Getenv("HCAPTCHA_SECRET_KEY")
	siteKey := os.Getenv("HCAPTCHA_SITE_KEY")

	service := &CaptchaService{
		secretKey: secretKey,
		siteKey:   siteKey,
		enabled:   enabled,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	log.Printf("[HCAPTCHA] CaptchaService initialized: %v", enabled)

	return service
}

// VerifyToken vérifie un token captcha auprès de hCaptcha
func (c *CaptchaService) VerifyToken(token string) error {
	// En développement, on ne vérifie pas le captcha
	if !c.enabled {
		log.Printf("[HCAPTCHA] Captcha verification skipped in development mode")
		return nil
	}

	if token == "" {
		return fmt.Errorf("captcha token is required")
	}

	if c.secretKey == "" {
		return fmt.Errorf("hCaptcha secret key is not configured")
	}

	// Préparer la requête vers hCaptcha
	data := url.Values{}
	data.Set("secret", c.secretKey)
	data.Set("response", token)

	log.Printf("[HCAPTCHA] Sending verification request to hCaptcha")

	resp, err := c.client.PostForm("https://hcaptcha.com/siteverify", data)
	if err != nil {
		log.Printf("[HCAPTCHA] HTTP request failed: %v", err)
		return fmt.Errorf("failed to verify captcha: %v", err)
	}
	defer resp.Body.Close()

	var hCaptchaResp hCaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&hCaptchaResp); err != nil {
		log.Printf("[HCAPTCHA] Failed to decode response: %v", err)
		return fmt.Errorf("failed to decode captcha response: %v", err)
	}

	if !hCaptchaResp.Success {
		log.Printf("[HCAPTCHA] Captcha verification failed: %v", hCaptchaResp.ErrorCodes)
		return fmt.Errorf("captcha verification failed: %v", hCaptchaResp.ErrorCodes)
	}

	log.Printf("[HCAPTCHA] Captcha verification successful")
	return nil
}

// IsEnabled vérifie si le captcha est activé
func (c *CaptchaService) IsEnabled() bool {
	return c.enabled
}

// GetSiteKey retourne la clé publique du captcha (Pour l'affichage du captcha dans le frontend si besoin)
func (c *CaptchaService) GetSiteKey() string {
	return c.siteKey
}

// ValidateConfig vérifie que la configuration est correcte
func (c *CaptchaService) ValidateConfig() error {
	if c.enabled {
		if c.secretKey == "" {
			return fmt.Errorf("HCAPTCHA_SECRET_KEY is required when captcha is enabled")
		}
		if c.siteKey == "" {
			return fmt.Errorf("HCAPTCHA_SITE_KEY is required when captcha is enabled")
		}
		log.Printf("[CAPTCHA] Configuration validated successfully")
	} else {
		log.Printf("[CAPTCHA] Running in disabled mode (development)")
	}
	return nil
}
