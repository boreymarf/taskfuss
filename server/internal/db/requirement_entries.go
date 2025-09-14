package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/models"
)

type RequirementEntries interface {
	UpsertTx(ctx context.Context, tx *sql.Tx, entry *models.RequirementEntry) (*models.RequirementEntry, error)
}

type requirementEntries struct {
	db *sql.DB
}

var _ RequirementEntries = (*requirementEntries)(nil)

func InitRequirementEntries(db *sql.DB) (RequirementEntries, error) {

	repo := &requirementEntries{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *requirementEntries) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_entries (
		id 							INTEGER NOT NULL PRIMARY KEY,
		requirement_id 	INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		revision_uuid		TEXT NOT NULL REFERENCES requirement_snapshots(revision_uuid) ON DELETE CASCADE,
		entry_date			DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		value TEXT
		UNIQUE(requirement_id, entry_date)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *requirementEntries) UpsertTx(ctx context.Context, tx *sql.Tx, entry *models.RequirementEntry) (*models.RequirementEntry, error) {
	query := `
        INSERT INTO requirement_entries (requirement_id, revision_uuid, entry_date, value)
        VALUES (?, ?, ?, ?)
        ON CONFLICT (requirement_id, entry_date) 
        DO UPDATE SET 
            revision_uuid = excluded.revision_uuid,
            value = excluded.value
        RETURNING id, requirement_id, revision_uuid, entry_date, value`

	row := tx.QueryRowContext(
		ctx,
		query,
		entry.RequirementID,
		entry.RevisionUUID,
		entry.EntryDate,
		entry.Value,
	)

	var updatedEntry models.RequirementEntry
	err := row.Scan(
		&updatedEntry.ID,
		&updatedEntry.RequirementID,
		&updatedEntry.RevisionUUID,
		&updatedEntry.EntryDate,
		&updatedEntry.Value,
	)
	if err != nil {
		// var sqliteErr sqlite3.Error
		// if errors.As(err, &sqliteErr) {
		// 	if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
		// 		return nil, apperrors.ErrForeignKeyViolation
		// 	}
		// }
		return nil, err
	}

	return &updatedEntry, nil
}
