package main

import (
	"log"
	"os"
	"time"
	authHandler "wowperf/internal/api/auth"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"
	"wowperf/internal/database"
	auth "wowperf/internal/services/auth"

	userHandler "wowperf/internal/api/user"
	userService "wowperf/internal/services/user"

	"wowperf/pkg/cache"
	"wowperf/pkg/middleware"

	mythicPlusRaiderioCache "wowperf/internal/api/raiderio/mythicplus"
	raidsRaiderioCache "wowperf/internal/api/raiderio/raids"

	models "wowperf/internal/models/raiderio/mythicrundetails"

	serviceBlizzard "wowperf/internal/services/blizzard"
	serviceRaiderio "wowperf/internal/services/raiderio"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func checkUpdateState(db *gorm.DB) {
	var state models.UpdateState
	result := db.First(&state)
	if result.Error != nil {
		log.Printf("Error fetching update state: %v", result.Error)
		return
	}
	log.Printf("Current update state: Last update was %v ago", time.Since(state.LastUpdateTime))
}

func main() {

	// Load environment variables
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

	if err := database.InitializeDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Services

	// Auth Service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set in the environment")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
	jwtExpiration := 24 * time.Hour // Default to 24 hours
	if jwtExpirationStr != "" {
		duration, err := time.ParseDuration(jwtExpirationStr)
		if err == nil {
			jwtExpiration = duration
		}
	}

	authService := auth.NewAuthService(
		db,
		jwtSecret,
		cache.GetRedisClient(),
		jwtExpiration,
		os.Getenv("BATTLE_NET_CLIENT_ID"),
		os.Getenv("BATTLE_NET_CLIENT_SECRET"),
		os.Getenv("BATTLE_NET_REDIRECT_URL"),
	)
	authHandler := authHandler.NewAuthHandler(authService)

	// User Service
	userService := userService.NewUserService(db)
	userHandler := userHandler.NewUserHandler(userService)

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

	// Raider.io Handler
	rioHandler := raiderio.NewHandler(rioService, db)

	// Blizzard Handler
	blizzardHandler := apiBlizzard.NewHandler(blizzardService, db)

	r := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	r.Use(gin.Logger())
	// API

	// Auth API
	r.GET("/csrf-token", middleware.CSRFToken())
	authHandler.RegisterRoutes(r)

	// Raider.io API
	rioHandler.RegisterRoutes(r)

	// Blizzard API
	blizzardHandler.RegisterRoutes(r)

	checkUpdateState(db)

	/*
		go func() {
			log.Println("Checking if dungeon stats update is needed...")
			if err := raiderioMythicPlus.UpdateDungeonStats(db, rioService); err != nil {
				log.Printf("Error during dungeon stats update check: %v", err)
			} else {
				log.Println("Dungeon stats update check completed")
			}
			checkUpdateState(db)

			log.Println("Setting up weekly dungeon stats update...")
			raiderioMythicPlus.StartWeeklyDungeonStatsUpdate(db, rioService)
		}()
	*/

	// Protected Routes
	userHandler.RegisterRoutes(r, authService.AuthMiddleware())

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
