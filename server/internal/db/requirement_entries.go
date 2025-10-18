package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type RequirementEntries interface {
	core.Repository[models.RequirementEntry]
	core.Creator[models.RequirementEntry]
	core.Getter[models.RequirementEntry]
	core.Updater[models.RequirementEntry]
	core.Deleter[models.RequirementEntry]
}

type requirementEntriesRepo struct {
	*core.BaseRepo[models.RequirementEntry]
}

const requirementEntriesSQL = `
	CREATE TABLE IF NOT EXISTS requirement_entries (
		id 							INTEGER NOT NULL PRIMARY KEY,
		requirement_id 	INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		revision_uuid		TEXT NOT NULL REFERENCES requirement_snapshots(revision_uuid) ON DELETE CASCADE,
		entry_date			DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		value TEXT,
		UNIQUE(requirement_id, entry_date)
	)`

func InitRequirementEntries(db *sqlx.DB) (RequirementEntries, error) {
	repo, err := core.InitRepo[models.RequirementEntry](db, "requirement_entries", requirementEntriesSQL)
	if err != nil {
		return nil, err
	}
	return &requirementEntriesRepo{BaseRepo: repo}, nil
}
