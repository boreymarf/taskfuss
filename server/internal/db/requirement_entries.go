package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
)

type RequirementEntries interface {
	WithTx(tx *sql.Tx) RequirementEntries

	Upsert(ctx context.Context, entry *models.RequirementEntry) (*models.RequirementEntry, error)

	GetByEntryID(ctx context.Context, entryID int64) (*models.RequirementEntry, error)
	GetByEntryIDs(ctx context.Context, entryIDs []int64) (map[int64]models.RequirementEntry, error)

	GetByRequirementID(ctx context.Context, requirementID int64, entryDate time.Time) (*models.RequirementEntry, error)
	GetByRequirementIDs(ctx context.Context, requirementIDs []int64, entryDates []time.Time) (map[int64]models.RequirementEntry, error)

	GetByRequirementIDInDateRange(ctx context.Context, requirementID int64, startDate, endDate time.Time) (*models.RequirementEntry, error)
	GetByRequirementIDsInDateRange(ctx context.Context, requirementIDs []int64, startDate, endDate time.Time) (map[int64]models.RequirementEntry, error)
}

type requirementEntries struct {
	db  *sql.DB
	tx  *sql.Tx
	ctx context.Context
}

var _ RequirementEntries = (*requirementEntries)(nil)

func InitRequirementEntries(db *sql.DB) (RequirementEntries, error) {

	repo := &requirementEntries{db: db, tx: nil}

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

func (r *requirementEntries) WithTx(tx *sql.Tx) RequirementEntries {
	return &requirementEntries{
		db:  r.db,
		tx:  tx,
		ctx: r.ctx,
	}
}

func (r *requirementEntries) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *requirementEntries) Upsert(ctx context.Context, entry *models.RequirementEntry) (*models.RequirementEntry, error) {
	logger.Log.Debug().
		Str("revisionUUID", entry.RevisionUUID.String()).
		Int64("requirementID", entry.RequirementID).
		Msg("Trying to upsert a requirement entry in db")

	// Use getExecutor() instead of hardcoded r.db
	executor := r.getExecutor()

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

	row := executor.QueryRowContext(
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
		return nil, err
	}

	return &updatedEntry, nil
}

func (r *requirementEntries) GetByEntryID(ctx context.Context, entryID int64) (*models.RequirementEntry, error) {
	if entryID == 0 {
		return nil, fmt.Errorf("entry ID cannot be zero")
	}

	entries, err := r.internalGetByEntryIDs(ctx, []int64{entryID})
	if err != nil {
		return nil, err
	}

	if entry, exists := entries[entryID]; exists {
		return &entry, nil
	}

	return nil, fmt.Errorf("requirement entry not found for entry ID %d", entryID)
}

func (r *requirementEntries) GetByEntryIDs(ctx context.Context, entryIDs []int64) (map[int64]models.RequirementEntry, error) {
	return r.internalGetByEntryIDs(ctx, entryIDs)
}

