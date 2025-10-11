package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type RequirementSkeletons interface {
	WithTx(tx *sqlx.Tx) RequirementEntries

	Create(ctx context.Context, requirementSkeleton *models.RequirementSkeleton) (*models.RequirementSkeleton, error)

	Get(ctx context.Context, uc *models.UserContext) *RequirementSkeletonsQuery
}

type requirementSkeletons struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

var _ RequirementSkeletons = (*requirementSkeletons)(nil)

func InitRequirementSkeletons(db *sqlx.DB) (RequirementSkeletons, error) {

	repo := &requirementSkeletons{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *requirementSkeletons) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_skeletons (
		id INTEGER PRIMARY KEY AUTOINCREMENT REFERENCES task_skeletons(id) ON DELETE CASCADE,
		task_id INTEGER NOT NULL
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *requirementSkeletons) WithTx(tx *sqlx.Tx) RequirementEntries {
	return &requirementEntries{
		db: r.db,
		tx: tx,
	}
}

func (r *requirementSkeletons) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// ---------------- //
// INSERT FUNCTIONS //
// ---------------- //

func (r *requirementSkeletons) Create(ctx context.Context, requirementSkeleton *models.RequirementSkeleton) (*models.RequirementSkeleton, error) {
	logger.Log.Debug().
		Int64("task_id", requirementSkeleton.TaskID).
		Msg("Trying to create new requirement skeleton in db")

	executor := r.getExecutor()

	query := `
        INSERT INTO requirement_skeletons (
            task_id
        )
        VALUES (?)
        RETURNING id, task_id`

	var createdRequirementSkeleton models.RequirementSkeleton
	err := executor.GetContext(
		ctx,
		&createdRequirementSkeleton,
		query,
		requirementSkeleton.TaskID,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, apperrors.ErrDuplicate
			}
		}
		return nil, err
	}

	logger.Log.Info().
		Int64("requirement_skeleton_id", createdRequirementSkeleton.ID).
		Msg("Created new requirement skeleton successfully")

	return &createdRequirementSkeleton, nil
}

// ------------- //
// GET FUNCTIONS //
// ------------- //

type RequirementSkeletonsParams struct {
	IDs     []int64
	TaskIDs []int64
}

type RequirementSkeletonsQuery struct {
	repo   *requirementSkeletons
	uc     *models.UserContext
	params *RequirementSkeletonsParams
	ctx    context.Context
}

func (r *requirementSkeletons) Get(ctx context.Context, uc *models.UserContext) *RequirementSkeletonsQuery {
	return &RequirementSkeletonsQuery{
		repo:   r,
		uc:     uc,
		params: &RequirementSkeletonsParams{},
		ctx:    ctx,
	}
}

func (q *RequirementSkeletonsQuery) WithIDs(ids ...any) *RequirementSkeletonsQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case int64:
			q.params.IDs = append(q.params.IDs, v)
		case []int64:
			q.params.IDs = append(q.params.IDs, v...)
		}
	}
	return q
}

func (q *RequirementSkeletonsQuery) WithTaskIDs(ids ...any) *RequirementSkeletonsQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case int64:
			q.params.TaskIDs = append(q.params.TaskIDs, v)
		case []int64:
			q.params.TaskIDs = append(q.params.TaskIDs, v...)
		}
	}
	return q
}

func (r *requirementSkeletons) BuildQuery(params *RequirementSkeletonsParams, user *models.UserContext) (string, []any, error) {
	var whereClauses []string
	var args []any
	var err error

	switch user.Role {
	case models.RoleAdmin:
		// no filter
	case models.RoleUser:
		// no filter
	}

	whereClauses, args, err = InQuery(whereClauses, args, "id", toAnySlice(params.IDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "task_id", toAnySlice(params.TaskIDs))
	if err != nil {
		return "", nil, err
	}

	query := "SELECT id, task_id FROM requirement_skeletons"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

func (q *RequirementSkeletonsQuery) All() ([]models.RequirementSkeleton, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Interface("query", query).
		Interface("args", args).
		Msg("Executing requirement skeletons query")

	var skeletons []models.RequirementSkeleton
	if err := q.repo.db.SelectContext(q.ctx, &skeletons, query, args...); err != nil {
		return nil, fmt.Errorf("failed to query requirement skeletons: %w", err)
	}
	return skeletons, nil
}

func (q *RequirementSkeletonsQuery) First() (*models.RequirementSkeleton, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	query += " LIMIT 1"

	var skeleton models.RequirementSkeleton
	if err := q.repo.db.GetContext(q.ctx, &skeleton, query, args...); err != nil {
		return nil, err
	}
	return &skeleton, nil
}
