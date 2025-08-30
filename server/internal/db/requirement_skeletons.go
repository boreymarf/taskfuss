package db

import (
	"database/sql"
	"fmt"
)

type RequirementSkeletons struct {
	db *sql.DB
}

func InitRequirementSkeletons(db *sql.DB) (*RequirementSkeletons, error) {

	repo := &RequirementSkeletons{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *RequirementSkeletons) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_skeletons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		FOREIGN KEY (task_id) REFERENCES task_skeletons(id) ON DELETE CASCADE
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
