package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type taskPeriods interface {
	core.Repository[models.TaskPeriod]
	core.Creator[models.TaskPeriod]
	core.Getter[models.TaskPeriod]
	core.Updater[models.TaskPeriod]
	core.Deleter[models.TaskPeriod]
}

type taskPeriodsRepo struct {
	*core.BaseRepo[models.TaskPeriod]
}

const taskPeriodsSQL = `
	CREATE TABLE IF NOT EXISTS task_periods (
		id							INTEGER NOT NULL PRIMARY KEY,
		task_id					INTEGER NOT NULL REFERENCES task_skeletons(id) ON DELETE CASCADE,
		start_date			DATETIME NOT NULL,
		end_date				DATETIME,
		is_active				BOOLEAN DEFAULT TRUE,
		created_at			DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at			DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

func InittaskPeriods(db *sqlx.DB) (taskPeriods, error) {
	repo, err := core.InitRepo[models.TaskPeriod](db, "task_periods", taskPeriodsSQL)
	if err != nil {
		return nil, err
	}
	return &taskPeriodsRepo{BaseRepo: repo}, nil
}
