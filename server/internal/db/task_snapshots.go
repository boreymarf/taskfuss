package db

import (
	"database/sql"
	"fmt"
)

type TaskSnapshots struct {
	db *sql.DB
}

func InitTaskSnapshots(db *sql.DB) (*TaskSnapshots, error) {

	repo := &TaskSnapshots{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *TaskSnapshots) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS task_snapshots (
		revision_uuid TEXT NOT NULL,
		task_id INTEGER  NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		description TEXT,
		start_date DATETIME,
		end_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_current BOOLEAN DEFAULT FALSE,
	PRIMARY KEY (revision_uuid, task_id)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