func (r *requirementEntries) internalGetByEntryIDs(
	ctx context.Context,
	entryIDs []int64,
) (map[int64]models.RequirementEntry, error) {

	logger.Log.Debug().
		Interface("entryIDs", entryIDs).
		Msg("Trying to get requirement entries by entry IDs")

	if len(entryIDs) == 0 {
		logger.Log.Warn().Msg("empty entry IDs slice passed to GetByEntryIDs")
		return map[int64]models.RequirementEntry{}, nil
	}

	executor := r.getExecutor()

	placeholders := make([]string, len(entryIDs))
	params := make([]any, len(entryIDs))

	for i, id := range entryIDs {
		placeholders[i] = "?"
		params[i] = id
	}

	query := fmt.Sprintf(`
        SELECT 
            id,
            requirement_id,
            revision_uuid,
            entry_date,
            value
        FROM requirement_entries
        WHERE id IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := executor.QueryContext(ctx, query, params...)
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
		requirementEntriesMap[requirementEntry.ID] = requirementEntry
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(requirementEntriesMap) == 0 {
		logger.Log.Warn().Msg("Retrieved 0 requirement entries by entry IDs")
	} else {
		logger.Log.Debug().
			Int("count", len(requirementEntriesMap)).
			Msg("Successfully retrieved requirement entries by entry IDs")
	}

	return requirementEntriesMap, nil
}

func (r *requirementEntries) GetByRequirementIDs(ctx context.Context, requirementIDs []int64, entryDates []time.Time) (map[int64]models.RequirementEntry, error) {
	return r.internalGetByRequirementIDs(ctx, requirementIDs, entryDates)
}

func (r *requirementEntries) GetByRequirementID(ctx context.Context, requirementID int64, entryDate time.Time) (*models.RequirementEntry, error) {
	if requirementID == 0 {
		return nil, fmt.Errorf("requirement ID cannot be zero")
	}

	entries, err := r.internalGetByRequirementIDs(ctx, []int64{requirementID}, []time.Time{entryDate})
	if err != nil {
		return nil, err
	}

	if entry, exists := entries[requirementID]; exists {
		return &entry, nil
	}

	return nil, fmt.Errorf("requirement entry not found for ID %d and date %s", requirementID, entryDate.Format("2006-01-02"))
}

func (r *requirementEntries) internalGetByRequirementIDs(
	ctx context.Context,
	requirementIDs []int64,
	entryDates []time.Time,
) (map[int64]models.RequirementEntry, error) {

	logger.Log.Debug().
		Interface("requirementIDs", requirementIDs).
		Interface("entryDates", entryDates).
		Msg("Trying to get requirement entries by requirement IDs and entry dates")

	if len(requirementIDs) == 0 {
		logger.Log.Warn().Msg("empty requirement IDs slice passed to GetByRequirementIDsAndEntryDates")
		return map[int64]models.RequirementEntry{}, nil
	}

	if len(entryDates) == 0 {
		logger.Log.Warn().Msg("empty entry dates slice passed to GetByRequirementIDsAndEntryDates")
		return map[int64]models.RequirementEntry{}, nil
	}

	// Use getExecutor() instead of hardcoded executor parameter
	executor := r.getExecutor()

	// Create placeholders and parameters for requirement IDs
	reqPlaceholders := make([]string, len(requirementIDs))
	params := make([]any, 0, len(requirementIDs)+len(entryDates))

	for i, id := range requirementIDs {
		reqPlaceholders[i] = "?"
		params = append(params, id)
	}

	// Create placeholders and parameters for entry dates
	datePlaceholders := make([]string, len(entryDates))
	for i, date := range entryDates {
		datePlaceholders[i] = "?"
		params = append(params, date)
	}

	query := fmt.Sprintf(`
        SELECT 
            id,
            requirement_id,
            revision_uuid,
            entry_date,
            value
        FROM requirement_entries
        WHERE requirement_id IN (%s) AND entry_date IN (%s)`,
		strings.Join(reqPlaceholders, ","),
		strings.Join(datePlaceholders, ","))

	rows, err := executor.QueryContext(ctx, query, params...)
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

func (r *requirementEntries) GetByRequirementIDsInDateRange(ctx context.Context, requirementIDs []int64, startDate, endDate time.Time) (map[int64]models.RequirementEntry, error) {
	return r.internalGetByRequirementIDsInDateRange(ctx, requirementIDs, startDate, endDate)
}

func (r *requirementEntries) GetByRequirementIDInDateRange(ctx context.Context, requirementID int64, startDate, endDate time.Time) (*models.RequirementEntry, error) {
	if requirementID == 0 {
		return nil, fmt.Errorf("requirement ID cannot be zero")
	}

	entries, err := r.internalGetByRequirementIDsInDateRange(ctx, []int64{requirementID}, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if entry, exists := entries[requirementID]; exists {
		return &entry, nil
	}

	return nil, fmt.Errorf("requirement entry not found for ID %d and date range %s to %s",
		requirementID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
}

func (r *requirementEntries) internalGetByRequirementIDsInDateRange(
	ctx context.Context,
	requirementIDs []int64,
	startDate, endDate time.Time,
) (map[int64]models.RequirementEntry, error) {

	logger.Log.Debug().
		Interface("requirementIDs", requirementIDs).
		Time("startDate", startDate).
		Time("endDate", endDate).
		Msg("Trying to get requirement entries by requirement IDs and date range")

	if len(requirementIDs) == 0 {
		logger.Log.Warn().Msg("empty requirement IDs slice passed to GetByRequirementIDsInDateRange")
		return map[int64]models.RequirementEntry{}, nil
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	// Use getExecutor() instead of hardcoded executor parameter
	executor := r.getExecutor()

	// Create placeholders and parameters for requirement IDs
	reqPlaceholders := make([]string, len(requirementIDs))
	params := make([]any, 0, len(requirementIDs)+2)

	for i, id := range requirementIDs {
		reqPlaceholders[i] = "?"
		params = append(params, id)
	}

	// Add date range parameters
	params = append(params, startDate, endDate)

	query := fmt.Sprintf(`
        SELECT 
            id,
            requirement_id,
            revision_uuid,
            entry_date,
            value
        FROM requirement_entries
        WHERE requirement_id IN (%s) AND entry_date BETWEEN ? AND ?`,
		strings.Join(reqPlaceholders, ","))

	rows, err := executor.QueryContext(ctx, query, params...)
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
