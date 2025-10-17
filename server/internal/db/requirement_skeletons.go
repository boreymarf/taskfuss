package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type requirementSkeletons interface {
	core.Repository[models.RequirementSkeleton]
	core.Creator[models.RequirementSkeleton]
	core.Getter[models.RequirementSkeleton]
	core.Updater[models.RequirementSkeleton]
	core.Deleter[models.RequirementSkeleton]
}

type requirementSkeletonsRepo struct {
	*core.BaseRepo[models.RequirementSkeleton]
}

const requirementSkeletonsSQL = `
	CREATE TABLE IF NOT EXISTS requirement_skeletons (
		id INTEGER PRIMARY KEY AUTOINCREMENT REFERENCES task_skeletons(id) ON DELETE CASCADE,
		task_id INTEGER NOT NULL
  )`

func InitrequirementSkeletons(db *sqlx.DB) (requirementSkeletons, error) {
	repo, err := core.InitRepo[models.RequirementSkeleton](db, "requirement_skeletons", requirementSkeletonsSQL)
	if err != nil {
		return nil, err
	}
	return &requirementSkeletonsRepo{BaseRepo: repo}, nil
}
