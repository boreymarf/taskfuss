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

	var created T
	err := r.GetExec().QueryRowxContext(ctx, query, obj).StructScan(&created)
	if err != nil {
		return nil, err
	}

	return &created, nil
}
