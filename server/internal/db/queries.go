package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type queryParams struct {
	IDs            []int64
	entryIDs       []int64
	requirementIDs []int64
	dates          []time.Time
	startDate      *time.Time
	endDate        *time.Time
}

type QueryOption func(*queryParams)

// Helper function
func toAnySlice[T any](s []T) []any {
	a := make([]any, len(s))
	for i, v := range s {
		a[i] = v
	}
	return a
}

// InQuery adds an IN clause to the existing whereClauses and args
func InQuery(whereClauses []string, args []any, column string, values []any) ([]string, []any, error) {
	if len(values) == 0 {
		return whereClauses, args, nil
	}

	queryPart, queryArgs, err := sqlx.In(fmt.Sprintf("%s IN (?)", column), values)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build IN clause for %s: %w", column, err)
	}

	whereClauses = append(whereClauses, queryPart)
	args = append(args, queryArgs...)
	return whereClauses, args, nil
}

// BetweenQuery adds a BETWEEN clause to the existing whereClauses and args
func BetweenQuery(whereClauses []string, args []any, column string, start, end any) ([]string, []any) {
	if start != nil && end != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN ? AND ?", column))
		args = append(args, start, end)
	}
	return whereClauses, args
}
