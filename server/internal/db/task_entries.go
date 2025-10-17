package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type taskEntries interface {
	core.Repository[models.TaskEntry]
	core.Creator[models.TaskEntry]
	core.Getter[models.TaskEntry]
	core.Updater[models.TaskEntry]
	core.Deleter[models.TaskEntry]
}

type taskEntriesRepo struct {
	*core.BaseRepo[models.TaskEntry]
}

const taskEntriesSQL = `
	CREATE TABLE IF NOT EXISTS task_entries (
		id 						INTEGER NOT NULL PRIMARY KEY,
		task_id 			INTEGER NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		entry_date 		DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed			BOOLEAN NOT NULL DEFAULT FALSE CHECK (completed IN (0, 1))
	)`

func InittaskEntries(db *sqlx.DB) (taskEntries, error) {
	repo, err := core.InitRepo[models.TaskEntry](db, "task_entries", taskEntriesSQL)
	if err != nil {
		return nil, err
	}
	return &taskEntriesRepo{BaseRepo: repo}, nil
}
