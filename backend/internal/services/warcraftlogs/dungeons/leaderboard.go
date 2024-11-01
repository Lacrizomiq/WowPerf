// package warcraftlogs/service.go
package warcraftlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	playerLeaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/player"
	teamLeaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/team"
	service "wowperf/internal/services/warcraftlogs"
)

const DungeonLeaderboardTeamQuery = `
query getDungeonLeaderboard($encounterId: Int!, $page: Int!) {
    worldData {
        encounter(id: $encounterId) {
            name
            fightRankings(page: $page)
        }
    }
}`

const DungeonLeaderboardPlayerQuery = `
query getDungeonLeaderboard(
    $encounterId: Int!,
    $page: Int!,
    $serverRegion: String,
		$serverSlug: String,
		$className: String,
		$specName: String
) {
    worldData {
        encounter(id: $encounterId) {
            name
            characterRankings(
							leaderboard: Any,
							page: $page,
							serverRegion: $serverRegion,
							serverSlug: $serverSlug,
							className: $className,
							specName: $specName
						)
        }
    }
}`

type LeaderboardParams struct {
	EncounterID  int
	Page         int
	ServerRegion string
	ServerSlug   string
	ClassName    string
	SpecName     string
}

// GetDungeonLeaderboardByPlayer returns the dungeon leaderboard for a given encounter, region and page
func GetDungeonLeaderboardByPlayer(s *service.WarcraftLogsClientService, params LeaderboardParams) (*playerLeaderboardModels.DungeonLogs, error) {
	log.Printf("Getting dungeon leaderboard for encounter %d, page %d", params.EncounterID, params.Page)

	variables := map[string]interface{}{
		"encounterId": params.EncounterID,
		"page":        params.Page,
	}

	// add optional parameters if they are provided
	if params.ServerRegion != "" {
		variables["serverRegion"] = params.ServerRegion
	}
	if params.ServerSlug != "" {
		variables["serverSlug"] = params.ServerSlug
	}
	if params.ClassName != "" {
		variables["className"] = params.ClassName
	}
	if params.SpecName != "" {
		variables["specName"] = params.SpecName
	}

	// make the request
	response, err := s.MakeRequest(context.Background(), DungeonLeaderboardPlayerQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get dungeon leaderboard: %w", err)
	}

	// log raw response for debugging
	log.Printf("Raw response: %s", string(response))

	// Check for GraphQL errors
	var errorResponse struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(response, &errorResponse); err == nil && len(errorResponse.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", errorResponse.Errors[0].Message)
	}

	// Define the response struct
	var result struct {
		Data struct {
			WorldData struct {
				Encounter struct {
					Name              string          `json:"name"`
					CharacterRankings json.RawMessage `json:"characterRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		} `json:"data"`
	}

	// Unmarshal the response into the struct
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verify if CharacterRankings is null or empty
	if len(result.Data.WorldData.Encounter.CharacterRankings) == 0 {
		return &playerLeaderboardModels.DungeonLogs{
			Page:         params.Page,
			HasMorePages: false,
			Count:        0,
			Rankings:     make([]playerLeaderboardModels.Ranking, 0),
		}, nil
	}

	// Now unmarshal the CharacterRankings JSON into our DungeonLeaderboard struct
	var leaderboard playerLeaderboardModels.DungeonLogs
	if err := json.Unmarshal(result.Data.WorldData.Encounter.CharacterRankings, &leaderboard); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character rankings: %w", err)
	}

	// log the leaderboard
	log.Printf("Leaderboard for %s: Page %d, Count %d, HasMore %v, Rankings: %d",
		result.Data.WorldData.Encounter.Name,
		leaderboard.Page,
		leaderboard.Count,
		leaderboard.HasMorePages,
		len(leaderboard.Rankings))

	return &leaderboard, nil
}

// GetDungeonLeaderboardByTeam returns the dungeon leaderboard for a given encounter and page
func GetDungeonLeaderboardByTeam(s *service.WarcraftLogsClientService, encounterID int, page int) (*teamLeaderboardModels.DungeonLeaderboard, error) {
	log.Printf("Getting dungeon leaderboard for encounter %d, page %d", encounterID, page)

	variables := map[string]interface{}{
		"encounterId": encounterID,
		"page":        page,
	}

	response, err := s.MakeRequest(context.Background(), DungeonLeaderboardTeamQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get dungeon leaderboard: %w", err)
	}

	// Log raw response for debugging
	log.Printf("Raw response: %s", string(response))

	// First check if there are any GraphQL errors
	var errorResponse struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(response, &errorResponse); err == nil && len(errorResponse.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", errorResponse.Errors[0].Message)
	}

	var result struct {
		Data struct {
			WorldData struct {
				Encounter struct {
					Name          string          `json:"name"`
					FightRankings json.RawMessage `json:"fightRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initial response: %w", err)
	}

	// Verify if FightRankings is null or empty
	if len(result.Data.WorldData.Encounter.FightRankings) == 0 {
		return &teamLeaderboardModels.DungeonLeaderboard{
			Page:         page,
			HasMorePages: false,
			Count:        0,
			Rankings:     make([]teamLeaderboardModels.Ranking, 0),
		}, nil
	}

	// Now unmarshal the FightRankings JSON into our DungeonLeaderboard struct
	var leaderboard teamLeaderboardModels.DungeonLeaderboard
	if err := json.Unmarshal(result.Data.WorldData.Encounter.FightRankings, &leaderboard); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fight rankings: %w", err)
	}

	log.Printf("Leaderboard for %s: Page %d, Count %d, HasMore %v, Rankings: %d",
		result.Data.WorldData.Encounter.Name,
		leaderboard.Page,
		leaderboard.Count,
		leaderboard.HasMorePages,
		len(leaderboard.Rankings))

	if leaderboard.Rankings != nil {
		for i, ranking := range leaderboard.Rankings {
			log.Printf("Ranking %d: Score %.2f, Medal %s, Team size %d",
				i+1, ranking.Score, ranking.Medal, len(ranking.Team))
		}
	}

	return &leaderboard, nil
}
