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

// RateLimitTracker surveille l'utilisation de l'API de manière thread-safe
type RateLimitTracker struct {
	mu            sync.RWMutex
	requests      []time.Time // Timestamps des requêtes
	windowSize    time.Duration
	maxRequests   int
	totalRequests int64 // Compteur total depuis le démarrage
}

// NewRateLimitTracker crée un nouveau tracker
func NewRateLimitTracker(maxRequests int, window time.Duration) *RateLimitTracker {
	return &RateLimitTracker{
		requests:    make([]time.Time, 0),
		windowSize:  window,
		maxRequests: maxRequests,
	}
}

// RecordRequest enregistre une nouvelle requête et nettoie les anciennes
func (rt *RateLimitTracker) RecordRequest() {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	now := time.Now()
	rt.totalRequests++

	// Ajoute la nouvelle requête
	rt.requests = append(rt.requests, now)

	// Nettoie les requêtes trop anciennes (hors de la fenêtre)
	cutoff := now.Add(-rt.windowSize)
	validRequests := make([]time.Time, 0, len(rt.requests))

	for _, reqTime := range rt.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	rt.requests = validRequests
}

// GetStats retourne les statistiques actuelles thread-safe
func (rt *RateLimitTracker) GetStats() (currentRequests int, totalRequests int64, remainingCapacity int) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	current := len(rt.requests)
	remaining := rt.maxRequests - current
	if remaining < 0 {
		remaining = 0
	}

	return current, rt.totalRequests, remaining
}

// IsNearLimit vérifie si on approche de la limite (80% par exemple)
func (rt *RateLimitTracker) IsNearLimit(threshold float64) bool {
	current, _, _ := rt.GetStats()
	return float64(current) >= float64(rt.maxRequests)*threshold
}

type RaiderIOClient struct {
	httpClient  *http.Client
	baseURL     string
	limiter     *rate.Limiter
	rateTracker *RateLimitTracker
	apiKey      string
	logger      *log.Logger
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
	var rateLimitInterval time.Duration

	if apiKey != "" {
		// Avec API key : limite plus élevée
		maxRequests = 1000
		rateLimitInterval = time.Minute
		log.Printf("[INFO] RaiderIO client initialized with API key (limit: %d req/min)", maxRequests)
	} else {
		// Sans API key : 300 req/min
		maxRequests = 300
		rateLimitInterval = time.Minute
		log.Printf("[INFO] RaiderIO client initialized without API key (limit: %d req/min)", maxRequests)
	}

	client := &RaiderIOClient{
		httpClient: &http.Client{
			Timeout: 3 * time.Minute,
		},
		baseURL:     apiURL,
		limiter:     rate.NewLimiter(rate.Every(1*time.Second), 1),
		rateTracker: NewRateLimitTracker(maxRequests, rateLimitInterval),
		apiKey:      apiKey,
		logger:      log.New(os.Stdout, "[RAIDERIO] ", log.LstdFlags),
	}

	// Démarre le monitoring périodique mais plus discrètement
	client.startPeriodicLogging()

	return client, nil
}

// startPeriodicLogging démarre un goroutine pour logger les stats périodiquement
func (c *RaiderIOClient) startPeriodicLogging() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Log toutes les 5 minutes au lieu de 1
		defer ticker.Stop()

		for range ticker.C {
			current, total, remaining := c.rateTracker.GetStats()

			// Log seulement si il y a eu des requêtes
			if total > 0 {
				c.logger.Printf("Rate limit stats: %d/%d requests in last 5min, %d remaining, %d total since start",
					current, c.rateTracker.maxRequests, remaining, total)

				// Alerte seulement si vraiment proche de la limite
				if c.rateTracker.IsNearLimit(0.9) {
					c.logger.Printf("[WARNING] Very close to rate limit: %d/%d requests (90%% threshold)",
						current, c.rateTracker.maxRequests)
				}
			}
		}
	}()
}

// GetRateLimitStats retourne les statistiques actuelles (pour monitoring externe)
func (c *RaiderIOClient) GetRateLimitStats() (current int, total int64, remaining int, maxRequests int) {
	current, total, remaining = c.rateTracker.GetStats()
	return current, total, remaining, c.rateTracker.maxRequests
}

// doRequestWithRetry performs a request with rate limiting and retries
func (c *RaiderIOClient) doRequestWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	maxRetries := 3 // Réduit les retries pour éviter trop d'attente

	for i := 0; i < maxRetries; i++ {
		// Rate limiting non bloquant - juste pour le tracking
		c.rateTracker.RecordRequest()

		// Logger seulement si vraiment proche de la limite
		if c.rateTracker.IsNearLimit(0.95) {
			current, _, remaining, max := c.GetRateLimitStats()
			c.logger.Printf("[WARNING] Very high API usage: %d/%d requests, %d remaining", current, max, remaining)
		}

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
