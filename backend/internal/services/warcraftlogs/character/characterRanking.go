// package warcraftlogs/character/characterRanking.go
package character

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	models "wowperf/internal/models/warcraftlogs/character/ranking"
	service "wowperf/internal/services/warcraftlogs"
)

const CharacterRankingQuery = `
query getCharacterRanking($characterName: String!, $serverSlug: String!, $serverRegion: String!, $zoneID: Int!) {
    characterData {
        character(name: $characterName, serverSlug: $serverSlug, serverRegion: $serverRegion) {
            name
						classID
						id
            zoneRankings(zoneID: $zoneID)
        }
    }
}
`

// decodeCharacterName handles the decoding of character names that might be URL encoded
func DecodeCharacterName(name string) (string, error) {
	// First check if the name is actually encoded
	if !strings.Contains(name, "%") {
		return name, nil
	}

	// Handle double encoding (e.g., %25C3%25B8 -> ø)
	decodedName := name
	var err error
	for i := 0; i < 2; i++ { // Maximum of 2 decode passes
		decodedName, err = url.QueryUnescape(decodedName)
		if err != nil {
			return "", fmt.Errorf("failed to decode character name on pass %d: %w", i+1, err)
		}
		// If no more % characters are found, stop decoding
		if !strings.Contains(decodedName, "%") {
			break
		}
	}

	return decodedName, nil
}

// GetCharacterRanking returns the character ranking for a given character name, server slug, server region and zone ID
func GetCharacterRanking(s *service.WarcraftLogsClientService, characterName, serverSlug, serverRegion string, zoneID int) (*models.CharacterData, error) {
	// decode the character name
	decodedCharacterName, err := DecodeCharacterName(characterName)
	if err != nil {
		return nil, fmt.Errorf("failed to decode character name: %w", err)
	}

	// variables
	variables := map[string]interface{}{
		"characterName": decodedCharacterName,
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
