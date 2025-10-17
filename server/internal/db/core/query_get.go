package core

import (
	"context"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type GetQuery[T any] struct {
	repo   *BaseRepo[T]
	ctx    context.Context
	uc     *models.UserContext
	params map[string]any
}

func (r *BaseRepo[T]) Get(ctx context.Context, uc *models.UserContext) *GetQuery[T] {
	return &GetQuery[T]{repo: r, ctx: ctx, uc: uc, params: make(map[string]any)}
}

func (q *GetQuery[T]) With(column string, value any) *GetQuery[T] {
	q.params[column] = value
	return q
}

func (q *GetQuery[T]) WithTx(tx *sqlx.Tx) *GetQuery[T] {
	q.repo.tx = tx
	return q
}

func (q *GetQuery[T]) All() ([]T, error) {
	query, args := q.repo.BuildQuery(q.params)
	var res []T
	err := q.repo.GetExec().SelectContext(q.ctx, &res, fmt.Sprintf("SELECT * FROM %s WHERE %s", q.repo.table, query), args...)
	return res, err
}

func (q *GetQuery[T]) First() (*T, error) {
	query, args := q.repo.BuildQuery(q.params)
	qstr := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", q.repo.table, query)
	var res T
	err := q.repo.GetExec().GetContext(q.ctx, &res, qstr, args...)
	return &res, err
}
