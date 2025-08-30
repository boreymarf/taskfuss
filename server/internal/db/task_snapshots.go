package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TaskSnapshots interface {
	CreateTx(ctx context.Context, tx *sql.Tx, taskSnapshot *models.TaskSnapshot, revision_uuid uuid.UUID) (*models.TaskSnapshot, error)
}

type taskSnapshots struct {
	db *sql.DB
}

var _ TaskSnapshots = (*taskSnapshots)(nil)

func InitTaskSnapshots(db *sql.DB) (TaskSnapshots, error) {
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
		task_id INTEGER  NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_current BOOLEAN DEFAULT FALSE,
	PRIMARY KEY (revision_uuid, task_id)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskSnapshots) CreateTx(ctx context.Context, tx *sql.Tx, taskSnapshot *models.TaskSnapshot, revision_uuid uuid.UUID) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Int64("task_id", taskSnapshot.TaskID).
		Msg("Trying to create new task snapshot in db via ctx")

	query := `INSERT INTO task_snapshots (
		revision_uuid,
		task_id,
		title,
		description
	) VALUES (?, ?, ?, ?)`

	result, err := tx.ExecContext(
		ctx,
		query,
		revision_uuid,
		taskSnapshot.TaskID,
		taskSnapshot.Title,
		taskSnapshot.Description,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				logger.Log.Error().
					Str("revision_uuid", revision_uuid.String()).
					Msg("For some reason there's duplicate of the task snapshot!")
				return nil, apperrors.ErrDuplicate
			}
		}
		return nil, err
	}

	// General checks
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no rows affected during insert")
	}

	createdTaskSnapshot, err := r.getByRevisionUUIDTx(ctx, tx, revision_uuid)
	if err != nil {
		return nil, err
	}
	if createdTaskSnapshot == nil {
		return nil, fmt.Errorf("failed to retrieve created task snapshot")
	}

	logger.Log.Debug().
		Str("revision_uuid", createdTaskSnapshot.RevisionUUID.String()).
		Int64("task_id", createdTaskSnapshot.TaskID).
		Msg("Created new task snapshot successfully")
	return createdTaskSnapshot, nil
}

func (r *taskSnapshots) getByRevisionUUIDTx(ctx context.Context, tx *sql.Tx, revision_uuid uuid.UUID) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Str("revision_uuid", revision_uuid.String()).
		Msg("Trying to find the task snapshot in the db via ctx")

	query := `SELECT
		revision_uuid,
		task_id,
		title,
		description,
		created_at,
		is_current
	FROM task_snapshots
	WHERE revision_uuid = ?`

	var task models.TaskSnapshot

	err := tx.QueryRowContext(ctx, query, revision_uuid).Scan(
		&task.RevisionUUID,
		&task.TaskID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.IsCurrent,
	)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
