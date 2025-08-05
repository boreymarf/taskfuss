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

type RequirementsRepository struct {
	db *sql.DB
}

func InitRequirementsRepository(db *sql.DB) (*RequirementsRepository, error) {

	repo := &RequirementsRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("taskRepository initialization completed")

	return repo, nil
}

func (r *RequirementsRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS requirements (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id   INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE
)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *RequirementsRepository) CreateRequirement(requirement *models.Requirement) (models.Requirement, error) {
	logger.Log.Debug().
		Msg("Trying to create new requirement to the db...")

	query := `INSERT INTO requirements (task_id) VALUES (?)`

	result, err := r.db.Exec(
		query,
		requirement.TaskID,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		// If there's a dublicate requirement
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return models.Requirement{}, apperrors.ErrDuplicate
			}
		}
		return models.Requirement{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Requirement{}, err
	}

	createdRequirement, err := r.GetRequirementByID(id)
	if err != nil {
		return models.Requirement{}, err
	}

	return createdRequirement, nil
}

func (r *RequirementsRepository) GetRequirementByID(id int64) (models.Requirement, error) {

	var requirement models.Requirement
	logger.Log.Debug().Int64("id", id).Msg("taskRepository tries to find task")

	query := `SELECT id, task_id
	FROM requirements
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&requirement.ID,
		&requirement.TaskID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Warn().
			Int64("taskID", id).
			Msg("No task was found")
		return models.Requirement{}, fmt.Errorf("task %d not found: %w", id, err)
	} else if err != nil {
		return models.Requirement{}, err
	}

	return requirement, nil
}

// func (r *RequirementsRepository) GetRequirementsByTaskIDs(taskIDs []int64) ([]models.Requirement, error) {
//
// 	var stringIDs []string
//
// 	for _, id := range taskIDs {
// 		stringIDs = append(stringIDs, strconv.FormatInt(id, 10))
// 	}
// 	idQuery := strings.Join(stringIDs, ", ")
//
// 	query := fmt.Sprintf(`SELECT id, task_id`, idQuery)
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
