package db

import (
	"database/sql"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type TaskRepository struct {
	db *sql.DB
}

func InitTaskRepository(db *sql.DB) (*TaskRepository, error) {

	repo := &TaskRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("UserRepository initialization completed")

	return repo, nil
}

func (r *TaskRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
	userId 				INTEGER NOT NULL PRIMARY KEY,
	name 					VARCHAR(255) NOT NULL,
	email 				VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
