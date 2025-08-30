package db

import (
	"database/sql"
	"fmt"
)

type RequirementSnapshots struct {
	db *sql.DB
}

func InitRequirementSnapshots(db *sql.DB) (*RequirementSnapshots, error) {

	repo := &RequirementSnapshots{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *RequirementSnapshots) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_snapshots (
		revision_uuid TEXT NOT NULL REFERENCES task_snapshots(revision_uuid) ON DELETE CASCADE,
		skeleton_id INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		parent_id INTEGER DEFAULT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('atom', 'condition')),
		data_type TEXT CHECK (data_type IN ('bool', 'int', 'float', 'duration', 'none')),
		operator TEXT CHECK (operator IN ('or', 'not', 'and', '==', '>=', '<=', '!=', '>', '<')),
		target_value TEXT,
		sort_order INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (revision_uuid, skeleton_id)
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
