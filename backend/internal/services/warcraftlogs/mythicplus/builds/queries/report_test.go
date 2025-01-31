package warcraftlogsBuildsQueries

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"wowperf/internal/services/warcraftlogs"
)

func TestLiveReportQuery(t *testing.T) {
	// Skip in CI/automated tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Try to load .env from multiple possible locations
	err := godotenv.Load()
	if err != nil {
		// Try loading from backend directory
		err = godotenv.Load("../../../../../../.env")
		require.NoError(t, err, "Failed to load .env file")
	}

	// Check for required environment variables
	if os.Getenv("WARCRAFTLOGS_CLIENT_ID") == "" || os.Getenv("WARCRAFTLOGS_CLIENT_SECRET") == "" {
		t.Skip("Skipping test: WARCRAFTLOGS_CLIENT_ID and WARCRAFTLOGS_CLIENT_SECRET environment variables are required")
	}

	// Create the client
	client, err := warcraftlogs.NewClient()
	require.NoError(t, err)

	// Test variables
	reportCode := "HZMACFrw7P198yLb" // Un code de rapport réel
	fightID := 95                    // ID du combat
	encounterID := uint(12660)       // ID de la rencontre

	variables := map[string]interface{}{
		"code":        reportCode,
		"fightID":     fightID,
		"encounterID": encounterID,
	}

	// Make the request
	response, err := client.MakeGraphQLRequest(GetReportTableQuery, variables)
	require.NoError(t, err)

	// Log the raw response for inspection
	t.Logf("Raw API Response: %s", string(response))

	// Print the variables being used
	t.Logf("Variables used: %+v", variables)

	// Try to parse
	report, talentsQuery, err := ParseReportDetailsResponse(response, reportCode, fightID, encounterID)
	require.NoError(t, err)

	// Log the details
	t.Logf("Report details:")
	t.Logf("- Code: %s", report.Code)
	t.Logf("- Fight ID: %d", report.FightID)
	t.Logf("- Encounter ID: %d", report.EncounterID)
	t.Logf("- Total Time: %d", report.TotalTime)
	t.Logf("- Keystone Level: %d", report.KeystoneLevel)
	t.Logf("- Keystone Time: %d", report.KeystoneTime)
	t.Logf("- Generated Talents Query: %s", talentsQuery)

	// If we got talents query, test it as well
	if talentsQuery != "" {
		talentsResponse, err := client.MakeGraphQLRequest(talentsQuery, nil)
		require.NoError(t, err)

		t.Logf("Raw Talents Response: %s", string(talentsResponse))

		talentCodes, err := ParseReportTalentsResponse(talentsResponse)
		require.NoError(t, err)

		t.Logf("Found %d talent codes", len(talentCodes))
		for spec, code := range talentCodes {
			t.Logf("- %s: %s", spec, code)
		}
	}

	// Verify essential fields are present
	assert.NotEmpty(t, report.Code)
	assert.NotZero(t, report.FightID)
	assert.NotZero(t, report.EncounterID)
	assert.NotZero(t, report.TotalTime)
	assert.NotZero(t, report.KeystoneLevel)
	assert.NotZero(t, report.KeystoneTime)
	assert.NotNil(t, report.Composition)
	assert.NotNil(t, report.PlayerDetailsDps)
	assert.NotNil(t, report.PlayerDetailsHealers)
	assert.NotNil(t, report.PlayerDetailsTanks)
}
