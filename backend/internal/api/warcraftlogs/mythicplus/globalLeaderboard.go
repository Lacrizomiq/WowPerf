// internal/api/warcraftlogs/mythicplus/globalLeaderboard.go
package warcraftlogs

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	dungeons "wowperf/internal/services/warcraftlogs/dungeons"
	"wowperf/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type GlobalLeaderboardHandler struct {
	rankingsService *dungeons.RankingsService
}

func NewGlobalLeaderboardHandler(rankingsService *dungeons.RankingsService) *GlobalLeaderboardHandler {
	return &GlobalLeaderboardHandler{rankingsService: rankingsService}
}

func (h *GlobalLeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	cacheKey := fmt.Sprintf("warcraftlogs:global:limit:%d", limit)

	var leaderboard interface{}
	err = cache.Get(cacheKey, &leaderboard)
	if err == nil {
		c.JSON(http.StatusOK, leaderboard)
		return
	}

	if err != redis.Nil {
		log.Printf("Error getting global leaderboard from cache: %v", err)
	}

	leaderboard, err = h.rankingsService.GetGlobalLeaderboard(c.Request.Context(), limit)
	if err != nil {
		log.Printf("Error getting global leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get global leaderboard"})
		return
	}

	if err := cache.Set(cacheKey, leaderboard, 8*time.Hour); err != nil {
		log.Printf("Error setting global leaderboard in cache: %v", err)
	}

	c.JSON(http.StatusOK, leaderboard)
}

func (h *GlobalLeaderboardHandler) GetRoleLeaderboard(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	cacheKey := fmt.Sprintf("warcraftlogs:role:%s:limit:%d", role, limit)

	var leaderboard interface{}
	err = cache.Get(cacheKey, &leaderboard)
	if err == nil {
		c.JSON(http.StatusOK, leaderboard)
		return
	}

	if err != redis.Nil {
		log.Printf("Error getting role leaderboard from cache: %v", err)
	}

	leaderboard, err = h.rankingsService.GetGlobalLeaderboardByRole(c.Request.Context(), role, limit)
	if err != nil {
		log.Printf("Error getting role leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role leaderboard"})
		return
	}

	if err := cache.Set(cacheKey, leaderboard, 8*time.Hour); err != nil {
		log.Printf("Error setting role leaderboard in cache: %v", err)
	}

	c.JSON(http.StatusOK, leaderboard)
}

func (h *GlobalLeaderboardHandler) GetClassLeaderboard(c *gin.Context) {
	class := c.Query("class")
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	cacheKey := fmt.Sprintf("warcraftlogs:class:%s:limit:%d", class, limit)

	var leaderboard interface{}
	err = cache.Get(cacheKey, &leaderboard)
	if err == nil {
		c.JSON(http.StatusOK, leaderboard)
		return
	}

	if err != redis.Nil {
		log.Printf("Error getting class leaderboard from cache: %v", err)
	}

	leaderboard, err = h.rankingsService.GetGlobalLeaderboardByClass(c.Request.Context(), class, limit)
	if err != nil {
		log.Printf("Error getting class leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get class leaderboard"})
		return
	}

	if err := cache.Set(cacheKey, leaderboard, 8*time.Hour); err != nil {
		log.Printf("Error setting class leaderboard in cache: %v", err)
	}

	c.JSON(http.StatusOK, leaderboard)
}

func (h *GlobalLeaderboardHandler) GetSpecLeaderboard(c *gin.Context) {
	class := c.Query("class")
	if class == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Class parameter is required"})
		return
	}

	spec := c.Query("spec")
	if spec == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spec parameter is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 {
		limit = 100
	}

	cacheKey := fmt.Sprintf("warcraftlogs:spec:%s:%s:limit:%d", class, spec, limit)

	var leaderboard interface{}
	err = cache.Get(cacheKey, &leaderboard)
	if err == nil {
		c.JSON(http.StatusOK, leaderboard)
		return
	}

	if err != redis.Nil {
		log.Printf("Error getting spec leaderboard from cache: %v", err)
	}

	leaderboard, err = h.rankingsService.GetGlobalLeaderboardBySpec(c.Request.Context(), class, spec, limit)
	if err != nil {
		log.Printf("Error getting spec leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spec leaderboard"})
		return
	}

	if err := cache.Set(cacheKey, leaderboard, 8*time.Hour); err != nil {
		log.Printf("Error setting spec leaderboard in cache: %v", err)
	}

	c.JSON(http.StatusOK, leaderboard)
}
