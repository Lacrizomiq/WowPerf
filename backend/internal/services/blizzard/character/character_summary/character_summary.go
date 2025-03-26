// internal/services/blizzard/character/character_summary/character_summary.go
package characterSummary

import (
	"fmt"
	"time"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/common"
)

// GetCharacterSummaryDetails fetches character summary details from Blizzard API
// Returns the updated character without persisting it
func GetCharacterSummaryDetails(
	profileService *blizzard.ProfileService,
	character *models.UserCharacter,
) error {
	// Params for the API call
	namespace := fmt.Sprintf("profile-%s", character.Region)
	locale := "en_US"

	// Fetch character profile data
	characterProfile, err := common.FetchCharacterProfileData(
		profileService,
		character.Region,
		character.Realm,
		character.Name,
		namespace,
		locale,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch character profile: %w", err)
	}

	// Update the UserCharacter model with the data obtained
	updateCharacterFromProfile(character, characterProfile)

	return nil
}

// updateCharacterFromProfile updates a UserCharacter with data from CharacterProfile
func updateCharacterFromProfile(character *models.UserCharacter, profile *models.CharacterProfile) {
	// Basic information update
	character.Name = profile.Name
	character.Race = profile.Race
	character.Class = profile.Class
	character.Faction = profile.Faction
	character.Gender = profile.Gender
	character.ActiveSpecName = profile.ActiveSpecName
	character.ActiveSpecID = profile.SpecID
	character.ActiveSpecRole = profile.ActiveSpecRole

	// Stats update
	character.AchievementPoints = profile.AchievementPoints
	character.HonorableKills = profile.HonorableKills

	// Avatar URLs update
	character.AvatarURL = profile.AvatarURL
	character.InsetAvatarURL = profile.InsetAvatarURL
	character.MainRawURL = profile.MainRawUrl
	character.ProfileURL = profile.ProfileURL

	// Update timestamp
	character.LastAPIUpdate = time.Now()
}
