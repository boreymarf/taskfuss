package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type RequirementSnapshots interface {
	CreateTx(ctx context.Context, tx *sql.Tx, requirementSnapshot *models.RequirementSnapshot, revision_uuid uuid.UUID) (*models.RequirementSnapshot, error)
}

type requirementSnapshots struct {
	db *sql.DB
}

var _ RequirementSnapshots = (*requirementSnapshots)(nil)

func InitRequirementSnapshots(db *sql.DB) (*requirementSnapshots, error) {

	repo := &requirementSnapshots{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *requirementSnapshots) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS requirement_snapshots (
		revision_uuid TEXT NOT NULL REFERENCES task_snapshots(revision_uuid) ON DELETE CASCADE,
		skeleton_id INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		parent_id INTEGER DEFAULT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('atom', 'condition')),
		data_type TEXT CHECK (data_type IN ('bool', 'int', 'float', 'duration', 'none')),
		operator TEXT CHECK (operator IN ('or', 'not', 'and', '==', '>=', '<=', '!=', '>', '<')),
		target_value TEXT,
		sort_order INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (revision_uuid, skeleton_id)
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *requirementSnapshots) CreateTx(ctx context.Context, tx *sql.Tx, requirementSnapshot *models.RequirementSnapshot, revision_uuid uuid.UUID) (*models.RequirementSnapshot, error) {
	logger.Log.Debug().
		Int64("skeleton_id", requirementSnapshot.SkeletonID).
		Msg("Trying to create new requirement snapshot in db via ctx")

	query := `INSERT INTO requirement_snapshots (
		revision_uuid,
		skeleton_id,
		parent_id,
		title,
		type,
		data_type,
		operator,
		target_value,
		sort_order
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	spew.Dump(requirementSnapshot)

	result, err := tx.ExecContext(
		ctx,
		query,
		revision_uuid,
		requirementSnapshot.SkeletonID,
		requirementSnapshot.ParentID,
		requirementSnapshot.Title,
		requirementSnapshot.Type,
		requirementSnapshot.DataType,
		requirementSnapshot.Operator,
		requirementSnapshot.TargetValue,
		requirementSnapshot.SortOrder,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				logger.Log.Error().
					Str("revision_uuid", revision_uuid.String()).
					Msg("For some reason there's duplicate of the requirement snapshot!")
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

	createdRequirementSnapshot, err := r.getByRevisionUUIDTx(ctx, tx, revision_uuid)
	if err != nil {
		return nil, err
	}
	if createdRequirementSnapshot == nil {
		return nil, fmt.Errorf("failed to retrieve created requirement snapshot")
	}

	logger.Log.Debug().
		Str("revision_uuid", createdRequirementSnapshot.RevisionUUID.String()).
		Int64("skeleton_id", createdRequirementSnapshot.SkeletonID).
		Msg("Created new task snapshot successfully")
	return createdRequirementSnapshot, nil
}

func (r *requirementSnapshots) getByRevisionUUIDTx(ctx context.Context, tx *sql.Tx, revision_uuid uuid.UUID) (*models.RequirementSnapshot, error) {
	logger.Log.Debug().
		Str("revision_uuid", revision_uuid.String()).
		Msg("Trying to find the requirement snapshot in the db via ctx")

	query := `SELECT
		revision_uuid,
		skeleton_id,
		parent_id,
		title,
		type,
		data_type,
		operator,
		target_value,
		sort_order
	FROM requirement_snapshots
	WHERE revision_uuid = ?`

	var requirementSnapshot models.RequirementSnapshot

	err := tx.QueryRowContext(ctx, query, revision_uuid).Scan(
		&requirementSnapshot.RevisionUUID,
		&requirementSnapshot.SkeletonID,
		&requirementSnapshot.ParentID,
		&requirementSnapshot.Title,
		&requirementSnapshot.Type,
		&requirementSnapshot.DataType,
		&requirementSnapshot.Operator,
		// FIXME: Returns "" for some reason
		&requirementSnapshot.TargetValue,
		&requirementSnapshot.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	spew.Dump(requirementSnapshot)

	return &requirementSnapshot, nil
}
