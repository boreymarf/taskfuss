package db

import (
	"database/sql"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type TaskEntryRepository struct {
	db *sql.DB
}

func InitTaskEntryRepository(db *sql.DB) (*TaskEntryRepository, error) {

	repo := &TaskEntryRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("taskRepository initialization completed")

	return repo, nil
}

func (r *TaskEntryRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS task_entries (
	id 						INTEGER NOT NULL PRIMARY KEY,
	task_id 			INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	date 					DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	done          BOOLEAN NOT NULL DEFAULT FALSE
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
