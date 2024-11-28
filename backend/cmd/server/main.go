package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	// Internal Packages - API Handlers
	authHandler "wowperf/internal/api/auth"
	apiBlizzard "wowperf/internal/api/blizzard"
	bnetAuthHandler "wowperf/internal/api/blizzard/auth"
	protectedProfileHandler "wowperf/internal/api/blizzard/protected/profile"
	"wowperf/internal/api/raiderio"
	userHandler "wowperf/internal/api/user"
	apiWarcraftlogs "wowperf/internal/api/warcraftlogs"

	// Internal Packages - Services
	auth "wowperf/internal/services/auth"
	serviceBlizzard "wowperf/internal/services/blizzard"
	bnetAuth "wowperf/internal/services/blizzard/auth"
	serviceRaiderio "wowperf/internal/services/raiderio"
	mythicplusUpdate "wowperf/internal/services/raiderio/mythicplus"
	userService "wowperf/internal/services/user"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	warcraftLogsLeaderboard "wowperf/internal/services/warcraftlogs/dungeons"

	// Internal Packages - Database
	"wowperf/internal/database"
	migrations "wowperf/internal/database/migrations"

	// Internal Packages - Middleware & Utils
	cacheMiddleware "wowperf/middleware/cache"
	"wowperf/pkg/cache"
	csrfMiddleware "wowperf/pkg/middleware"
	authMiddleware "wowperf/pkg/middleware/auth"             // JWT middleware
	blizzardAuthMiddleware "wowperf/pkg/middleware/blizzard" // Battle.net middleware
)

// Structures pour regrouper les services et handlers
type AppServices struct {
	Auth            *auth.AuthService
	BattleNet       *bnetAuth.BattleNetAuthService
	User            *userService.UserService
	Blizzard        *serviceBlizzard.Service
	RaiderIO        *serviceRaiderio.RaiderIOService
	WarcraftLogs    *warcraftlogs.WarcraftLogsClientService
	LeaderBoard     *warcraftLogsLeaderboard.GlobalLeaderboardService
	RankingsUpdater *warcraftLogsLeaderboard.RankingsUpdater
}

type AppHandlers struct {
	Auth             *authHandler.AuthHandler
	User             *userHandler.UserHandler
	BattleNet        *bnetAuthHandler.BattleNetAuthHandler
	RaiderIO         *raiderio.Handler
	Blizzard         *apiBlizzard.Handler
	WarcraftLogs     *apiWarcraftlogs.Handler
	ProtectedProfile *protectedProfileHandler.Handler
}

type AppConfig struct {
	Environment    string
	AllowedOrigins []string
	Port           string
	JWTSecret      string
	CSRFSecret     string
}

// Cache managers structure
type CacheManagers struct {
	RaiderIO     *cacheMiddleware.CacheManager
	Blizzard     *cacheMiddleware.CacheManager
	WarcraftLogs *cacheMiddleware.CacheManager
}

// Initialisation des services
func initializeServices(db *gorm.DB, cacheService cache.CacheService, cacheManagers CacheManagers) (*AppServices, error) {
	// Get Redis client from cache service
	redisClient := cacheService.GetRedisClient()
	if redisClient == nil {
		return nil, fmt.Errorf("failed to get Redis client from cache service")
	}

	// Configuration for Battle.net auth service
	battleNetService, err := bnetAuth.NewBattleNetAuthService(db, redisClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize battle.net auth service: %w", err)
	}

	// Main authentication service
	authService := auth.NewAuthService(
		db,
		os.Getenv("JWT_SECRET"),
		redisClient,
	)

	// Other services...
	userSvc := userService.NewUserService(db)

	blizzardService, err := serviceBlizzard.NewService(db, redisClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize blizzard service: %w", err)
	}

	rioService, err := serviceRaiderio.NewRaiderIOService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize raiderio service: %w", err)
	}

	warcraftLogsService, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize warcraftlogs service: %w", err)
	}

	globalLeaderboardService := warcraftLogsLeaderboard.NewGlobalLeaderboardService(db)
	rankingsUpdater := warcraftLogsLeaderboard.NewRankingsUpdater(
		db,
		warcraftLogsService,
		cacheService,
		cacheManagers.WarcraftLogs,
	)

	return &AppServices{
		Auth:            authService,
		BattleNet:       battleNetService,
		User:            userSvc,
		Blizzard:        blizzardService,
		RaiderIO:        rioService,
		WarcraftLogs:    warcraftLogsService,
		LeaderBoard:     globalLeaderboardService,
		RankingsUpdater: rankingsUpdater,
	}, nil
}

