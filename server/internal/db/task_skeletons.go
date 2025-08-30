package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/mattn/go-sqlite3"
)

type TaskSkeletons interface {
	Create(task *models.TaskSkeleton) (*models.TaskSkeleton, error)
	GetByID(id int64) (*models.TaskSkeleton, error)
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

func (r *taskSkeletons) Create(task *models.TaskSkeleton) (*models.TaskSkeleton, error) {
	logger.Log.Debug().
		Int64("owner_id", task.OwnerID).
		Msg("Trying to create new task in db")

	query := `INSERT INTO tasks_skeletons (owner_id) VALUES (?)`

	result, err := r.db.Exec(query, task.OwnerID)

	if err != nil {
		var sqliteErr sqlite3.Error
		// If there's a dublicate task
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

	createdTask, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	logger.Log.Info().
		Int64("task_skeleton_id", createdTask.ID).
		Int64("owner_id", createdTask.OwnerID).
		Msg("Created new task skeleton successfully")
	return createdTask, nil
}

func (r *taskSkeletons) GetByID(id int64) (*models.TaskSkeleton, error) {

	var task models.TaskSkeleton
	logger.Log.Debug().
		Int64("id", id).
		Msg("Trying to find the task skeleton in the db")

	query := `SELECT id, owner_id, status
	FROM tasks_skeletons
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
