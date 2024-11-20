// package warcraftlogs/character/characterEncounterRanking.go

package character

import (
	"encoding/json"
	"fmt"
	characterRaidRankingByEncounter "wowperf/internal/models/warcraftlogs/character/raidRanking"
	service "wowperf/internal/services/warcraftlogs"
)

const CharacterEncounterRankingQuery = `
query getCharacterEncounterRanking($characterName: String!, $serverSlug: String!, $serverRegion: String!, $encounterID: Int!, $includeCombatantInfo: Boolean!) {
    characterData {
        character(name: $characterName, serverSlug: $serverSlug, serverRegion: $serverRegion) {
            name
            classID
            id
            encounterRankings(encounterID: $encounterID, includeCombatantInfo: $includeCombatantInfo)
        }
    }
}
`

// GetCharacterEncounterRanking returns the encounter ranking for a given character
func GetCharacterEncounterRanking(s *service.WarcraftLogsClientService, characterName, serverSlug, serverRegion string, encounterID int, includeCombatantInfo bool) (*characterRaidRankingByEncounter.Character, error) {
	// Prepare variables for the GraphQL query
	variables := map[string]interface{}{
		"characterName":        characterName,
		"serverSlug":           serverSlug,
		"serverRegion":         serverRegion,
		"encounterID":          encounterID,
		"includeCombatantInfo": includeCombatantInfo,
	}

	// Make the request
	response, err := s.Client.MakeGraphQLRequest(CharacterEncounterRankingQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get character encounter ranking: %w", err)
	}

	// Unmarshal the response
	var warcraftLogsResponse characterRaidRankingByEncounter.WarcraftLogsResponse
	if err := json.Unmarshal(response, &warcraftLogsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if warcraftLogsResponse.Data.CharacterData.Character.Name == "" {
		return nil, fmt.Errorf("character not found")
	}

	// Return the character data
	return &warcraftLogsResponse.Data.CharacterData.Character, nil
}
