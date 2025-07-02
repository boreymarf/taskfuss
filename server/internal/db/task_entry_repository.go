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
	query := `CREATE TABLE IF NOT EXISTS tasks (
	id 						INTEGER NOT NULL PRIMARY KEY,
	task_id 			INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	requirements  TEXT NOT NULL,
	name 					VARCHAR(255) NOT NULL,
	description 	VARCHAR(255),
	created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
	start_date 		DATETIME DEFAULT CURRENT_TIMESTAMP,
	end_date 			DATETIME DEFAULT NULL
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
