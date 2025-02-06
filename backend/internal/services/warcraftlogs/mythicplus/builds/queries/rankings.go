package warcraftlogsBuildsQueries

import (
	"encoding/json"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"github.com/lib/pq"
)

// ClassRankingsQuery defines the GraphQL query to fetch rankings
// Note : The API return 100 rankings but i filter after to only get the first 20 char
const ClassRankingsQuery = `
query getClassRankings($encounterId: Int!, $className: String!, $specName: String!, $page: Int!) {
    worldData {
        encounter(id: $encounterId) {
            name
            characterRankings(
                leaderboard: LogsOnly
                className: $className
                specName: $specName
                page: $page
            )
        }
    }
}`

// ParseRankingsResponse processes the API response and returns only the top 20 rankings
// The rankings are already sorted by score from the API
func ParseRankingsResponse(response []byte, encounterId uint) ([]*warcraftlogsBuilds.ClassRanking, error) {
	log.Printf("Parsing response data for encounter %d", encounterId)

	// Define the response structure
	var result struct {
		WorldData struct {
			Encounter struct {
				Name              string
				CharacterRankings struct {
					Page         int               `json:"page"`
					HasMorePages bool              `json:"hasMorePages"`
					Count        int               `json:"count"`
					Rankings     []json.RawMessage `json:"rankings"`
				} `json:"characterRankings"`
			} `json:"encounter"`
		} `json:"worldData"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	// Unmarshal the response
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rankings response: %w", err)
	}

	log.Printf("Encounter name: %s", result.WorldData.Encounter.Name)
	log.Printf("Count: %d", result.WorldData.Encounter.CharacterRankings.Count)

	// Check for API errors
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Get all rankings first
	rankings := result.WorldData.Encounter.CharacterRankings.Rankings
	log.Printf("Found %d rankings in response", len(rankings))

	// Determine how many rankings to process (minimum of 20 or available rankings)
	rankingsToProcess := 20
	if len(rankings) < 20 {
		rankingsToProcess = len(rankings)
	}

	// Process only the top rankings
	processedRankings := make([]*warcraftlogsBuilds.ClassRanking, 0, rankingsToProcess)

	for i := 0; i < rankingsToProcess; i++ {
		var rankingResponse struct {
			Name          string  `json:"name"`
			Class         string  `json:"class"`
			Spec          string  `json:"spec"`
			Amount        float64 `json:"amount"`
			HardModeLevel int     `json:"hardModeLevel"`
			Duration      int64   `json:"duration"`
			StartTime     int64   `json:"startTime"`
			Report        struct {
				Code      string `json:"code"`
				FightID   int    `json:"fightID"`
				StartTime int64  `json:"startTime"`
			} `json:"report"`
			Server struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Region string `json:"region"`
			} `json:"server"`
			Guild struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Faction int    `json:"faction"`
			} `json:"guild"`
			Faction int     `json:"faction"`
			Affixes []int   `json:"affixes"`
			Medal   string  `json:"medal"`
			Score   float64 `json:"score"`
		}

		if err := json.Unmarshal(rankings[i], &rankingResponse); err != nil {
			log.Printf("Warning: Failed to unmarshal ranking at position %d: %v", i, err)
			continue
		}

		ranking := &warcraftlogsBuilds.ClassRanking{
			PlayerName:    rankingResponse.Name,
			Class:         rankingResponse.Class,
			Spec:          rankingResponse.Spec,
			EncounterID:   encounterId,
			Amount:        rankingResponse.Amount,
			HardModeLevel: rankingResponse.HardModeLevel,
			Duration:      rankingResponse.Duration,
			StartTime:     rankingResponse.StartTime,
			ReportCode:    rankingResponse.Report.Code,
			ReportFightID: rankingResponse.Report.FightID,
			ServerID:      rankingResponse.Server.ID,
			ServerName:    rankingResponse.Server.Name,
			ServerRegion:  rankingResponse.Server.Region,
			GuildID:       &rankingResponse.Guild.ID,
			GuildName:     &rankingResponse.Guild.Name,
			GuildFaction:  &rankingResponse.Guild.Faction,
			Medal:         rankingResponse.Medal,
			Score:         rankingResponse.Score,
			Faction:       rankingResponse.Faction,
			Affixes:       intSliceToInt64Array(rankingResponse.Affixes),
		}

		processedRankings = append(processedRankings, ranking)
	}

	log.Printf("Successfully processed top %d rankings for encounter %d", len(processedRankings), encounterId)

	return processedRankings, nil
}

func intSliceToInt64Array(ints []int) pq.Int64Array {
	int64s := make([]int64, len(ints))
	for i, v := range ints {
		int64s[i] = int64(v)
	}
	return int64s
}
