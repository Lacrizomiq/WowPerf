package main

import (
	"log"
	"os"
	"time"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"
	"wowperf/internal/database"
	"wowperf/pkg/cache"

	apiAuth "wowperf/internal/api/auth"
	authService "wowperf/internal/services/auth"

	mythicPlusRaiderioCache "wowperf/internal/api/raiderio/mythicplus"
	raidsRaiderioCache "wowperf/internal/api/raiderio/raids"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"

	serviceBlizzard "wowperf/internal/services/blizzard"
	serviceRaiderio "wowperf/internal/services/raiderio"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	// Init Cache
	cache.InitCache()

	// Wait for Redis to be ready
	waitForRedis()

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Use the new migration function
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.SeedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Services

	// Auth Service
	tokenExpiry := 24 * time.Hour
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatalf("JWT_KEY is not set")
	}
	authService := authService.NewAuthService(db, jwtKey, cache.GetRedisClient(), tokenExpiry)

	// Blizzard Service
	blizzardService, err := serviceBlizzard.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize blizzard service: %v", err)
	}

	// Raider.io Service
	rioService, err := serviceRaiderio.NewRaiderIOService()
	if err != nil {
		log.Fatalf("Failed to initialize raiderio service: %v", err)
	}

	// Cache Updater
	startCacheUpdater(blizzardService, rioService)

	// Handlers
	// Auth Handler
	authHandler := apiAuth.NewAuthHandler(authService)

	// Raider.io Handler
	rioHandler := raiderio.NewHandler(rioService, db)

	// Blizzard Handler
	blizzardHandler := apiBlizzard.NewHandler(blizzardService, db)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	// API

	// Auth API
	authHandler.RegisterRoutes(r)

	// Raider.io API
	rioHandler.RegisterRoutes(r)

	// Blizzard API
	blizzardHandler.RegisterRoutes(r)

	// Start Dungeon Stats Update
	go func() {
		log.Println("Checking if dungeon stats update is needed...")
		if err := raiderioMythicPlus.UpdateDungeonStats(db, rioService); err != nil {
			log.Printf("Error updating dungeon stats: %v", err)
		} else {
			log.Println("Dungeon stats update check completed")
		}

		log.Println("Setting up weekly dungeon stats update...")
		raiderioMythicPlus.StartWeeklyDungeonStatsUpdate(db, rioService)
	}()

	log.Println("Server is starting on :8080")
	log.Fatal(r.Run(":8080"))
}

func startCacheUpdater(blizzardService *serviceBlizzard.Service, rioService *serviceRaiderio.RaiderIOService) {
	raidsRaiderioCache.StartRaidLeaderboardCacheUpdater(rioService)
	mythicPlusRaiderioCache.StartMythicPlusBestRunsCacheUpdater(rioService)
}

func waitForRedis() {
	for i := 0; i < 30; i++ {
		err := cache.Ping()
		if err == nil {
			log.Println("Redis is ready")
			return
		}
		log.Println("Waiting for Redis to be ready")
		time.Sleep(time.Second)
	}
	log.Println("Redis is not ready")
}
