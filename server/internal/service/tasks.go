package service

import (
	"sort"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
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

	if req.Task.Title == "" {
		return apperrors.NewValidationError("EMPTY_FIELD", "title", "Field 'title' cannot be empty")
	}
	// TODO: This can fail I think
	if req.Task.Requirement == nil {
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

	if err := s.addRequirement(req.Task.Requirement, task.ID, nil); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) addRequirement(requirement *dto.Requirement, task_id int64, parent_id *int64) error {

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
		for _, operand := range requirement.Operands {
			s.addRequirement(&operand, task_id, &r.ID)
		}
	}

	return nil

}

func (s *TaskService) GetAllTasks(opts *db.GetAllTasksOptions) ([]dto.Task, error) {

	// Get tasks
	modelTasks, err := s.taskRepo.GetAllTasks(opts)
	if err != nil {
		logger.Log.Err(err).Msg("Failed to get all tasks!")
		return nil, err
	}

	// Get requirements
	var tasksIDs []int64
	for _, modelTask := range modelTasks {
		tasksIDs = append(tasksIDs, modelTask.ID)
	}
	modelRequirements, err := s.requirementRepo.GetRequirementsByTaskIDs(tasksIDs)
	if err != nil {
		return nil, err
	}

	// TODO: Check later if req can be a pointer instead
	modelRequirementsByTask := make(map[int64][]models.Requirement)
	for _, req := range modelRequirements {
		modelRequirementsByTask[req.TaskID] = append(modelRequirementsByTask[req.TaskID], req)
	}

	var result []dto.Task
	for _, modelTask := range modelTasks {
		dtoTask := dto.Task{
			ID:          modelTask.ID,
			Title:       modelTask.Title,
			Description: modelTask.Description,
		}

		// Add requirement if exists
		if reqs, exists := modelRequirementsByTask[modelTask.ID]; exists {
			dtoTask.Requirement = buildTree(reqs, modelTask.ID)
		}

		if modelTask.CreatedAt.Valid {
			dtoTask.CreatedAt = &modelTask.CreatedAt.Time
		}
		if modelTask.UpdatedAt.Valid {
			dtoTask.UpdatedAt = &modelTask.UpdatedAt.Time
		}
		if modelTask.StartDate.Valid {
			dtoTask.StartDate = &modelTask.StartDate.Time
		}
		if modelTask.EndDate.Valid {
			dtoTask.EndDate = &modelTask.EndDate.Time
		}

		result = append(result, dtoTask)
	}

	return result, nil

}

func buildTree(modelRequirements []models.Requirement, taskID int64) *dto.Requirement {
	// A map for quick access by ID
	nodeMap := make(map[int64]*models.Requirement)
	// A map for quick access by ParentID
	childrenMap := make(map[int64][]*models.Requirement)

	var root *models.Requirement

	// Fill the maps with pointers
	for i := range modelRequirements {
		req := &modelRequirements[i]
		nodeMap[req.ID] = req

		if req.ParentID == nil {
			if root != nil {
				logger.Log.Warn().Int64("Task ID", taskID).Msg("Multiple root requirements found")
			}
			root = req
		} else {
			parentID := *req.ParentID
			childrenMap[parentID] = append(childrenMap[parentID], req)
		}
	}

	if root == nil {
		return nil
	}

	var convert func(*models.Requirement) *dto.Requirement
	convert = func(r *models.Requirement) *dto.Requirement {
		dtoReq := &dto.Requirement{
			ID:        r.ID,
			Title:     r.Title,
			Type:      r.Type,
			SortOrder: r.SortOrder,
		}

		if r.Type == "condition" && r.Operator != nil {
			dtoReq.Operator = r.Operator
		} else if r.Type == "atom" {
			if r.DataType != nil {
				dtoReq.DataType = r.DataType
			}
			if r.Operator != nil {
				dtoReq.Operator = r.Operator
			}
			if r.TargetValue != nil {
				dtoReq.TargetValue = r.TargetValue
			}
		}

		if children, exists := childrenMap[r.ID]; exists {
			sort.Slice(children, func(i, j int) bool {
				return children[i].SortOrder < children[j].SortOrder
			})

			for _, child := range children {
				dtoReq.Operands = append(dtoReq.Operands, *convert(child))
			}
		}

		return dtoReq
	}

	return convert(root)
}
