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

	logger.Log.Debug().Msg("taskRepository initialization completed")

	return repo, nil
}

func (r *RequirementEntryRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS task_entries (
	id 						INTEGER NOT NULL PRIMARY KEY,
	task_id 			INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	date 					DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	value
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
