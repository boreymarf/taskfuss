package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SQLExecutor interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

var _ SQLExecutor = (*sqlx.DB)(nil)
var _ SQLExecutor = (*sqlx.Tx)(nil)
