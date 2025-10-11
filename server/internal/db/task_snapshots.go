package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type TaskSnapshots interface {
	WithTx(tx *sqlx.Tx) RequirementEntries

	Create(ctx context.Context, taskSnapshot *models.TaskSnapshot, revisionUUID uuid.UUID, taskID int64) (*models.TaskSnapshot, error)

	SetCurrentRevisionTx(ctx context.Context, taskID int64, revisionUUID uuid.UUID) error
}

type taskSnapshots struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

var _ TaskSnapshots = (*taskSnapshots)(nil)

func InitTaskSnapshots(db *sqlx.DB) (TaskSnapshots, error) {
	repo := &taskSnapshots{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *taskSnapshots) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS task_snapshots (
		revision_uuid TEXT NOT NULL,
		skeleton_id INTEGER  NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		description TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_current BOOLEAN DEFAULT FALSE,
	PRIMARY KEY (revision_uuid, skeleton_id)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskSnapshots) WithTx(tx *sqlx.Tx) RequirementEntries {
	return &requirementEntries{
		db: r.db,
		tx: tx,
	}
}

func (r *taskSnapshots) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// ---------------- //
// INSERT FUNCTIONS //
// ---------------- //

func (r *taskSnapshots) Create(ctx context.Context, taskSnapshot *models.TaskSnapshot, revisionUUID uuid.UUID, taskID int64) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Str("revision_uuid", revisionUUID.String()).
		Int64("skeleton_id", taskSnapshot.SkeletonID).
		Msg("Trying to create new task snapshot in db")

	executor := r.getExecutor()

	query := `
        INSERT INTO task_snapshots (
            revision_uuid,
            skeleton_id,
            title,
            description
        )
        VALUES (?, ?, ?, ?)
        RETURNING revision_uuid, skeleton_id, title, description`

	var createdTaskSnapshot models.TaskSnapshot
	err := executor.GetContext(
		ctx,
		&createdTaskSnapshot,
		query,
		revisionUUID,
		taskID,
		taskSnapshot.Title,
		taskSnapshot.Description,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			logger.Log.Error().
				Str("revision_uuid", revisionUUID.String()).
				Msg("Duplicate task snapshot detected")
			return nil, apperrors.ErrDuplicate
		}
		return nil, err
	}

	logger.Log.Debug().
		Str("revision_uuid", createdTaskSnapshot.RevisionUUID.String()).
		Int64("skeleton_id", createdTaskSnapshot.SkeletonID).
		Msg("Created new task snapshot successfully")

	return &createdTaskSnapshot, nil
}

// ------------- //
// GET FUNCTIONS //
// ------------- //

type TaskSnapshotsParams struct {
	RevisionUUIDs []uuid.UUID
	SkeletonIDs   []int64
	CreatedAts    []time.Time
	IsCurrents    []bool
}

type TaskSnapshotsQuery struct {
	repo   *taskSnapshots
	uc     *models.UserContext
	params *TaskSnapshotsParams
	ctx    context.Context
}

func (r *taskSnapshots) Get(ctx context.Context, uc *models.UserContext) *TaskSnapshotsQuery {
	return &TaskSnapshotsQuery{
		repo:   r,
		uc:     uc,
		params: &TaskSnapshotsParams{},
		ctx:    ctx,
	}
}

func (q *TaskSnapshotsQuery) WithRevisionUUIDs(ids ...any) *TaskSnapshotsQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case uuid.UUID:
			q.params.RevisionUUIDs = append(q.params.RevisionUUIDs, v)
		case []uuid.UUID:
			q.params.RevisionUUIDs = append(q.params.RevisionUUIDs, v...)
		}
	}
	return q
}

func (q *TaskSnapshotsQuery) WithSkeletonIDs(ids ...any) *TaskSnapshotsQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case int64:
			q.params.SkeletonIDs = append(q.params.SkeletonIDs, v)
		case []int64:
			q.params.SkeletonIDs = append(q.params.SkeletonIDs, v...)
		}
	}
	return q
}

func (q *TaskSnapshotsQuery) WithCreatedAts(times ...any) *TaskSnapshotsQuery {
	for _, t := range times {
		switch v := t.(type) {
		case time.Time:
			q.params.CreatedAts = append(q.params.CreatedAts, v)
		case []time.Time:
			q.params.CreatedAts = append(q.params.CreatedAts, v...)
		}
	}
	return q
}

func (q *TaskSnapshotsQuery) WithIsCurrents(flags ...any) *TaskSnapshotsQuery {
	for _, f := range flags {
		switch v := f.(type) {
		case bool:
			q.params.IsCurrents = append(q.params.IsCurrents, v)
		case []bool:
			q.params.IsCurrents = append(q.params.IsCurrents, v...)
		}
	}
	return q
}

func (r *taskSnapshots) BuildQuery(params *TaskSnapshotsParams, user *models.UserContext) (string, []any, error) {
	var whereClauses []string
	var args []any
	var err error

	switch user.Role {
	case models.RoleAdmin:
		// no filter
	case models.RoleUser:
		// no filter
	}

	whereClauses, args, err = InQuery(whereClauses, args, "revision_uuid", toAnySlice(params.RevisionUUIDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "skeleton_id", toAnySlice(params.SkeletonIDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "created_at", toAnySlice(params.CreatedAts))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "is_current", toAnySlice(params.IsCurrents))
	if err != nil {
		return "", nil, err
	}

	query := "SELECT revision_uuid, skeleton_id, title, description, created_at, is_current FROM task_snapshots"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

func (q *TaskSnapshotsQuery) Send(ctx context.Context) ([]models.TaskSnapshot, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Interface("query", query).
		Interface("args", args).
		Msg("Executing task snapshots query")

	var snapshots []models.TaskSnapshot
	if err := q.repo.db.SelectContext(ctx, &snapshots, query, args...); err != nil {
		return nil, fmt.Errorf("failed to query task snapshots: %w", err)
	}
	return snapshots, nil
}

// ---------------- //
// UPDATE FUNCTIONS //
// ---------------- //

func (r *taskSnapshots) SetCurrentRevisionTx(ctx context.Context, taskID int64, revisionUUID uuid.UUID) error {
	logger.Log.Debug().
		Int64("skeleton_id", taskID).
		Str("revision_uuid", revisionUUID.String()).
		Msg("Trying to set current revision")

	resetQuery := `
		UPDATE task_snapshots 
		SET is_current = FALSE 
		WHERE skeleton_id = ? AND is_current = TRUE`

	executor := r.getExecutor()

	_, err := executor.ExecContext(ctx, resetQuery, taskID)
	if err != nil {
		return err
	}

	setQuery := `
		UPDATE task_snapshots 
		SET is_current = TRUE 
		WHERE skeleton_id = ? AND revision_uuid = ?`

	_, err = executor.ExecContext(ctx, setQuery, taskID, revisionUUID)
	return err
}
