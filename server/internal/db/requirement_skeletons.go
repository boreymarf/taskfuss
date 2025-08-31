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
	"github.com/rs/zerolog"
)

type RequirementSkeletons interface {
	CreateTx(ctx context.Context, tx *sql.Tx, requirementSkeleton *models.RequirementSkeleton) (*models.RequirementSkeleton, error)
	GetByID(id int64) (*models.TaskSkeleton, error)
	GetByTaskIDs(task_ids []int64) ([]*models.RequirementSkeleton, error)
}

type requirementSkeletons struct {
	db *sql.DB
}

var _ TaskSkeletons = (*taskSkeletons)(nil)

func InitRequirementSkeletons(db *sql.DB) (RequirementSkeletons, error) {

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

func (r *requirementSkeletons) CreateTx(ctx context.Context, tx *sql.Tx, requirementSkeleton *models.RequirementSkeleton) (*models.RequirementSkeleton, error) {
	logger.Log.Debug().
		Int64("task_id", requirementSkeleton.TaskID).
		Msg("Trying to create new requirement skeleton in db via ctx")

	query := `INSERT INTO requirement_skeletons (
		task_id
	) VALUES (?)`

	result, err := tx.ExecContext(ctx, query, requirementSkeleton.TaskID)
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

	createdRequirementSkeleton, err := r.getByIDTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	logger.Log.Info().
		Int64("requirement_skeleton_id", createdRequirementSkeleton.ID).
		Msg("Created new requirement skeleton successfully")
	return createdRequirementSkeleton, nil
}

func (r *requirementSkeletons) getByIDTx(ctx context.Context, tx *sql.Tx, id int64) (*models.RequirementSkeleton, error) {
	logger.Log.Debug().
		Int64("id", id).
		Msg("Trying to find the task skeleton in the db via ctx")

	query := `SELECT 
		id,
		task_id
	FROM requirement_skeletons 
	WHERE id = ?`

	var requirementSkeleton models.RequirementSkeleton

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&requirementSkeleton.ID,
		&requirementSkeleton.TaskID,
	)
	if err != nil {
		return nil, err
	}

	return &requirementSkeleton, nil
}

func (r *requirementSkeletons) GetByID(id int64) (*models.TaskSkeleton, error) {

	var task models.TaskSkeleton
	logger.Log.Debug().
		Int64("id", id).
		Msg("Trying to find the task skeleton in the db")

	query := `SELECT
		id,
		task_id
	FROM requirement_skeletons
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

func (r *requirementSkeletons) GetByTaskIDs(task_ids []int64) ([]*models.RequirementSkeleton, error) {
	if len(task_ids) > 0 {
		arr := zerolog.Arr()
		for _, id := range task_ids {
			arr.Int64(id)
		}
		logger.Log.Debug().
			Array("task_ids", arr).
			Msg("Trying to get all requirement skeletons from the db")
	} else {
		logger.Log.Debug().
			Msg("Task ids list is empty, skipping database query")
	}

	if len(task_ids) == 0 {
		return []*models.RequirementSkeleton{}, nil
	}

	query := `SELECT
		id,
		task_id
	FROM requirement_skeletons
	WHERE task_id IN (` + strings.Repeat("?,", len(task_ids)-1) + "?)"

	args := make([]any, len(task_ids))
	for i, id := range task_ids {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requirements []*models.RequirementSkeleton
	for rows.Next() {
		var req models.RequirementSkeleton
		err := rows.Scan(
			&req.ID,
			&req.TaskID,
		)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requirements, nil
}
