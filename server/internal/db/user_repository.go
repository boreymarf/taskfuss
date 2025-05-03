package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string // TODO: Make it hashed later
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) (*UserRepository, error) {
	repo := &UserRepository{db: db}

	// Выполняем миграции при создании репозитория
	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *UserRepository) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
	userId 				INT NOT NULL PRIMARY KEY,
	name 					VARCHAR(255) NOT NULL,
	email 				VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL
	created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP
	updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateUser(user *User) error {

	query := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id

	logger.Log.Info().
		Str("name", user.Name).
		Int64("id", user.ID).
		Time("created_at", user.CreatedAt).
		Msg("User was created")
	return nil
}

func (r *UserRepository) GetUserByID(id int64, user *User) error {

	query := `SELECT id, name, email, created_at, updated_at 
	FROM users 
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error().
			Int64("userID", id).
			Msg("No User was found")
		return fmt.Errorf("user %d not found: %w", id, err)
	} else if err != nil {
		return err
	}

	return nil
}
