package service

import (
	"database/sql"

	"github.com/boreymarf/task-fuss/server/internal/db"
)

type EntriesService struct {
	db                 *sql.DB
	requirementEntries db.RequirementEntries
}

func InitEntriesService(
	requirementEntries db.RequirementEntries,
) (*)
