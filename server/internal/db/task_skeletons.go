package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/mattn/go-sqlite3"
)

type TaskSkeletons interface {
	CreateTx(ctx context.Context, tx *sql.Tx, task *models.TaskSkeleton) (*models.TaskSkeleton, error)
	GetByID(id int64) (*models.TaskSkeleton, error)
	GetAll(showActive, showArchived bool) ([]*models.TaskSkeleton, error)
}

type taskSkeletons struct {
	db *sql.DB
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

func (r *taskSkeletons) CreateTx(ctx context.Context, tx *sql.Tx, taskSkeleton *models.TaskSkeleton) (*models.TaskSkeleton, error) {
	logger.Log.Debug().
		Int64("owner_id", taskSkeleton.OwnerID).
		Msg("Trying to create new task in db via ctx")

	query := `INSERT INTO task_skeletons (owner_id) VALUES (?)`

	result, err := tx.ExecContext(ctx, query, taskSkeleton.OwnerID)
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

	createdTaskSkeleton, err := r.getByIDTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	logger.Log.Info().
		Int64("task_skeleton_id", createdTaskSkeleton.ID).
		Int64("owner_id", createdTaskSkeleton.OwnerID).
		Msg("Created new task skeleton successfully")
	return createdTaskSkeleton, nil
}

func (r *taskSkeletons) getByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*models.TaskSkeleton, error) {
	logger.Log.Debug().
		Int64("id", id).
		Msg("Trying to find the task skeleton in the db via ctx")

	query := `SELECT id, owner_id, status 
	FROM task_skeletons 
	WHERE id = ?`

	var task models.TaskSkeleton

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.OwnerID,
		&task.Status,
	)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskSkeletons) GetByID(id int64) (*models.TaskSkeleton, error) {

	var task models.TaskSkeleton
	logger.Log.Debug().
		Int64("id", id).
		Msg("Trying to find the task skeleton in the db")

	query := `SELECT id, owner_id, status
	FROM task_skeletons
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&task.ID,
		&task.OwnerID,
		&task.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Warn().
			Int64("taskID", id).
			Msg("No task was found")
		return nil, fmt.Errorf("task %d not found: %w", id, err)
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskSkeletons) GetAll(showActive, showArchived bool) ([]*models.TaskSkeleton, error) {
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

	rows, err := r.db.Query(query, args...)
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

	return tasks, nil
}
