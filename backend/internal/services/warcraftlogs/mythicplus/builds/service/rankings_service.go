package warcraftlogsBuildsService

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	rankingsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	"gorm.io/gorm"
)

type DungeonMapping struct {
	ID          uint
	EncounterID uint
}

type RankingsService struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *rankingsRepository.RankingsRepository
	db         *gorm.DB
}

func NewRankingsService(client *warcraftlogs.WarcraftLogsClientService, repo *rankingsRepository.RankingsRepository, db *gorm.DB) *RankingsService {
	return &RankingsService{
		client:     client,
		repository: repo,
		db:         db,
	}
}

var dungeonIDs = []uint{
	12660, // Ara-Kara
	12669, // City of Threads
	60670, // Grim Batol
	62290, // Mists of Tirna Scithe
	61822, // Siege of Boralus
	12662, // The Dawnbreaker
	62286, // The Necrotic Wake
	12652, // The Stonevault
}

// StartPeriodicCollection starts the periodic collection of rankings
func (s *RankingsService) StartPeriodicCollection(ctx context.Context) {
	log.Println("Starting periodic rankings collection...")

	go func() {
		for {
			log.Println("Starting builds collection cycle...")

			var wg sync.WaitGroup
			errorChan := make(chan error, len(dungeonIDs))

			for _, dungeonID := range dungeonIDs {
				wg.Add(1)
				go func(id uint) {
					defer wg.Done()

					if err := s.FetchAndStoreRankings(ctx, id); err != nil {
						errorChan <- fmt.Errorf("error collecting builds for dungeon %d: %w", id, err)
						return
					}
					log.Printf("Successfully collected builds for dungeon %d", id)
				}(dungeonID)

				// Pause entre les requêtes pour respecter le rate limit
				time.Sleep(2 * time.Second)
			}

			// Wait for all goroutines to finish
			go func() {
				wg.Wait()
				close(errorChan)
			}()

			// Collect errors
			var errors []error
			for err := range errorChan {
				errors = append(errors, err)
				log.Printf("Error: %v", err)
			}

			if len(errors) > 0 {
				log.Printf("Errors occurred during rankings collection: %v", errors)
			} else {
				log.Println("All rankings collected successfully")
			}

			select {
			case <-ctx.Done():
				log.Println("Rankings collection interrupted")
				return
			case <-time.After(7 * 24 * time.Hour):
				// Continue to the next cycle
				continue
			}
		}
	}()
}

func (s *RankingsService) FetchAndStoreRankings(ctx context.Context, encounterId uint) error {
	log.Printf("Fetching and storing rankings for encounter %d", encounterId)

	// Retrieve the last rankings for the encounter
	lastRanking, err := s.repository.GetLastRankingForEncounter(ctx, encounterId)
	if err != nil {
		return fmt.Errorf("failed to get last rankings: %w", err)
	}

	// if the updates are too recent, we don't fetch new rankings, i skip the update
	if lastRanking != nil {
		timeSinceLastUpdate := time.Since(lastRanking.CreatedAt)
		if timeSinceLastUpdate < 7*24*time.Hour {
			log.Printf("Skipping update, last update was too recent: %s", timeSinceLastUpdate)
			return nil
		}
	} else {
		log.Printf("No previous rankings found for encounter %d", encounterId)
	}

	// retrieve new rankings
	rankings, err := s.fetchAllRankings(ctx, encounterId)
	if err != nil {
		return fmt.Errorf("failed to fetch rankings: %w", err)
	}

	log.Printf("Fetched %d rankings for encounter %d", len(rankings), encounterId)

	// Check if the rankings have changed and store them if they have
	changes := s.detectChanges(rankings, lastRanking)
	if len(changes) > 0 {
		if err := s.repository.StoreRankings(ctx, encounterId, changes); err != nil {
			return fmt.Errorf("failed to store rankings: %w", err)
		}
		log.Printf("Stored %d new rankings for encounter %d", len(changes), encounterId)
	} else {
		log.Printf("No changes detected for encounter %d", encounterId)
	}

	return nil
}

// fetchAllRankings retrieves all rankings for a given encounter
func (s *RankingsService) fetchAllRankings(ctx context.Context, encounterId uint) ([]*warcraftlogsBuilds.ClassRanking, error) {

	var dungeonId uint
	if err := s.db.Table("dungeons").
		Where("encounter_id = ?", encounterId).
		Select("id").
		Scan(&dungeonId).Error; err != nil {
		return nil, fmt.Errorf("no dungeon mapping found for encounter %d: %w", encounterId, err)
	}
	var allRankings []*warcraftlogsBuilds.ClassRanking
	page := 1
	hasMore := true

	log.Printf("Starting to fetch rankings for encounter %d", encounterId)

	for hasMore {
		response, err := s.client.MakeRequest(ctx, rankingsQueries.ClassRankingsQuery, map[string]interface{}{
			"encounterId": encounterId,
			"className":   "Priest",
			"specName":    "Discipline",
			"page":        page,
		})
		if err != nil {
			log.Printf("Error fetching rankings: %v", err)
			return nil, fmt.Errorf("failed to fetch rankings: %w", err)
		}

		rankings, hasMorePage, err := rankingsQueries.ParseRankingsResponse(response, encounterId, dungeonId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rankings: %w", err)
		}

		log.Printf("Fetched %d rankings from page %d", len(rankings), page)

		allRankings = append(allRankings, rankings...)
		hasMore = hasMorePage && page < 2 // 2 is the max page for now
		page++

		time.Sleep(500 * time.Millisecond)
	}

	return allRankings, nil
}

// detectChanges compares the new rankings with the last ones and returns the changes
func (s *RankingsService) detectChanges(newRankings []*warcraftlogsBuilds.ClassRanking, lastRanking *warcraftlogsBuilds.ClassRanking) []*warcraftlogsBuilds.ClassRanking {
	// if there is no last ranking, we consider all new rankings as changes
	if lastRanking == nil {
		log.Printf("No previous rankings found, inserting %d new rankings", len(newRankings))
		return newRankings
	}

	log.Printf("Found %d rankings to compare with last ranking", len(newRankings))
	var changes []*warcraftlogsBuilds.ClassRanking
	for _, ranking := range newRankings {
		if hasChanged(ranking, lastRanking) {
			log.Printf("Change detected for player %s: Score %.2f -> %.2f",
				ranking.PlayerName, lastRanking.Score, ranking.Score)
			changes = append(changes, ranking)
		}
	}

	log.Printf("Detected %d changes in rankings", len(changes))
	return changes
}

// hasChanged checks if a ranking has changed compared to the last one
func hasChanged(new *warcraftlogsBuilds.ClassRanking, old *warcraftlogsBuilds.ClassRanking) bool {
	return new.Score != old.Score ||
		new.HardModeLevel != old.HardModeLevel ||
		new.Medal != old.Medal
}
