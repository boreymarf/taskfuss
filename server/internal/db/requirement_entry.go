package db

import (
	"database/sql"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type RequirementEntryRepository struct {
	db *sql.DB
}

func InitRequirementEntryRepository(db *sql.DB) (*RequirementEntryRepository, error) {

	repo := &RequirementEntryRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("Repository initialization completed")

	return repo, nil
}

func (r *RequirementEntryRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS requirement_entries (
	id 							INTEGER NOT NULL PRIMARY KEY,
	requirement_id 	INTEGER NOT NULL REFERENCES requirements(id) ON DELETE CASCADE,
	entry_date 			DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	value						TEXT NOT NULL
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
