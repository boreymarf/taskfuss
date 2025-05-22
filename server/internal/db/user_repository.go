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

type UserRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) (*UserRepository, error) {
	repo := &UserRepository{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repo, nil
}

func (r *UserRepository) CreateTable() error {
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

func (r *UserRepository) CreateUser(user *models.User) error {

	query := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		// If there's a dublicate user
		var sqliteErr sqlite3.Error
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

	user.ID = id

	err = r.GetUserByID(id, user)
	if err != nil {
		return err
	}

	logger.Log.Info().
		Str("name", user.Username).
		Int64("id", user.ID).
		Time("created_at", user.CreatedAt).
		Msg("User was created")
	return nil
}

func (r *UserRepository) GetUserByID(id int64, user *models.User) error {

	query := `SELECT userId, name, email, created_at, updated_at 
	FROM users 
	WHERE userId = ?`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
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

func (r *UserRepository) GetUserByEmail(email string, user *models.User) error {

	query := `SELECT userId, name, password_hash, email, created_at, updated_at 
	FROM users 
	WHERE email = ?`

	row := r.db.QueryRow(query, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Error().
			Str("email", email).
			Msg("No User was found")
		return apperrors.ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {

	query := `SELECT userId, name, email, created_at, updated_at 
              FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to execute users query")
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Failed to scan user row")
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Log.Error().Err(err).Msg("Rows iteration error")
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}
