package builds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
)

const ClassRankingsQuery = `
query getClassRankings($encounterId: Int!, $className: String!, $specName: String!) {
    worldData {
        encounter(id: $encounterId) {
            name
            characterRankings(
                leaderboard: LogsOnly
                className: $className
                specName: $specName
            )
        }
    }
}
`

type BuildsService struct {
	client *warcraftlogs.WarcraftLogsClientService
}

func NewBuildsService(client *warcraftlogs.WarcraftLogsClientService) *BuildsService {
	return &BuildsService{
		client: client,
	}
}

func (s *BuildsService) FetchClassRankings(ctx context.Context, encounterID int) error {
	variables := map[string]interface{}{
		"encounterId": encounterID,
		"className":   "Priest",
		"specName":    "Discipline",
	}

	response, err := s.client.MakeRequest(ctx, ClassRankingsQuery, variables)
	if err != nil {
		return fmt.Errorf("failed to fetch class rankings: %w", err)
	}

	// Parse response
	var result struct {
		Data struct {
			WorldData struct {
				Encounter struct {
					Name              string
					CharacterRankings struct {
						Rankings []struct {
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
							} `json:"guild,omitempty"`
							Affixes []int   `json:"affixes"`
							Medal   string  `json:"medal"`
							Score   float64 `json:"score"`
						} `json:"rankings"`
					} `json:"characterRankings"`
				} `json:"encounter"`
			} `json:"worldData"`
		} `json:"data"`
	}

	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to unmarshal class rankings response: %w", err)
	}

	// Convert to struct
	rankings := make([]*warcraftlogsBuilds.ClassRanking, 0)
	for _, r := range result.Data.WorldData.Encounter.CharacterRankings.Rankings {
		ranking := &warcraftlogsBuilds.ClassRanking{
			PlayerName:    r.Name,
			Class:         r.Class,
			Spec:          r.Spec,
			EncounterID:   uint(encounterID),
			Amount:        r.Amount,
			HardModeLevel: r.HardModeLevel,
			Duration:      r.Duration,
			StartTime:     r.StartTime,
			ReportCode:    r.Report.Code,
			ReportFightID: r.Report.FightID,
			ServerID:      r.Server.ID,
			ServerName:    r.Server.Name,
			ServerRegion:  r.Server.Region,
			GuildID:       &r.Guild.ID,
			GuildName:     &r.Guild.Name,
			GuildFaction:  &r.Guild.Faction,
			Medal:         r.Medal,
			Score:         r.Score,
			Affixes:       r.Affixes,
		}
		rankings = append(rankings, ranking)
	}

	// TODO: Save to database
	log.Printf("Fetched %d class rankings for encounter %d", len(rankings), encounterID)

	return nil
}
