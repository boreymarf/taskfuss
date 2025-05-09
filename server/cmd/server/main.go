package main

import (
	"os"

	"github.com/boreymarf/task-fuss/server/internal/config"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/middleware"
	"github.com/boreymarf/task-fuss/server/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {

	logger.Log.Info().Msg("Starting server...")
	logger.Log.Debug().Msg("Current log level is set to Debug")

	// Loading .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal().Msg("Error loading .env file")
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		if logLevel, err := zerolog.ParseLevel(level); err == nil {
			logger.Log = logger.Log.Level(logLevel)
		}
	}

	if os.Getenv("APP_ENV") == config.EnvDevelopment {
		logger.Log.Warn().Msg("APP_ENV is set to production!")
	}

	// Connecting to database
	database, err := db.InitDB()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to the database")
	} else {
		logger.Log.Info().Msg("Connected to the database successfully!")
	}
	defer database.Close()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Middleware
	r.Use(middleware.SimpleMiddleware())

	// TODO: Потом заменить на Prod и Dev вариации
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Handlers and repositories
	userRepository, err := db.InitUserRepository(database)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Unable to initialize user repository")
	}

	authHandler, err := handlers.InitAuthHanlder(userRepository)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Unable to initialize auth handler")
	}

	routes.SetupAPIRoutes(r, authHandler)

	// Initializing
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	if err := r.Run(":" + port); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
