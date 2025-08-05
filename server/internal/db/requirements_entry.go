package db

import (
	"database/sql"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type RequirementsEntryRepository struct {
	db *sql.DB
}

func InitRequirementsEntryRepository(db *sql.DB) (*RequirementsEntryRepository, error) {

	repo := &RequirementsEntryRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("Repository initialization completed")

	return repo, nil
}

func (r *RequirementsEntryRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS requirements_entries (
	id 											INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	requirement_id					INTEGER NOT NULL REFERENCES requirements(id) ON DELETE CASCADE,
	requirement_snapshot_id INTEGER NOT NULL REFERENCES requirements_snapshots(id) ON DELETE CASCADE,
	entry_date 							DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	value										TEXT NOT NULL,
	recorded_at							DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
