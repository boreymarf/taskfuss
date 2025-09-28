package db

import (
	"context"
	"database/sql"
)

type SQLExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

var _ SQLExecutor = (*sql.DB)(nil)
var _ SQLExecutor = (*sql.Tx)(nil)
