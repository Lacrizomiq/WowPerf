package enrichers

import (
	"context"
	"fmt"
	"time"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/common"
)

// SummaryEnricher enrichit les données de base d'un personnage
type SummaryEnricher struct {
	profileService *blizzard.ProfileService
}

// NewSummaryEnricher crée un nouvel enrichisseur de summary
func NewSummaryEnricher(profileService *blizzard.ProfileService) *SummaryEnricher {
	return &SummaryEnricher{
		profileService: profileService,
	}
}

// EnrichCharacter enrichit un personnage avec ses données de summary
func (e *SummaryEnricher) EnrichCharacter(ctx context.Context, character *models.UserCharacter) error {
	// Params pour l'appel API
	namespace := fmt.Sprintf("profile-%s", character.Region)
	locale := "en_US"

	// Récupérer les données du profil depuis l'API Blizzard
	characterProfile, err := common.FetchCharacterProfileData(
		e.profileService,
		character.Region,
		character.Realm,
		character.Name,
		namespace,
		locale,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch character profile: %w", err)
	}

	// Mettre à jour le UserCharacter avec les données obtenues
	updateCharacterFromProfile(character, characterProfile)

	return nil
}

// GetName retourne le nom de cet enrichisseur
func (e *SummaryEnricher) GetName() string {
	return "summary"
}

// GetPriority retourne la priorité d'exécution (summary en premier)
func (e *SummaryEnricher) GetPriority() int {
	return 1
}

// CanEnrich vérifie si cet enrichisseur peut traiter ce personnage
func (e *SummaryEnricher) CanEnrich(character *models.UserCharacter) bool {
	// Le summary peut enrichir tous les personnages
	return character.Name != "" && character.Realm != "" && character.Region != ""
}

// updateCharacterFromProfile met à jour un UserCharacter avec les données d'un CharacterProfile
func updateCharacterFromProfile(character *models.UserCharacter, profile *models.CharacterProfile) {
	// Mise à jour des informations de base
	character.Name = profile.Name
	character.Race = profile.Race
	character.Class = profile.Class
	character.Faction = profile.Faction
	character.Gender = profile.Gender
	character.ActiveSpecName = profile.ActiveSpecName
	character.ActiveSpecID = profile.SpecID
	character.ActiveSpecRole = profile.ActiveSpecRole

	// Mise à jour des statistiques
	character.AchievementPoints = profile.AchievementPoints
	character.HonorableKills = profile.HonorableKills

	// Mise à jour des URLs d'avatar
	character.AvatarURL = profile.AvatarURL
	character.InsetAvatarURL = profile.InsetAvatarURL
	character.MainRawURL = profile.MainRawUrl
	character.ProfileURL = profile.ProfileURL

	// Mise à jour du timestamp
	character.LastAPIUpdate = time.Now()
}
