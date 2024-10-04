package raiderio

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"wowperf/internal/services/raiderio"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"

	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis/v8"
)

type MythicPlusBestRunHandler struct {
	Service *raiderio.RaiderIOService
}

func NewMythicPlusBestRunHandler(service *raiderio.RaiderIOService) *MythicPlusBestRunHandler {
	return &MythicPlusBestRunHandler{
		Service: service,
	}
}

func (h *MythicPlusBestRunHandler) GetMythicPlusBestRuns(c *gin.Context) {
	season := c.Query("season")
	region := c.Query("region")
	dungeon := c.Query("dungeon")
	pageStr := c.Query("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	if season == "" || region == "" || dungeon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	cacheKey := fmt.Sprintf("mythic_plus_best_runs_%s_%s_%s_%d", season, region, dungeon, page)

	var bestRuns map[string]interface{}
	err = cache.Get(cacheKey, &bestRuns)
	if err == nil {
		c.JSON(http.StatusOK, bestRuns)
		return
	}

	if err != redis.Nil {
		log.Println("Error getting from cache", err)
	}

	bestRuns, err = raiderioMythicPlus.GetMythicPlusBestRuns(h.Service, season, region, dungeon, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = cache.Set(cacheKey, bestRuns, 12*time.Hour)
	if err != nil {
		log.Println("Error setting cache:", err)
	}

	c.JSON(http.StatusOK, bestRuns)
}

func StartMythicPlusBestRunsCacheUpdater(service *raiderio.RaiderIOService) {
	updateFunc := func() error {
		seasons := []string{"season-tww-1"}
		regions := []string{"world", "us", "eu", "tw", "kr", "cn"}
		dungeons := []string{"all", "arakara-city-of-echoes", "city-of-threads", "grim-batol", "mists-of-tirna-scithe", "siege-of-boralus", "the-dawnbreaker", "the-necrotic-wake", "the-stonevault"}
		pages := []int{0, 1, 2, 3, 4, 5}

		for _, season := range seasons {
			for _, region := range regions {
				for _, dungeon := range dungeons {
					for _, page := range pages {
						bestRuns, err := raiderioMythicPlus.GetMythicPlusBestRuns(service, season, region, dungeon, page)
						if err != nil {
							return fmt.Errorf("error updating cache for %s %s %s %d: %w", season, region, dungeon, page, err)
						}
						cacheKey := fmt.Sprintf("mythic_plus_best_runs_%s_%s_%s_%d", season, region, dungeon, page)
						cache.Set(cacheKey, bestRuns, 3*time.Hour)
					}
				}
			}
		}
		return nil
	}

	cache.StartPeriodicUpdate("mythic_plus_best_runs", updateFunc, 3*time.Hour)
}
