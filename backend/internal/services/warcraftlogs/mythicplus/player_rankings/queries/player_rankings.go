package warcraftlogsPlayerRankingsQueries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	playerLeaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/player"
	service "wowperf/internal/services/warcraftlogs"
	workflowModels "wowperf/internal/services/warcraftlogs/mythicplus/player_rankings/temporal/workflows/models"
)

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
	var specsToFetch []workflowModels.ClassSpec

	// If a class and a spec are provided, only fetch that one
	if params.ClassName != "" && params.SpecName != "" {
		specsToFetch = append(specsToFetch, workflowModels.ClassSpec{ClassName: params.ClassName, SpecName: params.SpecName})
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
			WorldData struct {
				Encounter struct {
					Name              string `json:"name"`
					CharacterRankings struct {
						Rankings     []playerLeaderboardModels.Ranking `json:"rankings"`
						Count        int                               `json:"count"`
						Page         int                               `json:"page"`
						HasMorePages bool                              `json:"hasMorePages"`
					} `json:"characterRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		}

		// Unmarshal the response
		if err := json.Unmarshal(response, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response for class %s, spec %s: %w", spec.ClassName, spec.SpecName, err)
		}

		// Verify rankings
		characterRankings := result.WorldData.Encounter.CharacterRankings
		if characterRankings.Rankings == nil {
			log.Printf("No rankings found for class %s, spec %s", spec.ClassName, spec.SpecName)
			continue
		}

		// Add the results to the general array
		allRankings = append(allRankings, characterRankings.Rankings...)
		totalCount += characterRankings.Count
		if characterRankings.HasMorePages {
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

// List of all specializations mapped to their classes, using the models.ClassSpec structure
var Specializations = []workflowModels.ClassSpec{
	// Priest
	{ClassName: "Priest", SpecName: "Discipline"},
	{ClassName: "Priest", SpecName: "Holy"},
	{ClassName: "Priest", SpecName: "Shadow"},

	// Death Knight
	{ClassName: "DeathKnight", SpecName: "Blood"},
	{ClassName: "DeathKnight", SpecName: "Frost"},
	{ClassName: "DeathKnight", SpecName: "Unholy"},

	// Druid
	{ClassName: "Druid", SpecName: "Balance"},
	{ClassName: "Druid", SpecName: "Feral"},
	{ClassName: "Druid", SpecName: "Guardian"},
	{ClassName: "Druid", SpecName: "Restoration"},

	// Hunter
	{ClassName: "Hunter", SpecName: "BeastMastery"},
	{ClassName: "Hunter", SpecName: "Marksmanship"},
	{ClassName: "Hunter", SpecName: "Survival"},

	// Mage
	{ClassName: "Mage", SpecName: "Arcane"},
	{ClassName: "Mage", SpecName: "Fire"},
	{ClassName: "Mage", SpecName: "Frost"},

	// Monk
	{ClassName: "Monk", SpecName: "Brewmaster"},
	{ClassName: "Monk", SpecName: "Mistweaver"},
	{ClassName: "Monk", SpecName: "Windwalker"},

	// Paladin
	{ClassName: "Paladin", SpecName: "Holy"},
	{ClassName: "Paladin", SpecName: "Protection"},
	{ClassName: "Paladin", SpecName: "Retribution"},

	// Rogue
	{ClassName: "Rogue", SpecName: "Assassination"},
	{ClassName: "Rogue", SpecName: "Subtlety"},
	{ClassName: "Rogue", SpecName: "Outlaw"},

	// Shaman
	{ClassName: "Shaman", SpecName: "Elemental"},
	{ClassName: "Shaman", SpecName: "Enhancement"},
	{ClassName: "Shaman", SpecName: "Restoration"},

	// Warlock
	{ClassName: "Warlock", SpecName: "Affliction"},
	{ClassName: "Warlock", SpecName: "Demonology"},
	{ClassName: "Warlock", SpecName: "Destruction"},

	// Warrior
	{ClassName: "Warrior", SpecName: "Arms"},
	{ClassName: "Warrior", SpecName: "Fury"},
	{ClassName: "Warrior", SpecName: "Protection"},

	// Demon Hunter
	{ClassName: "DemonHunter", SpecName: "Havoc"},
	{ClassName: "DemonHunter", SpecName: "Vengeance"},

	// Evoker
	{ClassName: "Evoker", SpecName: "Devastation"},
	{ClassName: "Evoker", SpecName: "Preservation"},
	{ClassName: "Evoker", SpecName: "Augmentation"},
}
