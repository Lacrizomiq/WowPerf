package warcraftlogsBuildsTemporalActivities

import (
	"encoding/json"
	"testing"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePlayerBuild(t *testing.T) {
	testCases := []struct {
		name         string
		reportData   string
		playerData   string
		expectedSpec string
		shouldError  bool
	}{
		{
			name: "Discipline Priest with active talents",
			reportData: `{
                "code": "test123",
                "fightID": 1,
                "talentCodes": {
                    "Priest_Discipline_talents": "test-talent-code",
                    "Priest_Holy_talents": null,
                    "Priest_Shadow_talents": null
                }
            }`,
			playerData: `{
                "id": 1,
                "name": "TestPlayer",
                "type": "Priest",
                "specs": ["Discipline", "Holy", "Shadow"],
                "maxItemLevel": 420,
                "combatantInfo": {
                    "stats": {},
                    "gear": {},
                    "talentTree": {}
                }
            }`,
			expectedSpec: "Discipline",
			shouldError:  false,
		},
		{
			name: "Holy Priest with active talents",
			reportData: `{
                "code": "test123",
                "fightID": 1,
                "talentCodes": {
                    "Priest_Discipline_talents": null,
                    "Priest_Holy_talents": "test-talent-code",
                    "Priest_Shadow_talents": null
                }
            }`,
			playerData: `{
                "id": 1,
                "name": "TestPlayer",
                "type": "Priest",
                "specs": ["Holy", "Discipline", "Shadow"],
                "maxItemLevel": 420,
                "combatantInfo": {
                    "stats": {},
                    "gear": {},
                    "talentTree": {}
                }
            }`,
			expectedSpec: "Holy",
			shouldError:  false,
		},
		{
			name: "Shadow Priest with active talents",
			reportData: `{
                "code": "test123",
                "fightID": 1,
                "talentCodes": {
                    "Priest_Discipline_talents": null,
                    "Priest_Holy_talents": null,
                    "Priest_Shadow_talents": "test-talent-code"
                }
            }`,
			playerData: `{
                "id": 1,
                "name": "TestPlayer",
                "type": "Priest",
                "specs": ["Shadow", "Discipline", "Holy"],
                "maxItemLevel": 420,
                "combatantInfo": {
                    "stats": {},
                    "gear": {},
                    "talentTree": {}
                }
            }`,
			expectedSpec: "Shadow",
			shouldError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse test data
			var report warcraftlogsBuilds.Report
			err := json.Unmarshal([]byte(tc.reportData), &report)
			require.NoError(t, err, "Failed to parse report data")

			var playerDetails PlayerDetails
			err = json.Unmarshal([]byte(tc.playerData), &playerDetails)
			require.NoError(t, err, "Failed to parse player data")

			// Create activity instance with mock repository
			activity := &PlayerBuildsActivity{
				repository: nil, // We don't need repository for this test
			}

			// Execute the test
			build, err := activity.createPlayerBuild(&report, playerDetails)

			// Assert results
			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, build)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, build)
				assert.Equal(t, tc.expectedSpec, build.Spec)
				assert.Equal(t, playerDetails.Type, build.Class)
				assert.Equal(t, playerDetails.Name, build.PlayerName)
				assert.Equal(t, report.Code, build.ReportCode)
			}
		})
	}
}
