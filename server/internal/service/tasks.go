package service

import (
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
)

type TaskService struct {
	userRepo             *db.Users
	taskSkeletons        db.TaskSkeletons
	taskSnapshots        *db.TaskSnapshots
	taskPeriods          *db.TaskPeriods
	taskEntries          *db.TaskEntries
	requirementSkeletons *db.RequirementSkeletons
	requirementSnapshots *db.RequirementSnapshots
	requirementEntries   *db.RequirementEntries
}

func InitTaskService(
	userRepo *db.Users,
	taskSkeletons db.TaskSkeletons,
	taskSnapshots *db.TaskSnapshots,
	taskPeriods *db.TaskPeriods,
	taskEntries *db.TaskEntries,
	requirementSkeletons *db.RequirementSkeletons,
	requirementSnapshots *db.RequirementSnapshots,
	requirementEntries *db.RequirementEntries,

) (*TaskService, error) {
	return &TaskService{
		userRepo:             userRepo,
		taskSkeletons:        taskSkeletons,
		taskSnapshots:        taskSnapshots,
		taskPeriods:          taskPeriods,
		taskEntries:          taskEntries,
		requirementSkeletons: requirementSkeletons,
		requirementSnapshots: requirementSnapshots,
		requirementEntries:   requirementEntries,
	}, nil
}

func (s *TaskService) CreateTask(req *dto.CreateTaskRequest, user_id int64) (*dto.TaskResponse, error) {

	logger.Log.Debug().Msg("Trying to Create new task")

	return nil, nil

}
