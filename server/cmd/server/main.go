package main

import (
	"os"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	logger.Log.Info().Msg("Starting server...")

	// Loading .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal().Msg("Error loading .env file")
	}

	// Connecting to database
	database, err := db.InitDB()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to the database")
	} else {
		logger.Log.Info().Msg("Connected to the database successfully!")
	}
	defer database.Close()

	// Create repositories
	// userRepository := db.InitUserRepository(database)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// TODO: Потом заменить на Prod и Dev вариации
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.SetupAPIRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Запуск
	if err := r.Run(":" + port); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
