package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type Users interface {
	WithTx(tx *sqlx.Tx) Users

	Create(ctx context.Context, user *models.User) error

	Get(ctx context.Context, user *models.UserContext) *UsersQuery
	GetContextByID(ctx context.Context, id int64) (*models.UserContext, error)
	Exists(ctx context.Context, userID int64) (bool, error)
}

type users struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

var _ Users = (*users)(nil)

func InitUsers(db *sqlx.DB) (Users, error) {
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

func (r *users) WithTx(tx *sqlx.Tx) Users {
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

// ---------------- //
// INSERT FUNCTIONS //
// ---------------- //

func (r *users) Create(ctx context.Context, user *models.User) error {
	logger.Log.Info().
		Str("email", user.Email).
		Msg("Trying to create new user in the db")

	executor := r.getExecutor()

	query := `
        INSERT INTO users (
            name,
            email,
            password_hash
        )
        VALUES (?, ?, ?)
        RETURNING id, name, email, password_hash, created_at`

	err := executor.GetContext(
		ctx,
		user,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return apperrors.ErrDuplicate
			}
		}
		return err
	}

	logger.Log.Info().
		Str("name", user.Username).
		Int64("id", user.ID).
		Time("created_at", user.CreatedAt).
		Msg("User was created")

	return nil
}

// ------------- //
// GET FUNCTIONS //
// ------------- //

type UsersParams struct {
	IDs    []int64
	Emails []string
}

type UsersQuery struct {
	repo   *users
	uc     *models.UserContext
	params *UsersParams
	ctx    context.Context
}

func (r *users) Get(ctx context.Context, uc *models.UserContext) *UsersQuery {
	return &UsersQuery{
		repo:   r,
		uc:     uc,
		params: &UsersParams{},
		ctx:    ctx,
	}
}

func (q *UsersQuery) WithIDs(ids ...any) *UsersQuery {
	for _, id := range ids {
		switch v := id.(type) {
		case int64:
			q.params.IDs = append(q.params.IDs, v)
		case []int64:
			q.params.IDs = append(q.params.IDs, v...)
		}
	}
	return q
}

func (q *UsersQuery) WithEmails(emails ...any) *UsersQuery {
	for _, e := range emails {
		switch v := e.(type) {
		case string:
			q.params.Emails = append(q.params.Emails, v)
		case []string:
			q.params.Emails = append(q.params.Emails, v...)
		}
	}
	return q
}

func (r *users) BuildQuery(params *UsersParams, uc *models.UserContext) (string, []any, error) {
	var whereClauses []string
	var args []any
	var err error

	switch uc.Role {
	case models.RoleAdmin:
		// no filter
	case models.RoleUser:
		// no filter
	}

	whereClauses, args, err = InQuery(whereClauses, args, "id", toAnySlice(params.IDs))
	if err != nil {
		return "", nil, err
	}

	whereClauses, args, err = InQuery(whereClauses, args, "email", toAnySlice(params.Emails))
	if err != nil {
		return "", nil, err
	}

	query := "SELECT id, name, email, password_hash, created_at FROM users"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	return query, args, nil
}

func (q *UsersQuery) All() ([]models.User, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug().
		Interface("query", query).
		Interface("args", args).
		Msg("Executing users query")

	var users []models.User
	if err := q.repo.db.SelectContext(q.ctx, &users, query, args...); err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	return users, nil
}

func (q *UsersQuery) First() (*models.User, error) {
	query, args, err := q.repo.BuildQuery(q.params, q.uc)
	if err != nil {
		return nil, err
	}

	// force limit 1
	query += " LIMIT 1"

	var user models.User
	if err := q.repo.db.GetContext(q.ctx, &user, query, args...); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *users) GetContextByID(ctx context.Context, id int64) (*models.UserContext, error) {
	var uc models.UserContext
	query := "SELECT id, role FROM users WHERE id = ? LIMIT 1"
	if err := r.db.GetContext(ctx, &uc, query, id); err != nil {
		return nil, fmt.Errorf("failed to fetch user context: %w", err)
	}
	return &uc, nil
}

func (r *users) Exists(ctx context.Context, userID int64) (bool, error) {
	executor := r.getExecutor()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
	var exists bool
	err := executor.GetContext(ctx, &exists, query, userID)
	return exists, err
}
