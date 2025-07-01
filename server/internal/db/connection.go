package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	_ "github.com/mattn/go-sqlite3" // Used for sql.Open("sqlite3"), without it will panic
)

func InitDB() (*sql.DB, error) {
	// Create directory for the db
	db_path := os.Getenv("DB_PATH")
	dir := filepath.Dir(db_path)
	logger.Log.Info().Str("db_path", db_path).Msg("Creating database...")

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Log.Fatal().Err(err).Str("path", dir).Msg("Failed to create DB")
	}

	// Connection
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
