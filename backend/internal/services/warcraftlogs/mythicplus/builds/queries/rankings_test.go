package warcraftlogsBuildsQueries

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRankingsResponse(t *testing.T) {
	t.Run("should return only top 20 rankings when more than 20 are available", func(t *testing.T) {
		// Create a response with 100 rankings
		response := createMockResponse(100)

		// Parse the response
		rankings, err := ParseRankingsResponse(response, 12345)

		// Assert results
		require.NoError(t, err)
		assert.Len(t, rankings, 20, "Should return exactly 20 rankings")

		// Verify rankings are sorted by score (highest first)
		for i := 0; i < len(rankings)-1; i++ {
			assert.GreaterOrEqual(t,
				rankings[i].Score,
				rankings[i+1].Score,
				"Rankings should be ordered by score")
		}
	})

	t.Run("should return all rankings when less than 20 are available", func(t *testing.T) {
		// Create a response with 10 rankings
		response := createMockResponse(10)

		// Parse the response
		rankings, err := ParseRankingsResponse(response, 12345)

		// Assert results
		require.NoError(t, err)
		assert.Len(t, rankings, 10, "Should return all 10 rankings")
	})

	t.Run("should handle empty rankings array", func(t *testing.T) {
		response := createMockResponse(0)

		rankings, err := ParseRankingsResponse(response, 12345)

		require.NoError(t, err)
		assert.Empty(t, rankings, "Should return empty slice for no rankings")
	})

	t.Run("should handle invalid JSON", func(t *testing.T) {
		invalidJSON := []byte(`{"invalid": json}`)

		rankings, err := ParseRankingsResponse(invalidJSON, 12345)

		assert.Error(t, err)
		assert.Nil(t, rankings)
	})
}

// createMockResponse creates a mock API response with the specified number of rankings
func createMockResponse(numRankings int) []byte {
	rankings := make([]map[string]interface{}, numRankings)

	for i := 0; i < numRankings; i++ {
		rankings[i] = map[string]interface{}{
			"name":          fmt.Sprintf("Player%d", i),
			"class":         "Priest",
			"spec":          "Shadow",
			"amount":        float64(1000 - i), // Decreasing amounts
			"hardModeLevel": 15,
			"duration":      int64(1200),
			"startTime":     int64(1600000000),
			"report": map[string]interface{}{
				"code":      fmt.Sprintf("ABC%d", i),
				"fightID":   i,
				"startTime": int64(1600000000),
			},
			"server": map[string]interface{}{
				"id":     1,
				"name":   "Server1",
				"region": "EU",
			},
			"guild": map[string]interface{}{
				"id":      1,
				"name":    "Guild1",
				"faction": 1,
			},
			"faction": 1,
			"affixes": []int{1, 2, 3},
			"medal":   "gold",
			"score":   float64(100 - i), // Decreasing scores
		}
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"worldData": map[string]interface{}{
				"encounter": map[string]interface{}{
					"name": "Test Encounter",
					"characterRankings": map[string]interface{}{
						"rankings": rankings,
					},
				},
			},
		},
	}

	data, _ := json.Marshal(response)
	return data
}
