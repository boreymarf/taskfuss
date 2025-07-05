package service

import (
	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/sanity-io/litter"
)

type TaskService struct {
	taskRepo             *db.TaskRepository
	taskEntryRepo        *db.TaskEntryRepository
	requirementRepo      *db.RequirementRepository
	requirementEntryRepo *db.RequirementEntryRepository
}

func InitTaskService(
	taskRepo *db.TaskRepository,
	taskEntryRepo *db.TaskEntryRepository,
	requirementRepo *db.RequirementRepository,
	requirementEntryRepo *db.RequirementEntryRepository,
) (*TaskService, error) {

	repo := &TaskService{
		taskRepo:             taskRepo,
		taskEntryRepo:        taskEntryRepo,
		requirementRepo:      requirementRepo,
		requirementEntryRepo: requirementEntryRepo,
	}

	return repo, nil
}

func (s *TaskService) AddTask(req *dto.TaskAddRequest, user_id int64) error {

	logger.Log.Debug().Msg("Trying to add new task")

	litter.Dump(req)

	if req.Task.Title == "" {
		return apperrors.NewValidationError("EMPTY_FIELD", "title", "Field 'title' cannot be empty")
	}
	if req.Task.Requirement.IsEmpty() {
		return apperrors.NewValidationError("EMPTY_FIELD", "requirement", "Field 'requirement' cannot be empty")
	}

	task := models.Task{
		OwnerID:     user_id,
		Title:       req.Task.Title,
		Description: req.Task.Description,
	}

	if err := s.taskRepo.AddTask(&task); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to add new task")
		return err
	}

	logger.Log.Debug().Msg("Now trying to add requirements of the task...")

	if err := s.addRequirement(&req.Task.Requirement, task.ID, nil); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) addRequirement(requirement *dto.TaskRequirement, task_id int64, parent_id *int64) error {

	r := models.Requirement{
		TaskID:      task_id,
		ParentID:    parent_id,
		Title:       requirement.Title,
		Type:        requirement.Type,
		DataType:    requirement.DataType,
		Operator:    requirement.Operator,
		TargetValue: requirement.TargetValue,
		Value:       requirement.Value,
		SortOrder:   requirement.SortOrder,
	}

	// Returns id after
	if err := s.requirementRepo.AddRequirement(&r); err != nil {
		return err
	}

	logger.Log.Debug().Str("title", requirement.Title).Msg("Added requirement to the db!")

	if requirement.Type == "condition" {
		for _, operand := range *requirement.Operands {
			s.addRequirement(&operand, task_id, &r.ID)
		}
	}

	return nil

}