// Initialisation des handlers
func initializeHandlers(services *AppServices, db *gorm.DB, cacheService cache.CacheService, cacheManagers CacheManagers) *AppHandlers {
	return &AppHandlers{
		Auth:      authHandler.NewAuthHandler(services.Auth),
		User:      userHandler.NewUserHandler(services.User),
		BattleNet: bnetAuthHandler.NewBattleNetAuthHandler(services.BattleNet),
		RaiderIO:  raiderio.NewHandler(services.RaiderIO, db, cacheService, cacheManagers.RaiderIO),
		Blizzard:  apiBlizzard.NewHandler(services.Blizzard, db, cacheService, cacheManagers.Blizzard),
		WarcraftLogs: apiWarcraftlogs.NewHandler(
			services.LeaderBoard,
			services.WarcraftLogs,
			db,
			cacheService,
			cacheManagers.WarcraftLogs,
		),
		ProtectedProfile: protectedProfileHandler.NewHandler(services.Blizzard.ProtectedProfile),
	}
}

// Configuration des routes
func setupRoutes(
	r *gin.Engine,
	services *AppServices,
	handlers *AppHandlers,
) {
	// Middlewares
	jwtMiddleware := authMiddleware.JWTAuth(services.Auth)
	bnetMiddleware := blizzardAuthMiddleware.NewBattleNetMiddleware(services.BattleNet)

	// CSRF Token endpoint
	r.GET("/csrf-token", csrfMiddleware.GetCSRFToken())

	// Authentication routes
	handlers.Auth.RegisterRoutes(r)
	handlers.BattleNet.RegisterRoutes(r, jwtMiddleware)

	// Protected API routes
	apiGroup := r.Group("")
	{
		// Routes protected by JWT
		protected := apiGroup.Group("")
		protected.Use(jwtMiddleware)

		// Routes requiring Battle.net
		bnetProtected := protected.Group("")
		bnetProtected.Use(bnetMiddleware.RequireBattleNetAccount())
		bnetProtected.Use(bnetMiddleware.RequireValidToken())

		// Protected Blizzard API routes
		handlers.ProtectedProfile.RegisterRoutes(bnetProtected)

		// Other API routes
		handlers.RaiderIO.RegisterRoutes(r)
		handlers.Blizzard.RegisterRoutes(r)
		handlers.WarcraftLogs.RegisterRoutes(r)
		handlers.User.RegisterRoutes(r, services.Auth)
	}
}

// Helper function for CORS
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

// Redis initialization
func initializeCacheService() (cache.CacheService, error) {
	cacheService, err := cache.NewRedisCache(&cache.Config{
		URL:      os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache service: %w", err)
	}

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cacheService.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to cache: %w", err)
	}

	return cacheService, nil
}

// Initialize cache managers
func initializeCacheManagers(cacheService cache.CacheService) CacheManagers {
	return CacheManagers{
		RaiderIO: cacheMiddleware.NewCacheManager(cacheMiddleware.CacheConfig{
			Cache:      cacheService,
			Expiration: 24 * time.Hour,
			KeyPrefix:  "raiderio",
			Metrics:    true,
		}),
		Blizzard: cacheMiddleware.NewCacheManager(cacheMiddleware.CacheConfig{
			Cache:      cacheService,
			Expiration: 24 * time.Hour,
			KeyPrefix:  "blizzard",
			Metrics:    true,
		}),
		WarcraftLogs: cacheMiddleware.NewCacheManager(cacheMiddleware.CacheConfig{
			Cache:      cacheService,
			Expiration: 2 * time.Hour,
			KeyPrefix:  "warcraftlogs",
			Tags:       []string{"rankings", "leaderboard"},
			Metrics:    true,
		}),
	}
}

