package raiderio

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

const apiURL = "https://raider.io/api/v1"

// maskAPIKeyInURL masque l'API key dans l'URL pour les logs
func maskAPIKeyInURL(url string) string {
	if strings.Contains(url, "access_key=") {
		// Trouve la position de access_key=
		start := strings.Index(url, "access_key=")
		if start == -1 {
			return url
		}

		// Trouve la fin de la valeur (soit & soit fin de string)
		keyStart := start + len("access_key=")
		end := strings.Index(url[keyStart:], "&")

		var maskedURL string
		if end == -1 {
			// API key est à la fin de l'URL
			maskedURL = url[:keyStart] + "***MASKED***"
		} else {
			// Il y a d'autres paramètres après
			end += keyStart
			maskedURL = url[:keyStart] + "***MASKED***" + url[end:]
		}

		return maskedURL
	}
	return url
}

// UltraSimpleTracker - Juste un compteur avec alertes
type UltraSimpleTracker struct {
	mu            sync.RWMutex
	totalRequests int64
	startTime     time.Time
	logger        *log.Logger
	lastAlert     time.Time
}

func NewUltraSimpleTracker(logger *log.Logger) *UltraSimpleTracker {
	now := time.Now()
	tracker := &UltraSimpleTracker{
		startTime: now,
		logger:    logger,
		lastAlert: now,
	}

	// Log de démarrage
	logger.Printf("API client started - tracking requests")

	return tracker
}

func (ut *UltraSimpleTracker) RecordRequest() {
	ut.mu.Lock()
	defer ut.mu.Unlock()

	ut.totalRequests++
	now := time.Now()

	// Alertes occasionnelles basées sur les seuils
	ut.checkAlerts(now)
}

func (ut *UltraSimpleTracker) checkAlerts(now time.Time) {
	total := ut.totalRequests

	// Alertes à des seuils spécifiques (10, 50, 100, 200, etc.)
	alertThresholds := []int64{10, 50, 100, 200, 500, 1000}

	for _, threshold := range alertThresholds {
		if total == threshold {
			duration := now.Sub(ut.startTime)
			rate := float64(total) / duration.Hours()

			ut.logger.Printf("Milestone: %d requests completed (%.1f req/hour average since start)",
				total, rate)
			ut.lastAlert = now
			return
		}
	}

	// Alerte quotidienne si beaucoup d'activité
	if now.Sub(ut.lastAlert) >= 24*time.Hour && total > 0 {
		duration := now.Sub(ut.startTime)
		rate := float64(total) / duration.Hours()

		ut.logger.Printf("Daily summary: %d total requests (%.1f req/hour average)",
			total, rate)
		ut.lastAlert = now
	}
}

func (ut *UltraSimpleTracker) GetSummary() (total int64, duration time.Duration, avgPerHour float64) {
	ut.mu.RLock()
	defer ut.mu.RUnlock()

	duration = time.Since(ut.startTime)
	avgPerHour = float64(ut.totalRequests) / duration.Hours()

	return ut.totalRequests, duration, avgPerHour
}

// Méthode pour forcer un résumé (pour debug ou monitoring)
func (ut *UltraSimpleTracker) LogSummary() {
	total, duration, rate := ut.GetSummary()
	ut.logger.Printf("Current summary: %d requests in %s (%.1f req/hour)",
		total, duration.Round(time.Minute), rate)
}

type RaiderIOClient struct {
	httpClient     *http.Client
	baseURL        string
	limiter        *rate.Limiter
	requestTracker *UltraSimpleTracker
	apiKey         string
	logger         *log.Logger
}

type RaiderIOService struct {
	Client *RaiderIOClient
}

