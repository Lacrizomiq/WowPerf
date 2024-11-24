// services/blizzard/protected/profile/service.go

package protectedProfile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
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

	// Construction de l'URL
	baseURL := fmt.Sprintf(apiURL, params.Region)
	u, err := url.Parse(baseURL + accountEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := url.Values{}
	query.Add("namespace", params.Namespace)
	query.Add("locale", params.Locale)
	u.RawQuery = query.Encode()

	data, err := s.protectedClient.MakeProtectedRequest(
		ctx,
		userID,
		u.String(),
		params.Namespace,
		params.Locale,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get WoW profile: %w", err)
	}
	log.Printf("Profile data: %s", string(data))

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse profile data: %w", err)
	}

	return result, nil
}
