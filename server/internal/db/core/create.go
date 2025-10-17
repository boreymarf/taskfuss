package core

import (
	"context"
	"fmt"
	"strings"
)

func (r *BaseRepo[T]) Create(ctx context.Context, obj *T) error {
	cols := r.ColumnNames()
	ph := r.ColumnPlaceholders()

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.table,
		strings.Join(cols, ", "),
		strings.Join(ph, ", "),
	)

	_, err := r.GetExec().NamedExecContext(ctx, query, obj)
	return err
}
