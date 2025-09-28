package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type RequirementEntries interface {
	WithTx(tx *sqlx.Tx) RequirementEntries

	Upsert(ctx context.Context, entry *models.RequirementEntry) (*models.RequirementEntry, error)

	Find(ctx context.Context, user *models.UserContext, opts ...QueryOption) ([]models.RequirementEntry, error)
}

type requirementEntries struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

var _ RequirementEntries = (*requirementEntries)(nil)

func InitRequirementEntries(db *sqlx.DB) (RequirementEntries, error) {

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

func (r *requirementEntries) WithTx(tx *sqlx.Tx) RequirementEntries {
	return &requirementEntries{
		db: r.db,
		tx: tx,
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

	var updatedEntry models.RequirementEntry
	err := executor.GetContext(
		ctx,
		&updatedEntry,
		query,
		entry.RequirementID,
		entry.RevisionUUID,
		entry.EntryDate,
		entry.Value,
	)
	if err != nil {
		return nil, err
	}

	return &updatedEntry, nil
}

// ------------- //
// GET FUNCTIONS //
// ------------- //
type queryParams struct {
	entryIDs       []int64
	requirementIDs []int64
	dates          []time.Time
	startDate      *time.Time
	endDate        *time.Time
}

type QueryOption func(*queryParams)

func WithEntryIDs(ids ...any) QueryOption {
	return func(q *queryParams) {
		for _, id := range ids {
			switch v := id.(type) {
			case int64:
				q.entryIDs = append(q.entryIDs, v)
			case []int64:
				q.entryIDs = append(q.entryIDs, v...)
			}
		}
	}
}

func WithRequirementIDs(ids ...any) QueryOption {
	return func(q *queryParams) {
		for _, id := range ids {
			switch v := id.(type) {
			case int64:
				q.requirementIDs = append(q.requirementIDs, v)
			case []int64:
				q.requirementIDs = append(q.requirementIDs, v...)
			}
		}
	}
}

func WithDates(dates ...any) QueryOption {
	return func(q *queryParams) {
		for _, d := range dates {
			switch v := d.(type) {
			case time.Time:
				q.dates = append(q.dates, v)
			case []time.Time:
				q.dates = append(q.dates, v...)
			}
		}
	}
}

func WithStartDate(start time.Time) QueryOption {
	return func(q *queryParams) {
		q.startDate = &start
	}
}

func WithEndDate(end time.Time) QueryOption {
	return func(q *queryParams) {
		q.endDate = &end
	}
}

func (r *requirementEntries) buildQuery(params *queryParams, user *models.UserContext) (string, []any, error) {
	var whereClauses []string
	var args []any
	var err error

	switch user.Role {
	case models.RoleAdmin:
		// no filter
	case models.RoleUser:
		// no filter
	}

	// Defaults
	if params.startDate == nil {
		today := time.Now().Truncate(24 * time.Hour)
		params.startDate = &today
	}
	if params.endDate == nil {
		today := time.Now().Truncate(24 * time.Hour)
		params.endDate = &today
	}

	// IN clauses
	whereClauses, args, err = InQuery(whereClauses, args, "entry_id", toAnySlice(params.entryIDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "requirement_id", toAnySlice(params.requirementIDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "entry_date", toAnySlice(params.dates))
	if err != nil {
		return "", nil, err
	}

	// BETWEEN clause
	whereClauses, args = BetweenQuery(whereClauses, args, "entry_date", params.startDate, params.endDate)

	// Build final query
	query := "SELECT * FROM requirement_entries"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

func (r *requirementEntries) DebugBuildQuery(params *queryParams, user *models.UserContext) (string, []interface{}, error) {
	query, args, err := r.buildQuery(params, user)
	return query, args, err
}

func (r *requirementEntries) Find(ctx context.Context, user *models.UserContext, opts ...QueryOption) ([]models.RequirementEntry, error) {
	// Start with empty params
	params := &queryParams{}

	// Apply all options
	for _, opt := range opts {
		opt(params)
	}

	// Build and execute query
	query, args, err := r.buildQuery(params, user)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Interface("query", query).
		Interface("args", args).
		Msg("Executing requirement entries query")

	var entries []models.RequirementEntry
	err = r.db.SelectContext(ctx, &entries, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query requirement entries: %w", err)
	}

	return entries, nil
}
