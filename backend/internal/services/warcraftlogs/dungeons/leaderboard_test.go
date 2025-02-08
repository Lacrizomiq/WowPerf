package warcraftlogs

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	service "wowperf/internal/services/warcraftlogs"
)

func TestLiveLeaderboardQuery(t *testing.T) {
	// Skip in CI/automated tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load .env file
	err := godotenv.Load("../../../../.env")
	require.NoError(t, err, "Failed to load .env file")

	// Check for required environment variables
	if os.Getenv("WARCRAFTLOGS_CLIENT_ID") == "" || os.Getenv("WARCRAFTLOGS_CLIENT_SECRET") == "" {
		t.Skip("Skipping test: WARCRAFTLOGS_CLIENT_ID and WARCRAFTLOGS_CLIENT_SECRET required")
	}

	// Create the client
	client, err := service.NewWarcraftLogsClientService()
	require.NoError(t, err)

	// Test variables
	variables := map[string]interface{}{
		"encounterId":  12660, // Ara-Kara
		"page":         1,
		"className":    "Priest",
		"specName":     "Discipline",
		"serverRegion": "eu",
	}

	// Create context
	ctx := context.Background()

	// Make the request
	response, err := client.MakeRequest(ctx, DungeonLeaderboardPlayerQuery, variables)
	require.NoError(t, err)

	// Log raw response for inspection
	t.Logf("Raw API Response: %s", string(response))

	// Parse response to check ranking count
	var result struct {
		WorldData struct {
			Encounter struct {
				CharacterRankings struct {
					Rankings []interface{} `json:"rankings"`
					Count    int           `json:"count"`
				} `json:"characterRankings"`
			} `json:"encounter"`
		} `json:"worldData"`
	}

	err = json.Unmarshal(response, &result)
	require.NoError(t, err)

	t.Logf("Found %d rankings", len(result.WorldData.Encounter.CharacterRankings.Rankings))
	t.Logf("Total count: %d", result.WorldData.Encounter.CharacterRankings.Count)
}
