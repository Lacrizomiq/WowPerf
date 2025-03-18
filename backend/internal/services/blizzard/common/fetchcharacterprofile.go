// internal/services/blizzard/common/fetchcharacterprofile.go
package common

import (
	"errors"
	"fmt"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/profile"
	wrapper "wowperf/internal/wrapper/blizzard"
)

// FetchCharacterProfileData retrieves and processes character data in the same way
// as GetCharacterProfile but returns the data directly instead of sending an HTTP response.
// This can be used by internal services that need the character profile information.
func FetchCharacterProfileData(profileService *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (*models.CharacterProfile, error) {
	// Validate input parameters
	if region == "" {
		return nil, errors.New("region is required")
	}
	if realmSlug == "" {
		return nil, errors.New("realm slug is required")
	}
	if characterName == "" {
		return nil, errors.New("character name is required")
	}
	if namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if locale == "" {
		locale = "en_US" // Default locale if not provided
	}

	// Get character profile data
	characterData, err := profile.GetCharacterProfile(profileService, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve character profile: %w", err)
	}

	// Get character media data
	mediaData, err := profile.GetCharacterMedia(profileService, region, realmSlug, characterName, namespace, locale)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve character media: %w", err)
	}

	// Transform the data using the wrapper and return directly
	return wrapper.TransformCharacterInfo(characterData, mediaData)
}
