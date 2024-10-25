// package warcraftlogs/service.go
package warcraftlogs

import (
	"encoding/json"
	"fmt"
	"log"

	leaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/team"
	warcraftlogsService "wowperf/internal/services/warcraftlogs"
)

const DungeonLeaderboardQuery = `
query getDungeonLeaderboard($encounterId: Int!, $page: Int!) {
    worldData {
        encounter(id: $encounterId) {
            name
            fightRankings(page: $page)
        }
    }
}`

type DungeonService struct {
	client *warcraftlogsService.Client
}

func NewDungeonService(client *warcraftlogsService.Client) *DungeonService {
	return &DungeonService{client: client}
}

// GetDungeonLeaderboard returns the dungeon leaderboard for a given encounter, region and page
func (s *DungeonService) GetDungeonLeaderboard(encounterID int, page int) (*leaderboardModels.DungeonLeaderboard, error) {
	log.Printf("Getting dungeon leaderboard for encounter %d, page %d", encounterID, page)

	variables := map[string]interface{}{
		"encounterId": encounterID,
		"page":        page,
	}

	response, err := s.client.MakeGraphQLRequest(DungeonLeaderboardQuery, variables)
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
		return &leaderboardModels.DungeonLeaderboard{
			Page:         page,
			HasMorePages: false,
			Count:        0,
			Rankings:     make([]leaderboardModels.Ranking, 0),
		}, nil
	}

	// Now unmarshal the FightRankings JSON into our DungeonLeaderboard struct
	var leaderboard leaderboardModels.DungeonLeaderboard
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
