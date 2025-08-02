package db

import (
	"database/sql"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type EntriesRepository struct {
	db *sql.DB
}

func InitEntriesRepository(db *sql.DB) (*EntriesRepository, error) {

	repo := &EntriesRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("taskRepository initialization completed")

	return repo, nil
}

func (r *EntriesRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS entries (
	id             INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  task_id        INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	requirement_id INTEGER NOT NULL REFERENCES requirement(id) ON DELETE CASCADE,
	value          TEXT,
	value_type     TEXT NOT NULL CHECK(value_type IN ('int', 'float', 'time', 'bool')),
	timestamp      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
