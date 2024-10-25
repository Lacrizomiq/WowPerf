// package warcraftlogs/dungeons/rankingByRole.go
package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	leaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/team"

	"gorm.io/gorm"
)

// Base structures
type Run struct {
	DungeonID  int     `json:"dungeonId"`
	Score      float64 `json:"score"`
	Time       int64   `json:"time"`
	AffixLevel int     `json:"affixLevel"`
	Medal      string  `json:"medal"`
}

type PlayerScore struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Class      string  `json:"class"`
	Spec       string  `json:"spec"`
	Role       string  `json:"role"`
	TotalScore float64 `json:"totalScore"`
	Runs       []Run   `json:"runs"`
}

type RoleRankings struct {
	Players []PlayerScore `json:"players"`
	Count   int           `json:"count"`
}

type GlobalRankings struct {
	Tanks   RoleRankings `json:"tanks"`
	Healers RoleRankings `json:"healers"`
	DPS     RoleRankings `json:"dps"`
}

// Rankings service structure
type RankingsService struct {
	dungeonService *DungeonService
	db             *gorm.DB
}

// Wrapper structure for the ranking with the dungeon ID
type RankingWithDungeon struct {
	*leaderboardModels.Ranking
	DungeonID int
}

func NewRankingsService(dungeonService *DungeonService, db *gorm.DB) *RankingsService {
	return &RankingsService{
		dungeonService: dungeonService,
		db:             db,
	}
}

// Get the global rankings of each role for a list of dungeons
func (s *RankingsService) GetGlobalRankings(ctx context.Context, dungeonIDs []int, pagesPerDungeon int) (*GlobalRankings, error) {
	rankingsChan := make(chan RankingWithDungeon, len(dungeonIDs)*pagesPerDungeon*100)
	errorsChan := make(chan error, len(dungeonIDs)*pagesPerDungeon)

	var wg sync.WaitGroup

	// Map to store the best score by player and by dungeon
	type playerDungeonKey struct {
		playerID  int
		dungeonID int
	}
	bestScores := make(map[playerDungeonKey]*PlayerScore)
	var scoresMutex sync.Mutex

	// Launch the goroutines
	for _, dungeonID := range dungeonIDs {
		for page := 1; page <= pagesPerDungeon; page++ {
			wg.Add(1)
			go func(dID, p int) {
				defer wg.Done()

				select {
				case <-ctx.Done():
					errorsChan <- ctx.Err()
					return
				default:
					dungeonData, err := s.dungeonService.GetDungeonLeaderboard(dID, p)
					if err != nil {
						errorsChan <- fmt.Errorf("failed to get dungeon leaderboard for dungeon %d, page %d: %w", dID, p, err)
						return
					}

					for _, ranking := range dungeonData.Rankings {
						rankingsChan <- RankingWithDungeon{
							Ranking:   &ranking,
							DungeonID: dID,
						}
					}
				}
			}(dungeonID, page)
		}
	}

	go func() {
		wg.Wait()
		close(rankingsChan)
		close(errorsChan)
	}()

	// Process the rankings
	for wrappedRanking := range rankingsChan {
		if wrappedRanking.Ranking == nil || len(wrappedRanking.Ranking.Team) == 0 {
			continue
		}

		for _, member := range wrappedRanking.Ranking.Team {
			scoresMutex.Lock()
			key := playerDungeonKey{
				playerID:  member.ID,
				dungeonID: wrappedRanking.DungeonID,
			}

			run := Run{
				DungeonID:  wrappedRanking.DungeonID,
				Score:      wrappedRanking.Ranking.Score,
				Time:       wrappedRanking.Ranking.Duration,
				AffixLevel: wrappedRanking.Ranking.BracketData,
				Medal:      wrappedRanking.Ranking.Medal,
			}

			// Check if we already have a score for this player in this dungeon
			if existingScore, exists := bestScores[key]; exists {
				if wrappedRanking.Ranking.Score > existingScore.Runs[0].Score {
					// Update with the best score
					existingScore.Runs = []Run{run}
					existingScore.TotalScore = wrappedRanking.Ranking.Score
				}
			} else {
				// New player for this dungeon
				bestScores[key] = &PlayerScore{
					ID:         member.ID,
					Name:       member.Name,
					Class:      member.Class,
					Spec:       member.Spec,
					Role:       member.Role,
					TotalScore: wrappedRanking.Ranking.Score,
					Runs:       []Run{run},
				}
			}
			scoresMutex.Unlock()
		}
	}

	// Check for errors
	for err := range errorsChan {
		if err != nil {
			log.Printf("Error getting dungeon leaderboard: %v", err)
		}
	}

	// Convert map to slices grouped by role
	result := &GlobalRankings{
		Tanks:   RoleRankings{Players: make([]PlayerScore, 0)},
		Healers: RoleRankings{Players: make([]PlayerScore, 0)},
		DPS:     RoleRankings{Players: make([]PlayerScore, 0)},
	}

	// Group players by their ID to calculer leur score total
	playerScores := make(map[int]*PlayerScore)
	for _, score := range bestScores {
		if player, exists := playerScores[score.ID]; exists {
			// Add the run to the existing player
			player.Runs = append(player.Runs, score.Runs...)
			player.TotalScore += score.TotalScore
		} else {
			// Create a copy of the score for the new player
			newPlayer := *score
			playerScores[score.ID] = &newPlayer
		}
	}

	// Convert the map to slices by role
	for _, player := range playerScores {
		playerCopy := *player
		switch player.Role {
		case "Tank":
			result.Tanks.Players = append(result.Tanks.Players, playerCopy)
		case "Healer":
			result.Healers.Players = append(result.Healers.Players, playerCopy)
		case "DPS":
			result.DPS.Players = append(result.DPS.Players, playerCopy)
		}
	}

	// Sort each role's players by total score
	sort.Slice(result.Tanks.Players, func(i, j int) bool {
		return result.Tanks.Players[i].TotalScore > result.Tanks.Players[j].TotalScore
	})
	sort.Slice(result.Healers.Players, func(i, j int) bool {
		return result.Healers.Players[i].TotalScore > result.Healers.Players[j].TotalScore
	})
	sort.Slice(result.DPS.Players, func(i, j int) bool {
		return result.DPS.Players[i].TotalScore > result.DPS.Players[j].TotalScore
	})

	// Update counts
	result.Tanks.Count = len(result.Tanks.Players)
	result.Healers.Count = len(result.Healers.Players)
	result.DPS.Count = len(result.DPS.Players)

	return result, nil
}