func NewRaiderIOClient() (*RaiderIOClient, error) {
	// Charge les variables d'environnement avec dotenv
	err := godotenv.Load()
	if err != nil {
		// Pas d'erreur fatale si .env n'existe pas
		log.Printf("[INFO] No .env file found, using system environment variables")
	}

	// Récupère l'API key depuis les variables d'environnement
	apiKey := os.Getenv("RAIDER_IO_API_KEY")

	// Détermine les limites selon la présence de l'API key
	var maxRequests int

	if apiKey != "" {
		// Avec API key : limite plus élevée
		maxRequests = 1000
		log.Printf("[INFO] RaiderIO client initialized with API key (limit: %d req/min)", maxRequests)
	} else {
		// Sans API key : 300 req/min
		maxRequests = 300
		log.Printf("[INFO] RaiderIO client initialized without API key (limit: %d req/min)", maxRequests)
	}

	logger := log.New(os.Stdout, "[RAIDERIO] ", log.LstdFlags)

	client := &RaiderIOClient{
		httpClient: &http.Client{
			Timeout: 3 * time.Minute,
		},
		baseURL:        apiURL,
		limiter:        rate.NewLimiter(rate.Every(1*time.Second), 1),
		requestTracker: NewUltraSimpleTracker(logger),
		apiKey:         apiKey,
		logger:         logger,
	}

	return client, nil
}

// GetRequestStats retourne les statistiques de requêtes (pour monitoring externe)
func (c *RaiderIOClient) GetRequestStats() (total int64, duration time.Duration, avgPerHour float64) {
	return c.requestTracker.GetSummary()
}

// LogRequestSummary force un log des statistiques actuelles
func (c *RaiderIOClient) LogRequestSummary() {
	c.requestTracker.LogSummary()
}

// doRequestWithRetry performs a request with rate limiting and retries
func (c *RaiderIOClient) doRequestWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	maxRetries := 3 // Réduit les retries pour éviter trop d'attente

	for i := 0; i < maxRetries; i++ {
		// Enregistrement simple de la requête
		c.requestTracker.RecordRequest()

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode != http.StatusGatewayTimeout {
			return resp, nil
		}

		if i < maxRetries-1 {
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			time.Sleep(waitTime)
		}
	}
	return nil, err
}

// Get makes a GET request to the Raider.io API (méthode originale pour compatibilité)
func (c *RaiderIOClient) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Add query parameters to the URL if they are provided
	if len(params) > 0 {
		url += "?"
		for key, value := range params {
			url += fmt.Sprintf("%s=%s&", key, value)
		}
		url = url[:len(url)-1] // Remove the last '&'
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "wowperf/1.0")

	// Log avec URL masquée pour la sécurité
	log.Printf("Making GET request to: %s", maskAPIKeyInURL(url))

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	log.Printf("Received response with status code: %d", resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusGatewayTimeout {
		return nil, fmt.Errorf("gateway timeout: the server took too long to respond")
	}

	// Parse the JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return result, nil
}

// GetRaw makes a GET request and returns raw JSON bytes (optimized version)
func (c *RaiderIOClient) GetRaw(endpoint string, params map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Ajoute automatiquement l'API key si disponible ET si pas déjà présente
	if c.apiKey != "" {
		if params == nil {
			params = make(map[string]string)
		}
		// Vérifie que l'API key n'est pas déjà dans les params
		if _, exists := params["access_key"]; !exists {
			params["access_key"] = c.apiKey
		}
	}

	// Add query parameters to the URL if they are provided
	if len(params) > 0 {
		url += "?"
		for key, value := range params {
			url += fmt.Sprintf("%s=%s&", key, value)
		}
		url = url[:len(url)-1] // Remove the last '&'
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "wowperf/1.0")

	// Log avec URL masquée pour la sécurité
	log.Printf("Making GET request to: %s", maskAPIKeyInURL(url))

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Received response with status code: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusGatewayTimeout {
		return nil, fmt.Errorf("gateway timeout: the server took too long to respond")
	}

	return body, nil
}

// GetWithAPIKey makes a GET request avec API key explicite (pour les nouvelles features)
func (c *RaiderIOClient) GetWithAPIKey(endpoint string, params map[string]string) (map[string]interface{}, error) {
	if c.apiKey != "" {
		if params == nil {
			params = make(map[string]string)
		}
		params["access_key"] = c.apiKey
	}

	return c.Get(endpoint, params)
}
