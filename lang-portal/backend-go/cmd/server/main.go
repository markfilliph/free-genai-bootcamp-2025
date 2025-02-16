package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"lang-portal/internal/handlers"
	"lang-portal/internal/logger"
	"lang-portal/internal/middleware"
	"lang-portal/internal/models"
	"lang-portal/internal/service"
)

var (
	dashboardService *service.DashboardService
	groupService     *service.GroupServiceImpl
	studyService     *service.StudyService
	wordService      *service.WordService
)

func main() {
	// Load environment-specific configuration
	if err := loadEnvConfig(); err != nil {
		log.Fatal("Failed to load environment configuration:", err)
	}

	// Initialize logger
	logConfig := logger.Config{
		// Basic settings
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		LogFile:       getEnv("LOG_FILE", filepath.Join("logs", fmt.Sprintf("app.%s.log", getEnv("APP_ENV", "development")))),
		EnableConsole: getEnvBool("LOG_CONSOLE", true),

		// Rotation settings
		MaxSize:    getEnvInt("LOG_MAX_SIZE", 10),
		MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 30),
		MaxAge:     getEnvInt("LOG_MAX_AGE", 90),
		Compress:   getEnvBool("LOG_COMPRESS", true),
	}

	if err := logger.Initialize(logConfig); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize database with MySQL
	if err := models.InitDB(""); err != nil {
		logger.Fatal("Failed to initialize database", logger.Fields{"error": err.Error()})
	}
	defer models.CloseDB()

	// Initialize services
	dashboardService = service.NewDashboardService()
	groupService = service.NewGroupServiceImpl()
	studyService = service.NewStudyService()
	wordService = service.NewWordService()

	// Initialize handlers
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, studyService)
	groupHandler := handlers.NewGroupHandler(groupService)
	studySessionHandler := handlers.NewStudySessionHandler(studyService)
	wordHandler := handlers.NewWordHandler(wordService)

	// Initialize Gin router
	r := gin.New() // Use gin.New() instead of gin.Default() to avoid default logger

	// Add global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogging())
	r.Use(middleware.ValidationMiddleware())

	// Initialize routes
	initializeRoutes(r, dashboardHandler, groupHandler, studySessionHandler, wordHandler)

	// Get server port from environment variable, default to 8080
	port := ":8080"
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		port = ":" + envPort
	}

	// Start server
	logger.Info("Starting server", logger.Fields{"port": port})
	if err := r.Run(port); err != nil {
		logger.Fatal("Failed to start server", logger.Fields{"error": err.Error()})
	}
}

func initializeRoutes(
	r *gin.Engine,
	dashboardHandler *handlers.DashboardHandler,
	groupHandler *handlers.GroupHandler,
	studySessionHandler *handlers.StudySessionHandler,
	wordHandler *handlers.WordHandler,
) {
	// API group
	api := r.Group("/api")
	
	// Dashboard routes
	api.GET("/dashboard/last-study-session", dashboardHandler.GetLastStudySession)
	api.GET("/dashboard/study-progress", dashboardHandler.GetStudyProgress)
	api.GET("/dashboard/quick-stats", dashboardHandler.GetQuickStats)

	// Study activities routes
	api.GET("/study-activities/:id", studySessionHandler.GetStudyActivity)
	api.GET("/study-activities/:id/sessions", studySessionHandler.GetStudyActivitySessions)
	api.POST("/study-activities", studySessionHandler.CreateStudyActivity)

	// Words routes
	api.GET("/words", wordHandler.GetWords)
	api.GET("/words/:id", wordHandler.GetWord)

	// Groups routes
	api.GET("/groups", groupHandler.GetGroups)
	api.POST("/groups", groupHandler.CreateGroup)
	api.GET("/groups/:id", groupHandler.GetGroup)
	api.GET("/groups/:id/words", groupHandler.GetGroupWords)
	api.GET("/groups/:id/study-sessions", groupHandler.GetGroupStudySessions)

	// Study sessions routes
	api.GET("/study-sessions", studySessionHandler.GetStudySessions)
	api.GET("/study-sessions/:id", studySessionHandler.GetStudySession)
	api.GET("/study-sessions/:id/words", studySessionHandler.GetStudySessionWords)
	api.POST("/study-sessions/:id/words/:word_id/review", studySessionHandler.ReviewWordInSession)

	// Reset routes
	api.POST("/reset/history", dashboardHandler.ResetHistory)
	api.POST("/reset/full", dashboardHandler.FullReset)
}

// loadEnvConfig loads the appropriate .env file based on APP_ENV
func loadEnvConfig() error {
	// Default to development if APP_ENV is not set
	env := getEnv("APP_ENV", "development")

	// Load base .env file first (if exists)
	_ = godotenv.Load()

	// Load environment-specific .env file
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("error loading %s: %v", envFile, err)
	}

	// Set GIN_MODE based on environment
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as int or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
}

// getEnvBool gets an environment variable as bool or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Warning: Invalid boolean value for %s, using default: %v", key, defaultValue)
	}
	return defaultValue
}