// Initialize database
func initializeDatabase() (*gorm.DB, error) {
	db, err := database.InitDB()
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	if err := migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	if err := database.InitializeDatabase(db); err != nil {
		return nil, fmt.Errorf("database seeding failed: %w", err)
	}

	return db, nil
}

// Load configuration
func loadConfig() (*AppConfig, error) {
	// Ignore silently the .env loading error
	_ = godotenv.Load()

	// Required environment variables
	requiredEnvVars := []string{
		"JWT_SECRET",
		"CSRF_SECRET",
		"REDIS_URL",
		"BLIZZARD_CLIENT_ID",
		"BLIZZARD_CLIENT_SECRET",
		"BLIZZARD_REDIRECT_URL",
	}

	// Check required variables with more logs
	missingVars := []string{}
	for _, envVar := range requiredEnvVars {
		val := strings.TrimSpace(os.Getenv(envVar))
		if val == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	// If some variables are missing, log the details and return an error
	if len(missingVars) > 0 {
		log.Printf("Missing required environment variables: %v", missingVars)
		return nil, fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	// If all variables are present, create and return the config
	config := &AppConfig{
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
		AllowedOrigins: strings.Split(getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		Port:           getEnvOrDefault("PORT", "8080"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		CSRFSecret:     os.Getenv("CSRF_SECRET"),
	}

	return config, nil
}

// Setup server middleware
func setupMiddleware(r *gin.Engine, config *AppConfig) {
	// Security headers
	r.Use(securityHeaders())

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"Accept",
			"Origin",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-CSRF-Token", "Set-Cookie", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Logger
	if config.Environment == "development" {
		r.Use(gin.Logger())
	}
}

// Start periodic tasks
func startPeriodicTasks(db *gorm.DB, services *AppServices) {
	// Mythic+ Updates
	go func() {
		log.Println("Setting up dungeon stats update...")
		if mythicplusUpdate.CheckAndSetUpdateLock(db) {
			log.Println("Performing initial dungeon stats update...")
			if err := mythicplusUpdate.UpdateDungeonStats(db, services.RaiderIO); err != nil {
				log.Printf("Error during initial dungeon stats update: %v", err)
			} else {
				log.Println("Initial dungeon stats update completed")
			}
		}
		mythicplusUpdate.StartWeeklyDungeonStatsUpdate(db, services.RaiderIO)
	}()

	// WarcraftLogs Updates
	go func() {
		log.Println("Setting up WarcraftLogs rankings update scheduler...")
		time.Sleep(10 * time.Second) // Wait for DB readiness
		services.RankingsUpdater.StartPeriodicUpdate(context.Background())
	}()
}

func main() {
	log.Println("Starting server...")

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Initialize components
	db, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	cacheService, err := initializeCacheService()
	if err != nil {
		log.Fatalf("Failed to initialize cache service: %v", err)
	}
	defer cacheService.Close()
	log.Println("Cache service initialized successfully")
	cacheManagers := initializeCacheManagers(cacheService)

	// Initialize services
	services, err := initializeServices(db, cacheService, cacheManagers)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}
	log.Println("Services initialized successfully")

	// Initialize handlers
	handlers := initializeHandlers(services, db, cacheService, cacheManagers)
	log.Println("Handlers initialized successfully")

	// Setup server
	r := gin.New()
	setupMiddleware(r, config)
	setupHealthCheck(r)
	setupRoutes(r, services, handlers)
	log.Println("Routes setup successfully")

	// Start periodic tasks
	startPeriodicTasks(db, services)
	log.Println("Periodic tasks started successfully")

	// Start server
	serverAddr := fmt.Sprintf(":%s", config.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupHealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		if os.Getenv("ENVIRONMENT") == "production" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}
