package main

import (
	"os"

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

// @title TaskFuss API
// @version 1.0.0
// @description API for TaskFuss app

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @scheme bearer
// @bearerFormat JWT

// @host localhost:4000
// @BasePath /api
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
	r.Use(middleware.MetricsMiddleware())
	r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.82"})

	// Middleware
	r.Use(middleware.SimpleMiddleware())

	// TODO: Потом заменить на Prod и Dev вариации
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:5173", "http://192.168.1.82:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Repositories
	usersRepository, err := db.InitUsers(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize user repository")
	}

	taskSkeletonsRepository, err := db.InitTaskSkeletons(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task skeletons repository")
	}

	taskSnapshotsRepository, err := db.InitTaskSnapshots(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task snapshots repository")
	}

	taskEntriesRepository, err := db.InitTaskEntries(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task entry repository")
	}

	taskPeriodsRepository, err := db.InitTaskPeriods(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task periods repository")
	}

	requirementSkeletonsRepository, err := db.InitRequirementSkeletons(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize requirement repository")
	}

	requirementSnapshotsRepository, err := db.InitRequirementSnapshots(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize requirements snapshots repository")
	}

	requirementEntriesRepository, err := db.InitRequirementEntries(database)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize requirement entry repository")
	}

	// Services
	taskService, err := service.InitTaskService(
		database,
		usersRepository,
		taskSkeletonsRepository,
		taskSnapshotsRepository,
		taskPeriodsRepository,
		taskEntriesRepository,
		requirementSkeletonsRepository,
		requirementSnapshotsRepository,
		requirementEntriesRepository,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task service repository")
	}
	// Handlers
	authHandler, err := handlers.InitAuthHandler(usersRepository)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize auth handler")
	}

	profileHandler, err := handlers.InitProfileHandler(usersRepository)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize profile handler")
	}

	taskHandler, err := handlers.InitTaskHandler(
		taskService,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Unable to initialize task handler")
	}

	routes.SetupAPIRoutes(
		r,
		usersRepository,
		authHandler,
		profileHandler,
		taskHandler,
	)

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
