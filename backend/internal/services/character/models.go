package character

// Types spécifiques au domaine character qui ne sont pas dans models/
// Par exemple, résultats de sync, configurations, etc.

// SyncResult représente le résultat d'une synchronisation de personnages
type SyncResult struct {
	SyncedCount   int      `json:"synced_count"`
	EnrichedCount int      `json:"enriched_count"`
	UpdatedCount  int      `json:"updated_count"`
	Errors        []string `json:"errors,omitempty"`
}

// EnrichmentConfig configure quels enrichissements activer
type EnrichmentConfig struct {
	EnableSummary    bool `json:"enable_summary"`
	EnableEquipment  bool `json:"enable_equipment"`
	EnableMythicPlus bool `json:"enable_mythic_plus"`
	EnableRaids      bool `json:"enable_raids"`
}

// CharacterDetailLevel définit le niveau de détail souhaité
type CharacterDetailLevel int

const (
	Basic CharacterDetailLevel = iota
	WithSummary
	WithEquipment
	WithProgression
	Complete
)

// EnrichmentStatus représente le statut d'enrichissement d'un personnage
type EnrichmentStatus struct {
	Summary    bool `json:"summary"`
	Equipment  bool `json:"equipment"`
	MythicPlus bool `json:"mythic_plus"`
	Raids      bool `json:"raids"`
}

// CharacterUpdateRequest représente une demande de mise à jour de personnage
type CharacterUpdateRequest struct {
	CharacterID uint                 `json:"character_id"`
	DetailLevel CharacterDetailLevel `json:"detail_level"`
	ForceUpdate bool                 `json:"force_update"`
}

// BatchUpdateRequest représente une demande de mise à jour en lot
type BatchUpdateRequest struct {
	UserID      uint                     `json:"user_id"`
	Characters  []CharacterUpdateRequest `json:"characters,omitempty"`
	DetailLevel CharacterDetailLevel     `json:"detail_level"`
	ForceUpdate bool                     `json:"force_update"`
}
