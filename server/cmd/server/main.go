package main

import (
	"os"

	// "github.com/boreymarf/task-fuss/server/internal/config"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/middleware"
	"github.com/boreymarf/task-fuss/server/internal/routes"
	"github.com/boreymarf/task-fuss/server/internal/service"
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
			logger.Log.Info().Msgf("Current log level is set to %s", level)
		} else {
			logger.Log.Error().Str("log_level", level).Msg("Failed to set log level!")
		}
	}

	// if os.Getenv("APP_ENV") == config.EnvDevelopment {
	// 	logger.Log.Warn().Msg("APP_ENV is set to production!")
	// }

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
	r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.82"})

	// Middleware
	r.Use(middleware.SimpleMiddleware())

	// TODO: Потом заменить на Prod и Dev вариации
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:5173", "http://192.168.1.82:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Repositories
	userRepository, err := db.InitUserRepository(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize user repository")
	}

	taskRepository, err := db.InitTaskRepository(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task repository")
	}

	taskEntryRepository, err := db.InitTaskEntryRepository(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task entry repository")
	}

	requirementRepository, err := db.InitRequirementRepository(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize requirement repository")
	}

	requirementEntryRepository, err := db.InitRequirementEntryRepository(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize requirement entry repository")
	}

	// Services
	taskService, err := service.InitTaskService(
		taskRepository,
		taskEntryRepository,
		requirementRepository,
		requirementEntryRepository,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task service repository")
	}

	// Handlers
	authHandler, err := handlers.InitAuthHandler(userRepository)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize auth handler")
	}

	profileHandler, err := handlers.InitProfileHandler(userRepository)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize profile handler")
	}

	taskHandler, err := handlers.InitTaskHandler(
		userRepository,
		taskRepository,
		taskEntryRepository,
		requirementRepository,
		requirementEntryRepository,
		taskService,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task handler")
	}

	routes.SetupAPIRoutes(r, userRepository, authHandler, profileHandler, taskHandler)

	// Initializing
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	addr := "0.0.0.0:" + port
	logger.Log.Info().Msgf("Server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
