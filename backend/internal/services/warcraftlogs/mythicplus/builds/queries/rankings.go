package warcraftlogsBuildsQueries

import (
	"encoding/json"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"github.com/lib/pq"
)

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

func ParseRankingsResponse(response []byte, encounterId uint) ([]*warcraftlogsBuilds.ClassRanking, bool, error) {
	log.Printf("Parsing response data")

	var result struct {
		Data struct {
			WorldData struct {
				Encounter struct {
					Name              string
					CharacterRankings struct {
						HasMorePages bool              `json:"hasMorePages"`
						Rankings     []json.RawMessage `json:"rankings"`
					} `json:"characterRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal rankings response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, false, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	log.Printf("Found %d rankings in response", len(result.Data.WorldData.Encounter.CharacterRankings.Rankings))

	rankings := make([]*warcraftlogsBuilds.ClassRanking, 0, len(result.Data.WorldData.Encounter.CharacterRankings.Rankings))

	for _, rankingData := range result.Data.WorldData.Encounter.CharacterRankings.Rankings {
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
			BracketData int     `json:"bracketData"`
			Faction     int     `json:"faction"`
			Affixes     []int   `json:"affixes"`
			Medal       string  `json:"medal"`
			Score       float64 `json:"score"`
		}

		if err := json.Unmarshal(rankingData, &rankingResponse); err != nil {
			return nil, false, fmt.Errorf("failed to unmarshal ranking: %w", err)
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

		rankings = append(rankings, ranking)
	}

	return rankings, result.Data.WorldData.Encounter.CharacterRankings.HasMorePages, nil
}

func intSliceToInt64Array(ints []int) pq.Int64Array {
	int64s := make([]int64, len(ints))
	for i, v := range ints {
		int64s[i] = int64(v)
	}
	return int64s
}
