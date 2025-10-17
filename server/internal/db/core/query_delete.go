package core

import (
	"context"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type DeleteQuery[T any] struct {
	repo   *BaseRepo[T]
	ctx    context.Context
	uc     *models.UserContext
	params map[string]any
}

func (r *BaseRepo[T]) Delete(ctx context.Context, uc *models.UserContext) *DeleteQuery[T] {
	return &DeleteQuery[T]{repo: r, ctx: ctx, uc: uc, params: make(map[string]any)}
}

func (q *DeleteQuery[T]) With(column string, value any) *DeleteQuery[T] {
	q.params[column] = value
	return q
}

func (q *DeleteQuery[T]) WithTx(tx *sqlx.Tx) *DeleteQuery[T] {
	q.repo.tx = tx
	return q
}

func (q *DeleteQuery[T]) One() error {
	where, args := q.repo.BuildQuery(q.params)
	query := fmt.Sprintf("DELETE FROM %s", q.repo.table)
	if where != "" {
		query += " WHERE " + where
	}
	_, err := q.repo.GetExec().ExecContext(q.ctx, query, args...)
	return err
}
