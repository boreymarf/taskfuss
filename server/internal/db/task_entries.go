package db

import (
	"database/sql"
	"fmt"
)

type TaskEntries struct {
	db *sql.DB
}

func InitTaskEntries(db *sql.DB) (*TaskEntries, error) {

	repo := &TaskEntries{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *TaskEntries) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS task_entries (
		id 						INTEGER NOT NULL PRIMARY KEY,
		task_id 			INTEGER NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		entry_date 		DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed			BOOLEAN NOT NULL DEFAULT FALSE CHECK (completed IN (0, 1))
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
