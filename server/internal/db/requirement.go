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

type RequirementRepository struct {
	db *sql.DB
}

func InitRequirementRepository(db *sql.DB) (*RequirementRepository, error) {

	repo := &RequirementRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("taskRepository initialization completed")

	return repo, nil
}

func (r *RequirementRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS requirements (
  id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  task_id      INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  parent_id    INTEGER REFERENCES requirements(id) ON DELETE CASCADE,
	title        TEXT NOT NULL,
  type         TEXT NOT NULL CHECK (type IN ('atom', 'condition')),
  data_type    TEXT CHECK (data_type IN ('bool', 'int', 'float', 'duration', 'none')),
  operator     TEXT CHECK (operator IN ('or', 'not', 'and', '==', '>=', '<=', '!=', '>', '<')),
  target_value TEXT,
  sort_order   INTEGER NOT NULL DEFAULT 0
  )`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *RequirementRepository) AddRequirement(requirement *models.Requirement) error {
	logger.Log.Debug().
		Str("title", requirement.Title).
		Msg("Trying to add new requirement to the db...")

	query := `INSERT INTO requirements (
		task_id,
		parent_id,
		title,
		type,
		data_type,
		operator,
		target_value,
		sort_order
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		requirement.TaskID,
		requirement.ParentID,
		requirement.Title,
		requirement.Type,
		requirement.DataType,
		requirement.Operator,
		requirement.TargetValue,
		requirement.SortOrder,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		// If there's a dublicate requirement
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

	requirement.ID = id

	// err = r.GetTaskByID(id, task)
	// if err != nil {
	// 	return err
	// }

	return nil
}
