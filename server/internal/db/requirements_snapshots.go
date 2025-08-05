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

type RequirementsSnapshotsRepository struct {
	db *sql.DB
}

func InitRequirementsSnapshotsRepository(db *sql.DB) (*RequirementsSnapshotsRepository, error) {

	repo := &RequirementsSnapshotsRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("Repository initialization completed")

	return repo, nil
}

func (r *RequirementsSnapshotsRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS requirements_snapshots (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  requirement_id     INTEGER NOT NULL REFERENCES requirements(id) ON DELETE CASCADE,
  task_id            INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  type               TEXT NOT NULL CHECK (type IN ('atom', 'condition')),
  title              TEXT NOT NULL,
  data_type          TEXT CHECK (data_type IN ('bool', 'int', 'float', 'duration', 'none')),
  operator           TEXT CHECK (operator IN ('or', 'not', 'and', '==', '>=', '<=', '!=', '>', '<')),
  target_value       TEXT,
  parent_id          INTEGER REFERENCES requirements(id) ON DELETE SET NULL,
  parent_snapshot_id INTEGER REFERENCES requirements_snapshots(id) ON DELETE SET NULL DEFAULT NULL,
  sort_order         INTEGER NOT NULL DEFAULT 0,
  valid_from         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  valid_to           DATETIME DEFAULT NULL
)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *RequirementsSnapshotsRepository) CreateSnapshot(snapshot *models.RequirementSnapshot) error {
	logger.Log.Debug().
		Str("title", snapshot.Title).
		Msg("Trying to create new requirement snapshot to the db...")

	query := `INSERT INTO requirements (
	requirement_id,
	task_id,
	type,
	title,
	data_type,
	operator,
	target_value,
	parent_id,
	parent_snapshot_id,
	sort_order
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		snapshot.RequirementID,
		snapshot.TaskID,
		snapshot.Type,
		snapshot.Title,
		snapshot.DataType,
		snapshot.Operator,
		snapshot.TargetValue,
		snapshot.ParentID,
		snapshot.ParentSnapshotID,
		snapshot.SortOrder,
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

	return nil
}

func (r *RequirementsSnapshotsRepository) GetSnapshotByID(id int64) (models.RequirementSnapshot, error) {

	var snapshot models.RequirementSnapshot
	logger.Log.Debug().Int64("id", id).Msg("taskRepository tries to find task")

	query := `SELECT
	id,
	requirement_id,
	task_id,
	type,
	title,
	data_type,
	operator,
	target_value,
	parent_id,
	parent_snapshot_id,
	sort_order
	FROM requirements_snapshots
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&snapshot.ID,
		&snapshot.RequirementID,
		&snapshot.TaskID,
		&snapshot.Type,
		&snapshot.Title,
		&snapshot.DataType,
		&snapshot.Operator,
		&snapshot.TargetValue,
		&snapshot.ParentID,
		&snapshot.ParentSnapshotID,
		&snapshot.SortOrder,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Warn().
			Int64("taskID", id).
			Msg("No task was found")
		return models.RequirementSnapshot{}, fmt.Errorf("task %d not found: %w", id, err)
	} else if err != nil {
		return models.RequirementSnapshot{}, err
	}

	return snapshot, nil
}

// func (r *RequirementRepository) GetRequirementsByTaskIDs(taskIDs []int64) ([]models.Requirement, error) {
//
// 	var stringIDs []string
//
// 	for _, id := range taskIDs {
// 		stringIDs = append(stringIDs, strconv.FormatInt(id, 10))
// 	}
// 	idQuery := strings.Join(stringIDs, ", ")
//
// 	query := fmt.Sprintf(`SELECT id, task_id, parent_id, title, type, data_type, operator, target_value, sort_order
// 		FROM requirements WHERE task_id IN (%s)`, idQuery)
//
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query requirements: %w", err)
// 	}
// 	defer rows.Close()
//
// 	var requirements []models.Requirement
// 	for rows.Next() {
// 		var req models.Requirement
// 		err := rows.Scan(
// 			&req.ID,
// 			&req.TaskID,
// 			&req.ParentID,
// 			&req.Title,
// 			&req.Type,
// 			&req.DataType,
// 			&req.Operator,
// 			&req.TargetValue,
// 			&req.SortOrder,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan requirement: %w", err)
// 		}
// 		requirements = append(requirements, req)
// 	}
//
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error after scanning requirements: %w", err)
// 	}
//
// 	return requirements, nil
// }
