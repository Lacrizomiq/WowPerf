// package warcraftlogs/service.go
package warcraftlogs

import (
	"encoding/json"
	"fmt"
	"log"

	logsModels "wowperf/internal/models/warcraftlogs/mythicplus/ByLog"
	leaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/ByTeam"
	warcraftlogsService "wowperf/internal/services/warcraftlogs"
)

const (
	DungeonLeaderboardQuery = `
	query getDungeonLeaderboard($encounterId: Int!, $region: String!, $page: Int!) {
		worldData {
			encounter(id: $encounterId) {
				name
				fightRankings(page: $page) {
					page
					hasMorePages
					count
					rankings {
						server {
							id
							name
							region
						}
						duration
						startTime
						deaths
						tanks
						healers
						melee
						ranged
						bracketData
						affixes
						team {
							id
							name
							class
							spec
							role
						}
						medal
						score
						leaderboard
					}
				}
			}
		}
	}`

	DungeonLogsQuery = `
	query getDungeonLogs($encounterId: Int!, $metric: String!, $className: String!) {
		worldData {
			encounter(id: $encounterId) {
				name
				characterRankings(
					metric: $metric
					includeCombatantInfo: false
					className: $className
				) {
					page
					hasMorePages
					count
					rankings {
						name
						class
						spec
						amount
						hardModeLevel
						duration
						startTime
						report {
							code
							fightID
							startTime
						}
						server {
							id
							name
							region
						}
						bracketData
						faction
						affixes
						medal
						score
						leaderboard
					}
				}
			}
		}
	}`
)

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

	var result struct {
		Data struct {
			WorldData struct {
				Encounter struct {
					Name          string                               `json:"name"`
					FightRankings leaderboardModels.DungeonLeaderboard `json:"fightRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal leaderboard response: %w", err)
	}

	return &result.Data.WorldData.Encounter.FightRankings, nil
}

// GetDungeonLogs returns the leaderboard for users with logs for a given encounter, metric and class name
func (s *DungeonService) GetDungeonLogs(encounterID int, metric string, className string) (*logsModels.DungeonLogs, error) {
	log.Printf("Getting dungeon logs for encounter %d, metric %s, class %s", encounterID, metric, className)

	variables := map[string]interface{}{
		"encounterId": encounterID,
		"metric":      metric,
		"className":   className,
	}

	response, err := s.client.MakeGraphQLRequest(DungeonLogsQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get dungeon logs: %w", err)
	}

	var result struct {
		Data struct {
			WorldData struct {
				Encounter struct {
					Name              string                 `json:"name"`
					CharacterRankings logsModels.DungeonLogs `json:"characterRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs response: %w", err)
	}

	return &result.Data.WorldData.Encounter.CharacterRankings, nil
}
