package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/mattn/go-sqlite3"
)

type Users interface {
	WithTx(tx *sql.Tx) Users

	CreateTable() error

	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64, user *models.User) error
	GetUserByEmail(ctx context.Context, email string, user *models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	Exists(ctx context.Context, userID int64) (bool, error)
}

type users struct {
	db *sql.DB
	tx *sql.Tx
}

var _ Users = (*users)(nil)

func InitUsers(db *sql.DB) (Users, error) {
	repo := &users{db: db}

	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Log.Debug().Msg("Users repository initialization completed")
	return repo, nil
}

func (r *users) CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
	id 				INTEGER NOT NULL PRIMARY KEY,
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

func (r *users) WithTx(tx *sql.Tx) Users {
	return &users{
		db: r.db,
		tx: tx,
	}
}

func (r *users) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *users) CreateUser(ctx context.Context, user *models.User) error {
	logger.Log.Info().Str("email", user.Email).Msg("Users tries to create new user")

	// Use getExecutor() to support transactions
	executor := r.getExecutor()

	query := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`

	result, err := executor.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		// Handle duplicate user error
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

	// Use the same transaction for GetUserByID
	err = r.GetUserByID(ctx, id, user)
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

func (r *users) GetUserByID(ctx context.Context, id int64, user *models.User) error {
	// Use getExecutor() to support transactions
	executor := r.getExecutor()

	query := `SELECT id, name, email, password_hash, created_at, updated_at 
              FROM users WHERE id = ?`

	row := executor.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found with ID %d", id)
		}
		return err
	}

	return nil
}

func (r *users) GetUserByEmail(ctx context.Context, email string, user *models.User) error {
	logger.Log.Debug().Str("email", email).Msg("Users tries to find user")

	executor := r.getExecutor()

	query := `SELECT id, name, password_hash, email, created_at, updated_at 
    FROM users 
    WHERE email = ? COLLATE NOCASE`

	row := executor.QueryRowContext(ctx, query, email)

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

func (r *users) GetAllUsers(ctx context.Context) ([]models.User, error) {
	executor := r.getExecutor()

	query := `SELECT id, name, email, created_at, updated_at 
              FROM users`

	rows, err := executor.QueryContext(ctx, query)
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

func (r *users) Exists(ctx context.Context, userID int64) (bool, error) {
	executor := r.getExecutor()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	var exists bool
	err := executor.QueryRowContext(ctx, query, userID).Scan(&exists)
	return exists, err
}
