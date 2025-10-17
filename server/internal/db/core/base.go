package core

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SQLExecutor interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

var _ SQLExecutor = (*sqlx.DB)(nil)
var _ SQLExecutor = (*sqlx.Tx)(nil)

// --------------------
// BaseRepo
// --------------------

type BaseRepo[T any] struct {
	db    *sqlx.DB
	tx    *sqlx.Tx
	table string
}

func InitRepo[T any](db *sqlx.DB, table string, createSQL string) (*BaseRepo[T], error) {
	repo := &BaseRepo[T]{db: db, table: table}
	ctx := context.Background()
	if err := repo.CreateTable(ctx, createSQL); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *BaseRepo[T]) CreateTable(ctx context.Context, sqlDef string) error {
	sqlDef = strings.TrimSpace(sqlDef)
	if sqlDef == "" {
		return fmt.Errorf("CreateTable: empty SQL definition")
	}
	if !strings.HasPrefix(strings.ToUpper(sqlDef), "CREATE TABLE") {
		return fmt.Errorf("CreateTable: SQL does not start with CREATE TABLE")
	}
	_, err := r.GetExec().ExecContext(ctx, sqlDef)
	if err != nil {
		return fmt.Errorf("CreateTable execution failed: %w", err)
	}
	return nil
}

func (r *BaseRepo[T]) GetExec() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *BaseRepo[T]) WithTx(tx *sqlx.Tx) Repository[T] {
	return &BaseRepo[T]{db: r.db, tx: tx, table: r.table}
}

// --------------------
// BuildQuery helper
// --------------------

func (r *BaseRepo[T]) BuildQuery(params map[string]any) (string, []any) {
	where := []string{}
	args := []any{}

	for key, val := range params {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Slice {
			if v.Len() == 0 {
				continue
			}
			placeholders := strings.TrimRight(strings.Repeat("?,", v.Len()), ",")
			where = append(where, fmt.Sprintf("%s IN (%s)", key, placeholders))
			for i := 0; i < v.Len(); i++ {
				args = append(args, v.Index(i).Interface())
			}
		} else {
			where = append(where, fmt.Sprintf("%s = ?", key))
			args = append(args, val)
		}
	}

	return strings.Join(where, " AND "), args
}

// --------------------
// Helper methods
// --------------------

func (r *BaseRepo[T]) ColumnNames() []string {
	var cols []string
	t := reflect.TypeOf(*new(T))
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag != "" && tag != "-" {
			cols = append(cols, tag)
		}
	}
	return cols
}

func (r *BaseRepo[T]) ColumnPlaceholders() []string {
	var ph []string
	for _, col := range r.ColumnNames() {
		ph = append(ph, ":"+col)
	}
	return ph
}
