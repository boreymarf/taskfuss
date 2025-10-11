package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type RequirementSnapshots interface {
	WithTx(tx *sqlx.Tx) RequirementEntries

	Create(ctx context.Context, requirementSnapshot *models.RequirementSnapshot, revisionUUID uuid.UUID) (*models.RequirementSnapshot, error)

	Get(ctx context.Context, uc *models.UserContext) *RequirementSnapshotsQuery
}

type requirementSnapshots struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

var _ RequirementSnapshots = (*requirementSnapshots)(nil)

func InitRequirementSnapshots(db *sqlx.DB) (*requirementSnapshots, error) {

	repo := &requirementSnapshots{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *requirementSnapshots) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_snapshots (
		revision_uuid TEXT NOT NULL REFERENCES task_snapshots(revision_uuid) ON DELETE CASCADE,
		skeleton_id INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		parent_id INTEGER DEFAULT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('atom', 'condition')),
		data_type TEXT NOT NULL CHECK (data_type IN ('bool', 'int', 'float', 'duration', 'none')),
		operator TEXT CHECK (operator IN ('or', 'not', 'and', '==', '>=', '<=', '!=', '>', '<')),
		target_value TEXT NOT NULL,
		sort_order INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (revision_uuid, skeleton_id)
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *requirementSnapshots) WithTx(tx *sqlx.Tx) RequirementEntries {
	return &requirementEntries{
		db: r.db,
		tx: tx,
	}
}

func (r *requirementSnapshots) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// ---------------- //
// INSERT FUNCTIONS //
// ---------------- //

func (r *requirementSnapshots) Create(ctx context.Context, requirementSnapshot *models.RequirementSnapshot, revisionUUID uuid.UUID) (*models.RequirementSnapshot, error) {
	logger.Log.Debug().
		Int64("skeleton_id", requirementSnapshot.SkeletonID).
		Msg("Trying to create new requirement snapshot in db")

	executor := r.getExecutor()

	query := `
        INSERT INTO requirement_snapshots (
            revision_uuid,
            skeleton_id,
            parent_id,
            title,
            type,
            data_type,
            operator,
            target_value,
            sort_order
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        RETURNING id, revision_uuid, skeleton_id, parent_id, title, type, data_type, operator, target_value, sort_order
    `

	var createdRequirementSnapshot models.RequirementSnapshot
	err := executor.GetContext(
		ctx,
		&createdRequirementSnapshot,
		query,
		revisionUUID,
		requirementSnapshot.SkeletonID,
		requirementSnapshot.ParentID,
		requirementSnapshot.Title,
		requirementSnapshot.Type,
		requirementSnapshot.DataType,
		requirementSnapshot.Operator,
		requirementSnapshot.TargetValue,
		requirementSnapshot.SortOrder,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			logger.Log.Error().
				Str("revision_uuid", revisionUUID.String()).
				Msg("Duplicate requirement snapshot detected")
			return nil, apperrors.ErrDuplicate
		}
		return nil, err
	}

	logger.Log.Debug().
		Str("revision_uuid", createdRequirementSnapshot.RevisionUUID.String()).
		Int64("skeleton_id", createdRequirementSnapshot.SkeletonID).
		Msg("Created new requirement snapshot successfully")

	return &createdRequirementSnapshot, nil
}

// ------------- //
// GET FUNCTIONS //
// ------------- //

type RequirementSnapshotsParams struct {
	RevisionUUIDs []uuid.UUID
	SkeletonIDs   []int64
	ParentIDs     []int64
	Types         []string
	DataTypes     []string
	Operators     []string
}

type RequirementSnapshotsQuery struct {
	repo   *requirementSnapshots
	uc     *models.UserContext
	params *RequirementSnapshotsParams
	ctx    context.Context
}

func (r *requirementSnapshots) Get(ctx context.Context, uc *models.UserContext) *RequirementSnapshotsQuery {
	return &RequirementSnapshotsQuery{
		repo:   r,
		uc:     uc,
		params: &RequirementSnapshotsParams{},
	}
}

func (q *RequirementSnapshotsQuery) WithRevisionUUIDs(ids ...any) *RequirementSnapshotsQuery {
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

func (q *RequirementSnapshotsQuery) WithSkeletonIDs(ids ...any) *RequirementSnapshotsQuery {
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

func (q *RequirementSnapshotsQuery) WithParentIDs(ids ...any) *RequirementSnapshotsQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case int64:
			q.params.ParentIDs = append(q.params.ParentIDs, v)
		case []int64:
			q.params.ParentIDs = append(q.params.ParentIDs, v...)
		}
	}
	return q
}

func (q *RequirementSnapshotsQuery) WithTypes(types ...any) *RequirementSnapshotsQuery {
	for _, t := range types {
		switch v := t.(type) {
		case string:
			q.params.Types = append(q.params.Types, v)
		case []string:
			q.params.Types = append(q.params.Types, v...)
		}
	}
	return q
}

func (q *RequirementSnapshotsQuery) WithDataTypes(dataTypes ...any) *RequirementSnapshotsQuery {
	for _, dt := range dataTypes {
		switch v := dt.(type) {
		case string:
			q.params.DataTypes = append(q.params.DataTypes, v)
		case []string:
			q.params.DataTypes = append(q.params.DataTypes, v...)
		}
	}
	return q
}

func (q *RequirementSnapshotsQuery) WithOperators(operators ...any) *RequirementSnapshotsQuery {
	for _, op := range operators {
		switch v := op.(type) {
		case string:
			q.params.Operators = append(q.params.Operators, v)
		case []string:
			q.params.Operators = append(q.params.Operators, v...)
		}
	}
	return q
}

func (r *requirementSnapshots) BuildQuery(params *RequirementSnapshotsParams, user *models.UserContext) (string, []any, error) {
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

	whereClauses, args, err = InQuery(whereClauses, args, "parent_id", toAnySlice(params.ParentIDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "type", toAnySlice(params.Types))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "data_type", toAnySlice(params.DataTypes))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "operator", toAnySlice(params.Operators))
	if err != nil {
		return "", nil, err
	}

	query := "SELECT revision_uuid, skeleton_id, parent_id, title, type, data_type, operator, target_value, sort_order FROM requirement_snapshots"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

func (q *RequirementSnapshotsQuery) Send(ctx context.Context) ([]models.RequirementSnapshot, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Interface("query", query).
		Interface("args", args).
		Msg("Executing requirement snapshots query")

	var snapshots []models.RequirementSnapshot
	if err := q.repo.db.SelectContext(ctx, &snapshots, query, args...); err != nil {
		return nil, fmt.Errorf("failed to query requirement snapshots: %w", err)
	}
	return snapshots, nil
}
