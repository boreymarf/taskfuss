package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type TaskSkeletons interface {
	core.Repository[models.TaskSkeleton]
	core.Creator[models.TaskSkeleton]
	core.Getter[models.TaskSkeleton]
	core.Deleter[models.TaskSkeleton]
	core.AccessChecker
}

type taskSkeletonsRepo struct {
	*core.BaseRepo[models.TaskSkeleton]
}

const TaskSkeletonsSQL = `
	CREATE TABLE IF NOT EXISTS task_skeletons (
		id              INTEGER NOT NULL PRIMARY KEY,
		owner_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		status          VARCHAR(255) NOT NULL DEFAULT 'active' CHECK(status IN ('archived', 'active'))
	)`

func InitTaskSkeletons(db *sqlx.DB) (TaskSkeletons, error) {
	repo, err := core.InitRepo[models.TaskSkeleton](db, "task_skeletons", TaskSkeletonsSQL)
	if err != nil {
		return nil, err
	}
	return &taskSkeletonsRepo{BaseRepo: repo}, nil
}
