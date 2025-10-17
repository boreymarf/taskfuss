package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type requirementSnapshots interface {
	core.Repository[models.RequirementSnapshot]
	core.Creator[models.RequirementSnapshot]
	core.Getter[models.RequirementSnapshot]
	core.Updater[models.RequirementSnapshot]
	core.Deleter[models.RequirementSnapshot]
}

type requirementSnapshotsRepo struct {
	*core.BaseRepo[models.RequirementSnapshot]
}

const requirementSnapshotsSQL = `
	CREATE TABLE IF NOT EXISTS requirement_snapshots (
		revision_uuid TEXT NOT NULL REFERENCES task_snapshots(revision_uuid) ON DELETE CASCADE,
		skeleton_id INTEGER NOT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		parent_id INTEGER DEFAULT NULL REFERENCES requirement_skeletons(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('atom', 'condition')),
		data_type TEXT NOT NULL CHECK (data_type IN ('bool', 'int', 'float', 'duration', 'none')),
		operator TEXT CHECK (operator IN ('or', 'not', 'and', '==', '>=', '<=', '!=', '>', '<')),
		target_value TEXT NOT NULL,
		sort_order INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (revision_uuid, skeleton_id)
	)`

func InitRequirementSnapshots(db *sqlx.DB) (requirementSnapshots, error) {
	repo, err := core.InitRepo[models.RequirementSnapshot](db, "requirement_snapshots", requirementSnapshotsSQL)
	if err != nil {
		return nil, err
	}
	return &requirementSnapshotsRepo{BaseRepo: repo}, nil
}
