package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/mattn/go-sqlite3"
)

type TaskSkeletons interface {
	WithTx(tx *sql.Tx) TaskSkeletons

	Create(ctx context.Context, task *models.TaskSkeleton) (*models.TaskSkeleton, error)

	GetByID(ctx context.Context, id int64) (*models.TaskSkeleton, error)
	GetByIDs(ctx context.Context, ids []int64) (map[int64]models.TaskSkeleton, error)
	GetAll(ctx context.Context, showActive, showArchived bool) ([]*models.TaskSkeleton, error)
}

type taskSkeletons struct {
	db  *sql.DB
	tx  *sql.Tx
	ctx context.Context
}

var _ TaskSkeletons = (*taskSkeletons)(nil)

func InitTaskSkeletons(db *sql.DB) (TaskSkeletons, error) {
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

func (r *taskSkeletons) WithTx(tx *sql.Tx) TaskSkeletons {
	return &taskSkeletons{
		db:  r.db,
		tx:  tx,
		ctx: r.ctx,
	}
}

func (r *taskSkeletons) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *taskSkeletons) Create(ctx context.Context, taskSkeleton *models.TaskSkeleton) (*models.TaskSkeleton, error) {
	logger.Log.Debug().
		Int64("owner_id", taskSkeleton.OwnerID).
		Msg("Trying to create new task skeleton in db")

	query := `INSERT INTO task_skeletons (owner_id, status) VALUES (?, ?)`

	executor := r.getExecutor()

	result, err := executor.ExecContext(ctx, query, taskSkeleton.OwnerID, taskSkeleton.Status)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, apperrors.ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Use the same instance to get the created task skeleton
	createdTaskSkeleton, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Log.Info().
		Int64("task_skeleton_id", createdTaskSkeleton.ID).
		Int64("owner_id", createdTaskSkeleton.ID).
		Msg("Created new task skeleton successfully")
	return createdTaskSkeleton, nil
}

func (r *taskSkeletons) GetByID(ctx context.Context, id int64) (*models.TaskSkeleton, error) {
	if id == 0 {
		return nil, fmt.Errorf("task ID cannot be zero")
	}

	tasks, err := r.internalGetByIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}

	if task, exists := tasks[id]; exists {
		return &task, nil
	}

	return nil, fmt.Errorf("task skeleton not found for ID %d", id)
}

func (r *taskSkeletons) GetByIDs(ctx context.Context, ids []int64) (map[int64]models.TaskSkeleton, error) {
	return r.internalGetByIDs(ctx, ids)
}

func (r *taskSkeletons) internalGetByIDs(
	ctx context.Context,
	ids []int64,
) (map[int64]models.TaskSkeleton, error) {

	logger.Log.Debug().
		Interface("ids", ids).
		Msg("Trying to get task skeletons by IDs")

	if len(ids) == 0 {
		logger.Log.Warn().Msg("empty IDs slice passed to GetByIDs")
		return map[int64]models.TaskSkeleton{}, nil
	}

	executor := r.getExecutor()

	placeholders := make([]string, len(ids))
	params := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		params[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, owner_id, status
		FROM task_skeletons
		WHERE id IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := executor.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskSkeletonsMap := make(map[int64]models.TaskSkeleton)

	for rows.Next() {
		var task models.TaskSkeleton
		err := rows.Scan(
			&task.ID,
			&task.OwnerID,
			&task.Status,
		)
		if err != nil {
			return nil, err
		}
		taskSkeletonsMap[task.ID] = task
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(taskSkeletonsMap) == 0 {
		logger.Log.Warn().Msg("Retrieved 0 task skeletons by IDs")
	} else {
		logger.Log.Debug().
			Int("count", len(taskSkeletonsMap)).
			Msg("Successfully retrieved task skeletons by IDs")
	}

	return taskSkeletonsMap, nil
}

func (r *taskSkeletons) GetAll(ctx context.Context, showActive, showArchived bool) ([]*models.TaskSkeleton, error) {
	logger.Log.Debug().
		Bool("show_active", showActive).
		Bool("show_archived", showArchived).
		Msg("Trying to get all task skeletons from the db")

	query := `SELECT id, owner_id, status
	FROM task_skeletons
	WHERE 1=1`

	args := []any{}

	if showActive {
		query += " AND status = ?"
		args = append(args, "active")
	}

	if showArchived {
		query += " AND status = ?"
		args = append(args, "archived")
	}

	executor := r.getExecutor()

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.TaskSkeleton
	for rows.Next() {
		var task models.TaskSkeleton
		err := rows.Scan(
			&task.ID,
			&task.OwnerID,
			&task.Status,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Int("count", len(tasks)).
		Msg("Successfully retrieved task skeletons")

	return tasks, nil
}
