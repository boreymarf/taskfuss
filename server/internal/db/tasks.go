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

type TaskRepository struct {
	db *sql.DB
}

func InitTaskRepository(db *sql.DB) (*TaskRepository, error) {

	repo := &TaskRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("Repository initialization completed")

	return repo, nil
}

func (r *TaskRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
		id              INTEGER NOT NULL PRIMARY KEY,
		owner_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		status          VARCHAR(255) NOT NULL DEFAULT 'active' CHECK(status IN ('archived', 'active')),
		title           VARCHAR(255) NOT NULL,
		description     VARCHAR(255),
		created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
		start_date      DATETIME DEFAULT CURRENT_TIMESTAMP,
		end_date        DATETIME DEFAULT NULL
    )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) CreateTask(task models.Task) (*models.Task, error) {
	logger.Log.Debug().
		Str("title", task.Title).
		Int64("owner_id", task.OwnerID).
		Msg("Trying to Create new task to the db...")

	query := `INSERT INTO tasks (owner_id, title, description) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, task.OwnerID, task.Title, task.Description)

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

	createdTask, err := r.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	return &createdTask, nil
}

func (r *TaskRepository) GetTaskByID(id int64) (models.Task, error) {

	var task models.Task
	logger.Log.Debug().Int64("id", id).Msg("taskRepository tries to find task")

	query := `SELECT id, owner_id, title, description, created_at, updated_at, start_date, end_date, status
	FROM tasks
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&task.ID,
		&task.OwnerID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.StartDate,
		&task.EndDate,
		&task.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Warn().
			Int64("taskID", id).
			Msg("No task was found")
		return models.Task{}, fmt.Errorf("task %d not found: %w", id, err)
	} else if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

type GetAllTasksOptions struct {
	DetailLevel  string
	ShowArchived bool
	ShowActive   bool
	UserID       int64
}

func (r *TaskRepository) GetAllTasks(opts *GetAllTasksOptions) ([]models.Task, error) {

	query := `SELECT 
		id,
		owner_id,
		title,
		description,
		created_at,
		updated_at,
		start_date,
		end_date,
		status
	FROM 
		tasks
	WHERE 
		owner_id = ?
  AND (
    (status = 'archived' AND ?) 
    OR 
    (status = 'active' AND ?)
  )`

	rows, err := r.db.Query(query,
		opts.UserID,
		opts.ShowArchived,
		opts.ShowActive,
	)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to execute query")
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.OwnerID,
			&task.Title,
			&task.Description,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.StartDate,
			&task.EndDate,
			&task.Status,
		)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan task row")
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		filteredTask := models.Task{
			ID:      task.ID,
			OwnerID: task.OwnerID,
			Title:   task.Title,
			Status:  task.Status,
		}

		if opts.DetailLevel == "basic" || opts.DetailLevel == "full" {
			filteredTask.StartDate = task.StartDate
			filteredTask.EndDate = task.EndDate
		}

		if opts.DetailLevel == "full" {
			filteredTask.Description = task.Description
			filteredTask.CreatedAt = task.CreatedAt
			filteredTask.UpdatedAt = task.UpdatedAt
		}

		tasks = append(tasks, filteredTask)

	}

	return tasks, nil
}
