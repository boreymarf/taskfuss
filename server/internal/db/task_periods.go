package db

import (
	"database/sql"
	"fmt"
)

type TaskPeriods struct {
	db *sql.DB
}

func InitTaskPeriods(db *sql.DB) (*TaskPeriods, error) {

	repo := &TaskPeriods{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *TaskPeriods) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS task_periods (
		id							INTEGER NOT NULL PRIMARY KEY,
		task_id					INTEGER NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		start_date			DATETIME NOT NULL,
		end_date				DATETIME,
		is_active				BOOLEAN DEFAULT TRUE,
		created_at			DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at			DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
