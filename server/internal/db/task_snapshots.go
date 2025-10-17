package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type TaskSnapshots interface {
	core.Repository[models.TaskSnapshot]
	core.Creator[models.TaskSnapshot]
	core.Getter[models.TaskSnapshot]
	core.Updater[models.TaskSnapshot]
	core.Deleter[models.TaskSnapshot]
}

type taskSnapshotsRepo struct {
	*core.BaseRepo[models.TaskSnapshot]
}

const TaskSnapshotsSQL = `
	CREATE TABLE IF NOT EXISTS task_snapshots (
		revision_uuid TEXT NOT NULL,
		skeleton_id INTEGER  NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		description TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_current BOOLEAN DEFAULT FALSE,
	PRIMARY KEY (revision_uuid, skeleton_id)
	)`

func InitTaskSnapshots(db *sqlx.DB) (TaskSnapshots, error) {
	repo, err := core.InitRepo[models.TaskSnapshot](db, "task_snapshots", TaskSnapshotsSQL)
	if err != nil {
		return nil, err
	}
	return &taskSnapshotsRepo{BaseRepo: repo}, nil
}
