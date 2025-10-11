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

type TaskSkeletons interface {
	WithTx(tx *sqlx.Tx) TaskSkeletons

	Create(ctx context.Context, task *models.TaskSkeleton) (*models.TaskSkeleton, error)

	Get(ctx context.Context, uc *models.UserContext) *TaskSkeletonsQuery
}

type taskSkeletons struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

var _ TaskSkeletons = (*taskSkeletons)(nil)

func InitTaskSkeletons(db *sqlx.DB) (TaskSkeletons, error) {
	repo := &taskSkeletons{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *taskSkeletons) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS task_skeletons (
		id              INTEGER NOT NULL PRIMARY KEY,
		owner_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		status          VARCHAR(255) NOT NULL DEFAULT 'active' CHECK(status IN ('archived', 'active'))
    )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskSkeletons) WithTx(tx *sqlx.Tx) TaskSkeletons {
	return &taskSkeletons{
		db: r.db,
		tx: tx,
	}
}

func (r *taskSkeletons) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

// ---------------- //
// INSERT FUNCTIONS //
// ---------------- //

func (r *taskSkeletons) Create(ctx context.Context, taskSkeleton *models.TaskSkeleton) (*models.TaskSkeleton, error) {
	logger.Log.Debug().
		Int64("owner_id", taskSkeleton.OwnerID).
		Msg("Trying to create new task skeleton in db")

	query := `
        INSERT INTO task_skeletons (owner_id, status)
        VALUES (?, ?)
        RETURNING id, owner_id, status
    `

	executor := r.getExecutor()

	var createdTaskSkeleton models.TaskSkeleton
	err := executor.GetContext(ctx, &createdTaskSkeleton, query, taskSkeleton.OwnerID, taskSkeleton.Status)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, apperrors.ErrDuplicate
		}
		return nil, err
	}

	logger.Log.Info().
		Int64("task_skeleton_id", createdTaskSkeleton.ID).
		Int64("owner_id", createdTaskSkeleton.OwnerID).
		Msg("Created new task skeleton successfully")
	return &createdTaskSkeleton, nil
}

// ------------- //
// GET FUNCTIONS //
// ------------- //

type TaskSkeletonsParams struct {
	IDs      []int64
	OwnerIDs []int64
	Status   []string
}

type TaskSkeletonsQuery struct {
	repo   *taskSkeletons
	uc     *models.UserContext
	params *TaskSkeletonsParams
	ctx    context.Context
}

func (r *taskSkeletons) Get(ctx context.Context, uc *models.UserContext) *TaskSkeletonsQuery {
	return &TaskSkeletonsQuery{
		repo:   r,
		uc:     uc,
		params: &TaskSkeletonsParams{},
		ctx:    ctx,
	}
}

func (q *TaskSkeletonsQuery) WithIDs(ids ...any) *TaskSkeletonsQuery {
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

func (q *TaskSkeletonsQuery) WithOwnerIDs(ids ...any) *TaskSkeletonsQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case int64:
			q.params.OwnerIDs = append(q.params.OwnerIDs, v)
		case []int64:
			q.params.OwnerIDs = append(q.params.OwnerIDs, v...)
		}
	}
	return q
}

func (q *TaskSkeletonsQuery) WithStatus(statuses ...any) *TaskSkeletonsQuery {
	for _, s := range statuses {
		switch v := s.(type) {
		case string:
			q.params.Status = append(q.params.Status, v)
		case []string:
			q.params.Status = append(q.params.Status, v...)
		}
	}
	return q
}

func (r *taskSkeletons) BuildQuery(params *TaskSkeletonsParams, user *models.UserContext) (string, []any, error) {
	var whereClauses []string
	var args []any
	var err error

	// Always force owner filter for non-admin users
	if user.Role != models.RoleAdmin {
		params.OwnerIDs = []int64{user.ID}
	}

	whereClauses, args, err = InQuery(whereClauses, args, "id", toAnySlice(params.IDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "owner_id", toAnySlice(params.OwnerIDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "status", toAnySlice(params.Status))
	if err != nil {
		return "", nil, err
	}

	query := "SELECT id, owner_id, status FROM task_skeletons"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

func (q *TaskSkeletonsQuery) Send(ctx context.Context) ([]models.TaskSkeleton, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Interface("query", query).
		Interface("args", args).
		Msg("Executing task skeletons query")

	var skeletons []models.TaskSkeleton
	if err := q.repo.db.SelectContext(ctx, &skeletons, query, args...); err != nil {
		return nil, fmt.Errorf("failed to query task skeletons: %w", err)
	}
	return skeletons, nil
}
