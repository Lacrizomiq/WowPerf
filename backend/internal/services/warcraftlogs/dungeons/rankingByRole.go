package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	service "wowperf/internal/services/warcraftlogs"
)

func determineRole(class, spec string) string {
	tanks := map[string][]string{
		"Warrior":     {"Protection"},
		"Paladin":     {"Protection"},
		"DeathKnight": {"Blood"},
		"DemonHunter": {"Vengeance"},
		"Druid":       {"Guardian"},
		"Monk":        {"Brewmaster"},
	}

	healers := map[string][]string{
		"Priest":  {"Holy", "Discipline"},
		"Paladin": {"Holy"},
		"Druid":   {"Restoration"},
		"Shaman":  {"Restoration"},
		"Monk":    {"Mistweaver"},
		"Evoker":  {"Preservation"},
	}

	if specs, ok := tanks[class]; ok {
		for _, tankSpec := range specs {
			if spec == tankSpec {
				return "Tank"
			}
		}
	}

	if specs, ok := healers[class]; ok {
		for _, healSpec := range specs {
			if spec == healSpec {
				return "Healer"
			}
		}
	}

	return "DPS"
}

// GetGlobalRankings retrieve the global rankings for the given dungeon IDs and pages per dungeon
func GetGlobalRankings(s *service.WarcraftLogsClientService, ctx context.Context, dungeonIDs []int, pagesPerDungeon int) (*GlobalRankings, error) {
	log.Printf("Starting global rankings collection for %d dungeons, %d pages each", len(dungeonIDs), pagesPerDungeon)

	var wg sync.WaitGroup
	sem := make(chan struct{}, requestsPerSecond)

	type playerKey struct {
		name      string
		dungeonID int
	}
	bestScores := make(map[playerKey]*PlayerScore)
	var scoresMutex sync.Mutex
	errorsChan := make(chan error, len(dungeonIDs)*pagesPerDungeon)

	for _, dungeonID := range dungeonIDs {
		for page := 1; page <= pagesPerDungeon; page++ {
			wg.Add(1)
			go func(dID, p int) {
				defer wg.Done()

				sem <- struct{}{}
				defer func() { <-sem }()

				if err := s.RateLimiter.Wait(ctx); err != nil {
					errorsChan <- fmt.Errorf("rate limiter error: %w", err)
					return
				}

				reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
				defer cancel()

				select {
				case <-reqCtx.Done():
					errorsChan <- reqCtx.Err()
					return
				default:
					dungeonData, err := GetDungeonLeaderboardByPlayer(s, dID, p)
					if err != nil {
						errorsChan <- fmt.Errorf("failed to get dungeon leaderboard for dungeon %d, page %d: %w", dID, p, err)
						return
					}

					scoresMutex.Lock()
					for _, ranking := range dungeonData.Rankings {
						key := playerKey{
							name:      ranking.Name,
							dungeonID: dID,
						}

						role := determineRole(ranking.Class, ranking.Spec)

						// Update or create the player score
						playerScore := &PlayerScore{
							Name:       ranking.Name,
							Class:      ranking.Class,
							Spec:       ranking.Spec,
							Role:       role,
							TotalScore: ranking.Score,
							Amount:     ranking.Amount,
							Guild: Guild{
								ID:      ranking.Guild.ID,
								Name:    ranking.Guild.Name,
								Faction: ranking.Guild.Faction,
							},
							Server: Server{
								ID:     ranking.Server.ID,
								Name:   ranking.Server.Name,
								Region: ranking.Server.Region,
							},
							Faction: ranking.Faction,
							Runs: []Run{{
								DungeonID:     dID,
								Score:         ranking.Score,
								Duration:      ranking.Duration,
								StartTime:     ranking.StartTime,
								HardModeLevel: ranking.HardModeLevel,
								BracketData:   ranking.BracketData,
								Medal:         ranking.Medal,
								Affixes:       ranking.Affixes,
								Report: Report{
									Code:      ranking.Report.Code,
									FightID:   ranking.Report.FightID,
									StartTime: ranking.Report.StartTime,
								},
							}},
						}

						if existing, exists := bestScores[key]; exists {
							if ranking.Score > existing.TotalScore {
								bestScores[key] = playerScore
							}
						} else {
							bestScores[key] = playerScore
						}
					}
					scoresMutex.Unlock()

					log.Printf("Processed dungeon %d page %d: %d rankings", dID, p, len(dungeonData.Rankings))
				}
			}(dungeonID, page)
		}
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(errorsChan)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		// Continue processing
	}

	// Process errors
	var errors []error
	for err := range errorsChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("encountered %d errors during ranking collection", len(errors))
	}

	// Build the final result
	result := &GlobalRankings{
		Tanks:   RoleRankings{Players: make([]PlayerScore, 0)},
		Healers: RoleRankings{Players: make([]PlayerScore, 0)},
		DPS:     RoleRankings{Players: make([]PlayerScore, 0)},
	}

	// Aggregate scores by player
	playerTotalScores := make(map[string]*PlayerScore)
	for _, score := range bestScores {
		if existing, exists := playerTotalScores[score.Name]; exists {
			existing.TotalScore += score.TotalScore
			existing.Runs = append(existing.Runs, score.Runs...)
		} else {
			playerCopy := *score
			playerTotalScores[score.Name] = &playerCopy
		}
	}

	// Distribute scores by role
	for _, player := range playerTotalScores {
		switch player.Role {
		case "Tank":
			result.Tanks.Players = append(result.Tanks.Players, *player)
		case "Healer":
			result.Healers.Players = append(result.Healers.Players, *player)
		case "DPS":
			result.DPS.Players = append(result.DPS.Players, *player)
		}
	}

	// Final sorting
	sort.Slice(result.Tanks.Players, func(i, j int) bool {
		return result.Tanks.Players[i].TotalScore > result.Tanks.Players[j].TotalScore
	})
	sort.Slice(result.Healers.Players, func(i, j int) bool {
		return result.Healers.Players[i].TotalScore > result.Healers.Players[j].TotalScore
	})
	sort.Slice(result.DPS.Players, func(i, j int) bool {
		return result.DPS.Players[i].TotalScore > result.DPS.Players[j].TotalScore
	})

	// Update counters
	result.Tanks.Count = len(result.Tanks.Players)
	result.Healers.Count = len(result.Healers.Players)
	result.DPS.Count = len(result.DPS.Players)

	log.Printf("Rankings collection completed: Tanks: %d, Healers: %d, DPS: %d",
		result.Tanks.Count,
		result.Healers.Count,
		result.DPS.Count)

	return result, nil
}
