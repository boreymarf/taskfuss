package core

import (
	"context"
	"fmt"
	"strings"
)

func (r *BaseRepo[T]) Create(ctx context.Context, obj *T) (*T, error) {
	cols := r.ColumnNames()
	ph := r.ColumnPlaceholders()

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		r.table,
		strings.Join(cols, ", "),
		strings.Join(ph, ", "),
	)

	_, err := r.GetExec().NamedExecContext(ctx, query, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
