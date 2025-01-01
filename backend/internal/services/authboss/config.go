package authboss

import "os"

type Config struct {
	Environment string
	RootURL     string
	SessionKey  string
}

func NewConfig() *Config {
	env := os.Getenv("ENVIRONMENT")

	var rootURL string
	switch env {
	case "production":
		rootURL = "https://wowperf.com"
	case "test":
		rootURL = "https://test.wowperf.com"
	default:
		rootURL = "https://localhost/api"
	}

	return &Config{
		Environment: env,
		RootURL:     rootURL,
		SessionKey:  os.Getenv("SESSION_KEY"),
	}
}
