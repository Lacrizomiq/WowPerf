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

	var allRankings []playerLeaderboardModels.Ranking
	hasMorePages := false
	totalCount := 0

	// List of specializations to fetch
	var specsToFetch []Specialization

	// If a class and a spec are provided, only fetch that one
	if params.ClassName != "" && params.SpecName != "" {
		specsToFetch = append(specsToFetch, Specialization{params.ClassName, params.SpecName})
	} else {
		// Otherwise, fetch all specializations
		specsToFetch = Specializations
	}

	// Loop over each specialization to fetch the data
	for _, spec := range specsToFetch {
		variables := map[string]interface{}{
			"encounterId": params.EncounterID,
			"page":        params.Page,
			"className":   spec.ClassName,
			"specName":    spec.SpecName,
		}

		// Add optional parameters if provided
		if params.ServerRegion != "" {
			variables["serverRegion"] = params.ServerRegion
		}
		if params.ServerSlug != "" {
			variables["serverSlug"] = params.ServerSlug
		}

		// Execute the request
		response, err := s.MakeRequest(context.Background(), DungeonLeaderboardPlayerQuery, variables)
		if err != nil {
			log.Printf("Error fetching leaderboard for class %s, spec %s: %v", spec.ClassName, spec.SpecName, err)
			continue
		}

		// Check for GraphQL errors
		var errorResponse struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(response, &errorResponse); err == nil && len(errorResponse.Errors) > 0 {
			log.Printf("GraphQL error for class %s, spec %s: %s", spec.ClassName, spec.SpecName, errorResponse.Errors[0].Message)
			continue
		}

		// Define the response structure
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

		// Unmarshal the response
		if err := json.Unmarshal(response, &result); err != nil {
			log.Printf("Failed to unmarshal response for class %s, spec %s: %v", spec.ClassName, spec.SpecName, err)
			continue
		}

		// Check for empty results
		if len(result.Data.WorldData.Encounter.CharacterRankings) == 0 {
			continue
		}

		// Unmarshal the rankings
		var leaderboard playerLeaderboardModels.DungeonLogs
		if err := json.Unmarshal(result.Data.WorldData.Encounter.CharacterRankings, &leaderboard); err != nil {
			log.Printf("Failed to unmarshal character rankings for class %s, spec %s: %v", spec.ClassName, spec.SpecName, err)
			continue
		}

		// Add the results to the general array
		allRankings = append(allRankings, leaderboard.Rankings...)
		totalCount += leaderboard.Count
		if leaderboard.HasMorePages {
			hasMorePages = true
		}
	}

	// Return the leaderboard
	return &playerLeaderboardModels.DungeonLogs{
		Page:         params.Page,
		HasMorePages: hasMorePages,
		Count:        totalCount,
		Rankings:     allRankings,
	}, nil
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
