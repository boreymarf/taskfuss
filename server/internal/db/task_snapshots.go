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
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type TaskSnapshots interface {
	CreateTx(ctx context.Context, tx *sql.Tx, taskSnapshot *models.TaskSnapshot, revision_uuid uuid.UUID, skeleton_id int64) (*models.TaskSnapshot, error)

	GetByCompositeKey(ctx context.Context, revision_uuid uuid.UUID, skeleton_id int64) (*models.TaskSnapshot, error)
	GetLatest(ctx context.Context, skeleton_id int64) (*models.TaskSnapshot, error)
	GetAllLatest(ctx context.Context, skeleton_ids []int64) ([]*models.TaskSnapshot, error)

	SetCurrentRevisionTx(ctx context.Context, tx *sql.Tx, taskID int64, revisionUUID uuid.UUID) error
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
		skeleton_id INTEGER  NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_current BOOLEAN DEFAULT FALSE,
	PRIMARY KEY (revision_uuid, skeleton_id)
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskSnapshots) CreateTx(
	ctx context.Context,
	tx *sql.Tx,
	taskSnapshot *models.TaskSnapshot,
	revisionUUID uuid.UUID,
	taskID int64,
) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Int64("skeleton_id", taskSnapshot.SkeletonID).
		Msg("Trying to create new task snapshot in db via ctx")

	query := `INSERT INTO task_snapshots (
		revision_uuid,
		skeleton_id,
		title,
		description
	) VALUES (?, ?, ?, ?)`

	result, err := tx.ExecContext(
		ctx,
		query,
		revisionUUID,
		taskID,
		taskSnapshot.Title,
		taskSnapshot.Description,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				logger.Log.Error().
					Str("revision_uuid", revisionUUID.String()).
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

	createdTaskSnapshot, err := r.getByCompositeKeyTx(ctx, tx, revisionUUID, taskID)
	if err != nil {
		return nil, err
	}
	if createdTaskSnapshot == nil {
		return nil, fmt.Errorf("failed to retrieve created task snapshot")
	}

	logger.Log.Debug().
		Str("revision_uuid", createdTaskSnapshot.RevisionUUID.String()).
		Int64("skeleton_id", createdTaskSnapshot.SkeletonID).
		Msg("Created new task snapshot successfully")
	return createdTaskSnapshot, nil
}

func (r *taskSnapshots) getByCompositeKeyTx(ctx context.Context, tx *sql.Tx, revision_uuid uuid.UUID, skeleton_id int64) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Str("revision_uuid", revision_uuid.String()).
		Int64("skeleton_id", skeleton_id).
		Msg("Trying to find the task snapshot in the db via ctx")

	query := `SELECT
		revision_uuid,
		skeleton_id,
		title,
		description,
		created_at,
		is_current
	FROM task_snapshots
	WHERE revision_uuid = ? AND skeleton_id = ?`

	var task models.TaskSnapshot

	err := tx.QueryRowContext(ctx, query, revision_uuid, skeleton_id).Scan(
		&task.RevisionUUID,
		&task.SkeletonID,
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

func (r *taskSnapshots) GetByCompositeKey(ctx context.Context, revision_uuid uuid.UUID, skeleton_id int64) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Str("revision_uuid", revision_uuid.String()).
		Int64("skeleton_id", skeleton_id).
		Msg("Trying to find the task snapshot in the db")

	query := `SELECT
		revision_uuid,
		skeleton_id,
		title,
		description,
		created_at,
		is_current
	FROM task_snapshots
	WHERE revision_uuid = ? AND skeleton_id = ?`

	var task models.TaskSnapshot

	err := r.db.QueryRowContext(ctx, query, revision_uuid, skeleton_id).Scan(
		&task.RevisionUUID,
		&task.SkeletonID,
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

func (r *taskSnapshots) GetLatest(ctx context.Context, skeleton_id int64) (*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Int64("skeleton_id", skeleton_id).
		Msg("Trying to find the latest task snapshot in the db")

	query := `SELECT
		revision_uuid,
		skeleton_id,
		title,
		description,
		created_at,
		is_current
	FROM task_snapshots
	WHERE skeleton_id = ? AND is_current = TRUE`

	var task models.TaskSnapshot

	err := r.db.QueryRowContext(ctx, query, skeleton_id).Scan(
		&task.RevisionUUID,
		&task.SkeletonID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.IsCurrent,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("latest task snapshot for skeleton_id %d not found", skeleton_id)
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskSnapshots) GetAllLatest(ctx context.Context, skeleton_ids []int64) ([]*models.TaskSnapshot, error) {
	logger.Log.Debug().
		Interface("skeleton_ids", skeleton_ids).
		Msg("Trying to find all latest task snapshots in the db")

	query := `SELECT
		revision_uuid,
		skeleton_id,
		title,
		description,
		created_at,
		is_current
	FROM task_snapshots
	WHERE skeleton_id IN (` + strings.Repeat("?,", len(skeleton_ids)-1) + "?) AND is_current = TRUE"

	args := make([]any, len(skeleton_ids))
	for i, id := range skeleton_ids {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.TaskSnapshot
	for rows.Next() {
		var task models.TaskSnapshot
		err := rows.Scan(
			&task.RevisionUUID,
			&task.SkeletonID,
			&task.Title,
			&task.Description,
			&task.CreatedAt,
			&task.IsCurrent,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no latest task snapshots found for provided skeleton_ids")
	}

	return tasks, nil
}

func (r *taskSnapshots) SetCurrentRevisionTx(ctx context.Context, tx *sql.Tx, taskID int64, revisionUUID uuid.UUID) error {
	logger.Log.Debug().
		Int64("skeleton_id", taskID).
		Str("revision_uuid", revisionUUID.String()).
		Msg("Trying to set current revision")

	resetQuery := `
		UPDATE task_snapshots 
		SET is_current = FALSE 
		WHERE skeleton_id = ? AND is_current = TRUE`

	_, err := tx.ExecContext(ctx, resetQuery, taskID)
	if err != nil {
		return err
	}

	setQuery := `
		UPDATE task_snapshots 
		SET is_current = TRUE 
		WHERE skeleton_id = ? AND revision_uuid = ?`

	_, err = tx.ExecContext(ctx, setQuery, taskID, revisionUUID)
	return err
}
