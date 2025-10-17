package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type UpdateQuery[T any] struct {
	repo      *BaseRepo[T]
	ctx       context.Context
	uc        *models.UserContext
	params    map[string]any
	setParams map[string]any
}

func (r *BaseRepo[T]) Update(ctx context.Context, uc *models.UserContext) *UpdateQuery[T] {
	return &UpdateQuery[T]{repo: r, ctx: ctx, uc: uc, params: make(map[string]any), setParams: make(map[string]any)}
}

func (q *UpdateQuery[T]) With(column string, value any) *UpdateQuery[T] {
	q.params[column] = value
	return q
}

func (q *UpdateQuery[T]) WithTx(tx *sqlx.Tx) *UpdateQuery[T] {
	q.repo.tx = tx
	return q
}

func (q *UpdateQuery[T]) Set(column string, value any) *UpdateQuery[T] {
	q.setParams[column] = value
	return q
}

func (q *UpdateQuery[T]) One() error {
	if len(q.setParams) == 0 {
		return fmt.Errorf("no SET fields provided")
	}

	args := []any{}
	setClause := []string{}
	for col, val := range q.setParams {
		setClause = append(setClause, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	whereClause, whereArgs := q.repo.BuildQuery(q.params)
	args = append(args, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s", q.repo.table, strings.Join(setClause, ", "))
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	_, err := q.repo.GetExec().ExecContext(q.ctx, query, args...)
	return err
}
