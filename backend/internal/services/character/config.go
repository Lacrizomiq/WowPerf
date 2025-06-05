package character

import "time"

// Configuration constants for character services
const (
	// Rate limiting
	MaxSyncPerDay       = 10000           // Maximum syncs per user per day
	MinDelayBetweenSync = 0 * time.Minute // Minimum delay between syncs
	MaxEnrichPerHour    = 300000          // Maximum enrichments per user per hour
	MaxRetries          = 3               // Maximum retry attempts

	// Performance
	BatchSize      = 5                 // Number of characters to process in batch
	RequestTimeout = 300 * time.Second // Timeout for API requests

	// Enrichers configuration
	EnableSummary    = true  // Enable summary enrichment
	EnableEquipment  = false // Enable equipment enrichment (for future)
	EnableMythicPlus = false // Enable M+ enrichment (for future)
	EnableRaids      = false // Enable raids enrichment (for future)
)
