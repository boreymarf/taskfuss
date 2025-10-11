package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SQLExecutor interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
}

var _ SQLExecutor = (*sqlx.DB)(nil)
var _ SQLExecutor = (*sqlx.Tx)(nil)
