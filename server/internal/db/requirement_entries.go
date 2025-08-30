package db

import (
	"database/sql"
	"fmt"
)

type RequirementEntries struct {
	db *sql.DB
}

func InitRequirementEntries(db *sql.DB) (*RequirementEntries, error) {

	repo := &RequirementEntries{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *RequirementEntries) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_entries (
		id 							INTEGER NOT NULL PRIMARY KEY,
		requirement_id 	INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		revision_uuid		TEXT NOT NULL REFERENCES requirement_snapshots(revision_uuid) ON DELETE CASCADE,
		entry_date			DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		value TEXT
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
