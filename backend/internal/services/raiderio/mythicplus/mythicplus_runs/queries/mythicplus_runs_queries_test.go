// go test ./internal/services/raiderio/mythicplus/mythicplus_runs/queries -v
package raiderioMythicPlusRunsQueries

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	models "wowperf/internal/models/raiderio/mythicplus_runs"
)

func TestParseAPIResponse(t *testing.T) {
	tests := []struct {
		name           string
		apiResponse    map[string]interface{}
		expectedRuns   int
		expectedErrors bool
	}{
		{
			name: "Valid API response with 2 runs",
			apiResponse: map[string]interface{}{
				"rankings": []interface{}{
					map[string]interface{}{
						"rank":  1,
						"score": 493.0,
						"run": map[string]interface{}{
							"keystone_run_id": float64(17095427), // JSON numbers are float64
							"season":          "season-tww-2",
							"status":          "finished",
							"dungeon": map[string]interface{}{
								"slug": "darkflame-cleft",
								"name": "Darkflame Cleft",
							},
							"mythic_level":      float64(20),
							"clear_time_ms":     float64(1464036),
							"keystone_time_ms":  float64(1860999),
							"completed_at":      "2025-05-24T16:59:32.000Z",
							"num_chests":        float64(2),
							"time_remaining_ms": float64(396963),
							"roster": []interface{}{
								createTestRosterMember("dps", "Mage", "Arcane", "eu"),
								createTestRosterMember("dps", "Druid", "Balance", "eu"),
								createTestRosterMember("dps", "Death Knight", "Unholy", "eu"),
								createTestRosterMember("healer", "Priest", "Discipline", "eu"),
								createTestRosterMember("tank", "Demon Hunter", "Vengeance", "eu"),
							},
						},
					},
				},
			},
			expectedRuns:   1,
			expectedErrors: false,
		},
		{
			name: "Empty rankings",
			apiResponse: map[string]interface{}{
				"rankings": []interface{}{},
			},
			expectedRuns:   0,
			expectedErrors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.apiResponse)
			require.NoError(t, err)

			var response APIResponse
			err = json.Unmarshal(jsonBytes, &response)

			if tt.expectedErrors {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, response.Rankings, tt.expectedRuns)

			if tt.expectedRuns > 0 {
				firstRun := response.Rankings[0].Run
				assert.NotZero(t, firstRun.KeystoneRunID)
				assert.Equal(t, "season-tww-2", firstRun.Season)
				assert.Len(t, firstRun.Roster, 5)
			}
		})
	}
}

func TestExtractRuns(t *testing.T) {
	apiResponse := &APIResponse{
		Rankings: []struct {
			Rank  int        `json:"rank"`
			Score float64    `json:"score"`
			Run   models.Run `json:"run"`
		}{
			{
				Rank:  1,
				Score: 493.0,
				Run:   createTestRun(17095427),
			},
		},
	}

	extractedRuns := make([]*models.Run, len(apiResponse.Rankings))
	for i, ranking := range apiResponse.Rankings {
		extractedRuns[i] = &ranking.Run
	}

	assert.Len(t, extractedRuns, 1)
	assert.Equal(t, int64(17095427), extractedRuns[0].KeystoneRunID)
}

// Helper functions
func createTestRosterMember(role, className, specName, region string) map[string]interface{} {
	return map[string]interface{}{
		"role": role,
		"character": map[string]interface{}{
			"class": map[string]interface{}{
				"name": className,
			},
			"spec": map[string]interface{}{
				"name": specName,
			},
			"region": map[string]interface{}{
				"slug": region,
			},
		},
	}
}

func createTestRun(keystoneRunID int64) models.Run {
	return models.Run{
		KeystoneRunID: keystoneRunID,
		Season:        "season-tww-2",
		Status:        "finished",
		MythicLevel:   20,
		Dungeon: models.DungeonInfo{
			Slug: "darkflame-cleft",
			Name: "Darkflame Cleft",
		},
		Roster: createTestRoster(),
	}
}

func createTestRoster() []models.RosterMember {
	return []models.RosterMember{
		createTestRosterMemberData("tank", "Demon Hunter", "Vengeance", "eu"),
		createTestRosterMemberData("healer", "Priest", "Discipline", "eu"),
		createTestRosterMemberData("dps", "Mage", "Arcane", "eu"),
		createTestRosterMemberData("dps", "Warrior", "Arms", "eu"),
		createTestRosterMemberData("dps", "Death Knight", "Unholy", "eu"),
	}
}

func createTestRosterMemberData(role, className, specName, region string) models.RosterMember {
	return models.RosterMember{
		Role: role,
		Character: models.Character{
			Class: models.ClassInfo{
				Name: className,
			},
			Spec: models.SpecInfo{
				Name: specName,
			},
			Region: models.RegionInfo{
				Slug: region,
			},
		},
	}
}
