// package warcraftlogs/dungeons/rankingByRole.go
package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	leaderboardModels "wowperf/internal/models/warcraftlogs/mythicplus/ByTeam"
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
}

// Wrapper structure for the ranking with the dungeon ID
type RankingWithDungeon struct {
	*leaderboardModels.Ranking
	DungeonID int
}

func NewRankingsService(dungeonService *DungeonService) *RankingsService {
	return &RankingsService{
		dungeonService: dungeonService,
	}
}

// Get the global rankings of each role for a list of dungeons
func (s *RankingsService) GetGlobalRankings(ctx context.Context, dungeonIDs []int, pagesPerDungeon int) (*GlobalRankings, error) {
	rankingsChan := make(chan RankingWithDungeon, len(dungeonIDs)*pagesPerDungeon*100)
	errorsChan := make(chan error, len(dungeonIDs)*pagesPerDungeon)

	var wg sync.WaitGroup
	PlayerScores := make(map[int]*PlayerScore)
	var playersMutex sync.Mutex

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
					rankings, err := s.dungeonService.GetDungeonLeaderboard(dID, p)
					if err != nil {
						errorsChan <- fmt.Errorf("failed to get dungeon leaderboard for dungeon %d, page %d: %w", dID, p, err)
						return
					}

					for _, ranking := range rankings.Rankings {
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

	for wrappedRanking := range rankingsChan {
		for _, member := range wrappedRanking.Ranking.Team {
			playersMutex.Lock()
			if player, exists := PlayerScores[member.ID]; exists {
				player.Runs = append(player.Runs, Run{
					DungeonID:  wrappedRanking.DungeonID,
					Score:      wrappedRanking.Ranking.Score,
					Time:       wrappedRanking.Ranking.Duration,
					AffixLevel: wrappedRanking.Ranking.BracketData,
					Medal:      wrappedRanking.Ranking.Medal,
				})
				player.TotalScore += wrappedRanking.Ranking.Score
			} else {
				PlayerScores[member.ID] = &PlayerScore{
					ID:    member.ID,
					Name:  member.Name,
					Class: member.Class,
					Spec:  member.Spec,
					Role:  member.Role,
					Runs: []Run{{
						DungeonID:  wrappedRanking.DungeonID,
						Score:      wrappedRanking.Ranking.Score,
						Time:       wrappedRanking.Ranking.Duration,
						AffixLevel: wrappedRanking.Ranking.BracketData,
						Medal:      wrappedRanking.Ranking.Medal,
					}},
					TotalScore: wrappedRanking.Ranking.Score,
				}
			}
			playersMutex.Unlock()
		}
	}

	// Check for errors
	for err := range errorsChan {
		if err != nil {
			log.Printf("Error getting dungeon leaderboard: %v", err)
		}
	}

	// Prepare the result
	result := &GlobalRankings{}
	tanks := make([]PlayerScore, 0)
	healers := make([]PlayerScore, 0)
	dps := make([]PlayerScore, 0)

	for _, player := range PlayerScores {
		switch player.Role {
		case "Tank":
			tanks = append(tanks, *player)
		case "Healer":
			healers = append(healers, *player)
		case "DPS":
			dps = append(dps, *player)
		}
	}

	// Sort by total score
	sort.SliceStable(tanks, func(i, j int) bool {
		return tanks[i].TotalScore > tanks[j].TotalScore
	})
	sort.SliceStable(healers, func(i, j int) bool {
		return healers[i].TotalScore > healers[j].TotalScore
	})
	sort.SliceStable(dps, func(i, j int) bool {
		return dps[i].TotalScore > dps[j].TotalScore
	})

	// Assign the results
	result.Tanks.Players = tanks
	result.Tanks.Count = len(tanks)
	result.Healers.Players = healers
	result.Healers.Count = len(healers)
	result.DPS.Players = dps
	result.DPS.Count = len(dps)

	return result, nil
}
