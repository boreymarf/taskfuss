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

	logger.Log.Debug().Msg("taskRepository initialization completed")

	return repo, nil
}

func (r *TaskRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
	id 						INTEGER NOT NULL PRIMARY KEY,
	name 					VARCHAR(255) NOT NULL,
	description 	VARCHAR(255),
	created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
	start_date 		DATETIME DEFAULT CURRENT_TIMESTAMP,
 	end_date 			DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) CreateTask(task *models.Task) error {

	logger.Log.Info().Int64("id", task.ID).Msg("TaskRepository tries to create new task")

	query := `INSERT INTO tasks (name, description, start_date, end_date) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, task.Name, task.Description, task.StartDate, task.EndDate)

	if err != nil {
		// If there's a dublicate task
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return apperrors.ErrDuplicate
			}
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = id

	err = r.GetTaskByID(id, task)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetTaskByID(id int64, task *models.Task) error {

	logger.Log.Debug().Int64("id", id).Msg("taskRepository tries to find user")

	query := `SELECT id, name , description, created_at, updated_at, start_date, end_date  
	FROM tasks 
	WHERE taskId = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&task.ID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Warn().
			Int64("taskID", id).
			Msg("No task was found")
		return fmt.Errorf("task %d not found: %w", id, err)
	} else if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) GetAllTasks() error {
	return nil
}
