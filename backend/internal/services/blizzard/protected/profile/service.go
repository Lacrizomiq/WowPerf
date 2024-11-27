// services/blizzard/protected/profile/service.go

package protectedProfile

import (
	"context"
	"encoding/json"
	"fmt"
	"wowperf/internal/services/blizzard/types"
)

const (
	apiURL          = "https://%s.api.blizzard.com"
	accountEndpoint = "/profile/user/wow"
)

type ProtectedProfileService struct {
	protectedClient types.ProtectedClientInterface
}

func NewProtectedProfileService(protectedClient types.ProtectedClientInterface) *ProtectedProfileService {
	return &ProtectedProfileService{
		protectedClient: protectedClient,
	}
}

func (s *ProtectedProfileService) GetAccountProfile(ctx context.Context, userID uint, params types.ProfileServiceParams) (map[string]interface{}, error) {
	// Validation
	if params.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	if params.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if params.Locale == "" {
		params.Locale = "en_US"
	}

	// Construction de l'endpoint (sans le baseURL)
	endpoint := "/profile/user/wow"

	// Appel au client protégé avec juste l'endpoint
	data, err := s.protectedClient.MakeProtectedRequest(
		ctx,
		userID,
		endpoint, // Juste l'endpoint, pas l'URL complète
		params.Namespace,
		params.Locale,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get WoW profile: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse profile data: %w", err)
	}

	return result, nil
}

func (s *ProtectedProfileService) GetProtectedCharacterProfile(
	ctx context.Context,
	userID uint,
	realmId string,
	characterId string,
	params types.ProfileServiceParams,
) (map[string]interface{}, error) {
	// Validation
	if params.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	if params.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	if params.Locale == "" {
		params.Locale = "en_US"
	}

	// Build the endpoint
	endpoint := fmt.Sprintf("/profile/user/wow/protected-character/%s-%s", realmId, characterId)

	// Call the protected client
	data, err := s.protectedClient.MakeProtectedRequest(
		ctx,
		userID,
		endpoint,
		params.Namespace,
		params.Locale,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get protected character profile: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse profile data: %w", err)
	}

	return result, nil
}
