package core

import (
	"os"
	"path/filepath"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func InitDB() (*sqlx.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/db.sqlite"
	}

	dir := filepath.Dir(dbPath)
	logger.Log.Info().Str("db_path", dbPath).Msg("Creating database directory if missing")
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Log.Fatal().Err(err).Str("path", dir).Msg("Failed to create DB directory")
	}

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Log.Info().Str("db_path", dbPath).Msg("Database connected")
	return db, nil
}
