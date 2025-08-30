package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

func main() {
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

	database, err := db.InitDB()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to the database")
	} else {
		logger.Log.Info().Msg("Connected to the database successfully!")
	}
	defer database.Close()

	repo, err := db.InitTaskSkeletons(database)
	if err != nil {
		log.Fatal(err)
	}

	task := &models.TaskSkeleton{
		OwnerID: 1,
	}

	result, err := repo.Create(task)
	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("Success! Created task with ID: %d\n", result.ID)

	// Проверяем получение
	retrieved, err := repo.GetByID(result.ID)
	if err != nil {
		log.Fatal("Error getting task:", err)
	}

	fmt.Printf("Retrieved task: %+v\n", retrieved)
}
