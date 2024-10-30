// package warcraftlogs/character/characterRanking.go
package character

import (
	"encoding/json"
	"fmt"
	models "wowperf/internal/models/warcraftlogs/character/ranking"
	service "wowperf/internal/services/warcraftlogs"
)

const CharacterRankingQuery = `
query getCharacterRanking($characterName: String!, $serverSlug: String!, $serverRegion: String!, $zoneID: Int!) {
    characterData {
        character(name: $characterName, serverSlug: $serverSlug, serverRegion: $serverRegion) {
            name
            zoneRankings(zoneID: $zoneID)
        }
    }
}
`

// GetCharacterRanking returns the character ranking for a given character name, server slug, server region and zone ID
func GetCharacterRanking(s *service.WarcraftLogsClientService, characterName, serverSlug, serverRegion string, zoneID int) (*models.CharacterData, error) {
	// variables
	variables := map[string]interface{}{
		"characterName": characterName,
		"serverSlug":    serverSlug,
		"serverRegion":  serverRegion,
		"zoneID":        zoneID,
	}

	// make the request
	response, err := s.Client.MakeGraphQLRequest(CharacterRankingQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get character ranking: %w", err)
	}

	// unmarshal the response
	var warcraftLogsResponse models.WarcraftLogsResponse
	if err := json.Unmarshal(response, &warcraftLogsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// return the character ranking
	return &warcraftLogsResponse.Data.CharacterData.Character, nil
}
