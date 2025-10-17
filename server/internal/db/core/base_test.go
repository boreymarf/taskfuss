package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test struct for reflection
type TestModel struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Ignore string
	Skip   string `db:"-"`
}

func TestBaseRepo_ColumnNames(t *testing.T) {
	repo := &BaseRepo[TestModel]{table: "test_table"}

	cols := repo.ColumnNames()
	assert.Equal(t, []string{"id", "name"}, cols)
}

func TestBaseRepo_ColumnPlaceholders(t *testing.T) {
	repo := &BaseRepo[TestModel]{table: "test_table"}

	ph := repo.ColumnPlaceholders()
	assert.Equal(t, []string{":id", ":name"}, ph)
}
