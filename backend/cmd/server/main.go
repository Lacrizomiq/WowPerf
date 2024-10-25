package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	// Internal Packages - API Handlers
	authHandler "wowperf/internal/api/auth"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"
	mythicPlusRaiderioCache "wowperf/internal/api/raiderio/mythicplus"
	raidsRaiderioCache "wowperf/internal/api/raiderio/raids"
	userHandler "wowperf/internal/api/user"
	apiWarcraftlogs "wowperf/internal/api/warcraftlogs"

	// Internal Packages - Services
	auth "wowperf/internal/services/auth"
	serviceBlizzard "wowperf/internal/services/blizzard"
	serviceRaiderio "wowperf/internal/services/raiderio"
	mythicplusUpdate "wowperf/internal/services/raiderio/mythicplus"
	userService "wowperf/internal/services/user"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	warcraftLogsRankings "wowperf/internal/services/warcraftlogs/dungeons"

	// Internal Packages - Models & Database
	"wowperf/internal/database"
	models "wowperf/internal/models/raiderio/mythicrundetails"

	// Internal Packages - Utils
	"wowperf/pkg/cache"
	"wowperf/pkg/middleware"
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

func initializeServices(db *gorm.DB) (
	*auth.AuthService,
	*userService.UserService,
	*serviceBlizzard.Service,
	*serviceRaiderio.RaiderIOService,
	*warcraftLogsRankings.RankingsService,
	*warcraftLogsRankings.RankingsUpdater,
	*warcraftLogsRankings.DungeonService,
	error,
) {
	// Auth Service Setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("JWT_SECRET must be set in the environment")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
	jwtExpiration := 24 * time.Hour
	if jwtExpirationStr != "" {
		if duration, err := time.ParseDuration(jwtExpirationStr); err == nil {
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

	// User Service
	userSvc := userService.NewUserService(db)

	// Blizzard Service
	blizzardService, err := serviceBlizzard.NewService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize blizzard service: %v", err)
	}

	// Raider.io Service
	rioService, err := serviceRaiderio.NewRaiderIOService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize raiderio service: %v", err)
	}

	// WarcraftLogs Service Setup
	warcraftLogsClient, err := warcraftlogs.NewClient()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize warcraftlogs client: %v", err)
	}

	dungeonService := warcraftLogsRankings.NewDungeonService(warcraftLogsClient)
	rankingsService := warcraftLogsRankings.NewRankingsService(dungeonService, db)
	rankingsUpdater := warcraftLogsRankings.NewRankingsUpdater(db, rankingsService)

	return authService, userSvc, blizzardService, rioService, rankingsService, rankingsUpdater, dungeonService, nil
}

func setupRoutes(
	r *gin.Engine,
	authService *auth.AuthService,
	authHandler *authHandler.AuthHandler,
	userHandler *userHandler.UserHandler,
	rioHandler *raiderio.Handler,
	blizzardHandler *apiBlizzard.Handler,
	warcraftlogsHandler *apiWarcraftlogs.Handler,
) {
	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))
	r.Use(gin.Logger())

	// Auth Routes
	r.GET("/csrf-token", middleware.CSRFToken())
	authHandler.RegisterRoutes(r)

	// API Routes
	rioHandler.RegisterRoutes(r)
	blizzardHandler.RegisterRoutes(r)
	warcraftlogsHandler.RegisterRoutes(r)

	// Protected Routes
	userHandler.RegisterRoutes(r, authService.AuthMiddleware())
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		return
	}

	// Initialize Cache and Redis
	cache.InitCache()
	waitForRedis()

	// Initialize Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.InitializeDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize Services
	authService, userSvc, blizzardService, rioService, rankingsService, rankingsUpdater, dungeonService, err := initializeServices(db)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Initialize Handlers
	authHandler := authHandler.NewAuthHandler(authService)
	userHandler := userHandler.NewUserHandler(userSvc)
	rioHandler := raiderio.NewHandler(rioService, db)
	blizzardHandler := apiBlizzard.NewHandler(blizzardService, db)
	warcraftlogsHandler := apiWarcraftlogs.NewHandler(rankingsService, dungeonService, db)

	// Start Cache Updater
	startCacheUpdater(blizzardService, rioService)

	// Start Periodic Updates (Mythic+ Dungeon Stats)
	go func() {
		log.Println("Setting up dungeon stats update...")
		if mythicplusUpdate.CheckAndSetUpdateLock(db) {
			log.Println("Performing initial dungeon stats update...")
			if err := mythicplusUpdate.UpdateDungeonStats(db, rioService); err != nil {
				log.Printf("Error during initial dungeon stats update: %v", err)
			} else {
				log.Println("Initial dungeon stats update completed")
			}
		} else {
			log.Println("Dungeon stats are up to date")
		}
		checkUpdateState(db)

		log.Println("Setting up weekly dungeon stats update...")
		mythicplusUpdate.StartWeeklyDungeonStatsUpdate(db, rioService)
	}()

	// Start Periodic Updates (WarcraftLogs Dungeon Rankings)
	go func() {
		log.Println("Setting up rankings update...")
		rankingsUpdater.CheckAndUpdate()
		rankingsUpdater.StartPeriodicUpdate()
	}()

	// Setup and Start Server
	r := gin.Default()
	setupRoutes(r, authService, authHandler, userHandler, rioHandler, blizzardHandler, warcraftlogsHandler)

	log.Println("Server is starting on :8080")
	log.Fatal(r.Run(":8080"))
}
