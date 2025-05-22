package googleauth

import (
	"fmt"
	"os"
)

// Config contient la configuration Google OAuth
type Config struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	FrontendURL   string
	DashboardPath string
	ErrorPath     string
}

// LoadConfig charge la configuration depuis les variables d'environnement
func LoadConfig() (*Config, error) {
	// Variable obligatoires
	required := map[string]string{
		"GOOGLE_CLIENT_ID":     os.Getenv("GOOGLE_CLIENT_ID"),
		"GOOGLE_CLIENT_SECRET": os.Getenv("GOOGLE_CLIENT_SECRET"),
		"GOOGLE_REDIRECT_URL":  os.Getenv("GOOGLE_REDIRECT_URL"),
		"FRONTEND_URL":         os.Getenv("FRONTEND_URL"),
	}

	// Vérifier que toutes les variables obligatoires sont présentes
	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("missing required env variable: %s", key)
		}
	}

	return &Config{
		ClientID:      required["GOOGLE_CLIENT_ID"],
		ClientSecret:  required["GOOGLE_CLIENT_SECRET"],
		RedirectURL:   required["GOOGLE_REDIRECT_URL"],
		FrontendURL:   required["FRONTEND_URL"],
		DashboardPath: getOrDefault("FRONTEND_DASHBOARD_PATH", "/"),
		ErrorPath:     getOrDefault("FRONTEND_AUTH_ERROR_PATH", "/login"),
	}, nil
}

// getOrDefault retourne la valeur de la variable d'environnement ou une valeur par défaut
func getOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate valide la configuration
func (c *Config) Validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("Google Client ID is required")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("Google Client Secret is required")
	}
	if c.RedirectURL == "" {
		return fmt.Errorf("Google Redirect URL is required")
	}
	if c.FrontendURL == "" {
		return fmt.Errorf("Frontend URL is required")
	}
	return nil
}
