package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/google/uuid"
)

type RequirementEntries interface {
	UpsertTx(ctx context.Context, tx *sql.Tx, entry *models.RequirementEntry) (*models.RequirementEntry, error)

	GetByRequirementIDsTx(ctx context.Context, tx *sql.Tx, requirementIDs []int64, revisionUUID uuid.UUID) (map[int64]models.RequirementEntry, error)
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
		value TEXT,
		UNIQUE(requirement_id, entry_date)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *requirementEntries) UpsertTx(ctx context.Context, tx *sql.Tx, entry *models.RequirementEntry) (*models.RequirementEntry, error) {

	logger.Log.Debug().
		Str("revisionUUID", entry.RevisionUUID.String()).
		Int64("requirementID", entry.RequirementID).
		Msg("Trying to upsert a requirement entry in db via tx")

	query := `
		INSERT INTO requirement_entries (
			requirement_id,
			revision_uuid,
			entry_date,
			value
		)
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

func (r *requirementEntries) GetByRequirementIDsTx(ctx context.Context, tx *sql.Tx, requirementIDs []int64, revisionUUID uuid.UUID) (map[int64]models.RequirementEntry, error) {
	logger.Log.Debug().
		Interface("requirementIDs", requirementIDs).
		Str("revisionUUID", revisionUUID.String()).
		Msg("Trying to get requirement entries by requirement IDs and revision UUID via tx")

	if len(requirementIDs) == 0 {
		logger.Log.Warn().Msg("empty requirement IDs slice passed to GetByRequirementIDsAndRevision")
		return map[int64]models.RequirementEntry{}, nil
	}

	placeholders := make([]string, len(requirementIDs))
	params := make([]any, len(requirementIDs)+1)
	for i, id := range requirementIDs {
		placeholders[i] = "?"
		params[i] = id
	}
	params[len(requirementIDs)] = revisionUUID

	query := fmt.Sprintf(`
		SELECT 
			id,
			requirement_id,
			revision_uuid,
			entry_date,
			value
		FROM requirement_entries
		WHERE requirement_id IN (%s) AND revision_uuid = ?`,
		strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requirementEntriesMap := make(map[int64]models.RequirementEntry)

	for rows.Next() {
		var requirementEntry models.RequirementEntry
		err := rows.Scan(
			&requirementEntry.ID,
			&requirementEntry.RequirementID,
			&requirementEntry.RevisionUUID,
			&requirementEntry.EntryDate,
			&requirementEntry.Value,
		)
		if err != nil {
			return nil, err
		}
		requirementEntriesMap[requirementEntry.RequirementID] = requirementEntry
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(requirementEntriesMap) == 0 {
		logger.Log.Warn().Msg("Retrieved 0 requirement entries")
	} else {
		logger.Log.Debug().
			Int("count", len(requirementEntriesMap)).
			Msg("Successfully retrieved requirement entries")
	}

	return requirementEntriesMap, nil
}
