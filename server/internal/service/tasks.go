package service

import (
	"fmt"
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

func (s *TaskService) CreateTask(req *dto.TaskCreateRequest, user_id int64) (*models.Task, error) {

	logger.Log.Debug().Msg("Trying to Create new task")

	if req.Task.Title == "" {
		return nil, apperrors.NewValidationError("EMPTY_FIELD", "title", "Field 'title' cannot be empty")
	}
	// TODO: This can fail I think
	if req.Task.Requirement == nil {
		return nil, apperrors.NewValidationError("EMPTY_FIELD", "requirement", "Field 'requirement' cannot be empty")
	}

	task := models.Task{
		OwnerID:     user_id,
		Title:       req.Task.Title,
		Description: req.Task.Description,
	}

	createdTask, err := s.taskRepo.CreateTask(task)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to Create new task")
		return nil, err
	}

	logger.Log.Debug().Msg("Now trying to Create requirements of the task...")

	if err := s.CreateRequirement(req.Task.Requirement, task.ID, nil); err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (s *TaskService) CreateRequirement(requirement *dto.Requirement, task_id int64, parent_id *int64) error {

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
	if err := s.requirementRepo.CreateRequirement(&r); err != nil {
		return err
	}

	logger.Log.Debug().Str("title", requirement.Title).Msg("Created requirement to the db!")

	if requirement.Type == "condition" {
		for _, operand := range requirement.Operands {
			s.CreateRequirement(&operand, task_id, &r.ID)
		}
	}

	return nil

}

// FIXME: Service should not return DTO, I'll fix it later
func (s *TaskService) GetTaskByID(taskID int64, userID int64) (dto.Task, error) {

	modelTask, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return dto.Task{}, nil
	}

	if modelTask.OwnerID != userID {
		return dto.Task{}, apperrors.ErrForbidden
	}

	var modelRequirements []models.Requirement
	modelRequirements, err = s.requirementRepo.GetRequirementsByTaskIDs([]int64{taskID})
	if err != nil {
		return dto.Task{}, err
	}

	dtoTask := dto.Task{
		ID:          modelTask.ID,
		Title:       modelTask.Title,
		Description: modelTask.Description,
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

	dtoRequirement, err := buildTree(modelRequirements, taskID)
	if err != nil {
		return dto.Task{}, err
	}

	dtoTask.Requirement = dtoRequirement

	return dtoTask, nil
}

type GetAllTasksOptions struct {
	DetailLevel   string
	ShowActive    bool
	ShowArchived  bool
	ShowCompleted bool
}

// FIXME: Service shouldn't return dto, I'll fix it later
func (s *TaskService) GetAllTasks(opts *GetAllTasksOptions, userID int64) ([]dto.Task, error) {

	dbOpts := db.GetAllTasksOptions{
		DetailLevel:  opts.DetailLevel,
		ShowActive:   opts.ShowActive,
		ShowArchived: opts.ShowArchived,
		UserID:       userID,
	}

	// Get tasks
	modelTasks, err := s.taskRepo.GetAllTasks(&dbOpts)
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

	modelRequirementsByTask := make(map[int64][]models.Requirement)
	for i := range modelRequirements {
		req := modelRequirements[i]
		modelRequirementsByTask[req.TaskID] = append(modelRequirementsByTask[req.TaskID], req)
	}

	var result []dto.Task
	for _, modelTask := range modelTasks {
		dtoTask := dto.Task{
			ID:          modelTask.ID,
			Title:       modelTask.Title,
			Description: modelTask.Description,
		}

		// Create requirement if exists
		if reqs, exists := modelRequirementsByTask[modelTask.ID]; exists {
			dtoTask.Requirement, err = buildTree(reqs, modelTask.ID)
			if err != nil {
				return []dto.Task{}, err
			}
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

func buildTree(modelRequirements []models.Requirement, taskID int64) (*dto.Requirement, error) {
	// A map for quick access by ID
	nodeMap := make(map[int64]*models.Requirement)
	// A map for quick access by ParentID
	childrenMap := make(map[int64][]*models.Requirement)

	var root *models.Requirement

	// Fill the maps with pointers
	for _, req := range modelRequirements {
		nodeMap[req.ID] = &req

		if req.ParentID == nil {
			if root != nil {
				logger.Log.Error().Int64("taskID", taskID).Msg("Multiple requirement roots found!")
				return nil, fmt.Errorf("multiple root requirements found")
			}
			root = &req
		} else {
			parentID := *req.ParentID
			childrenMap[parentID] = append(childrenMap[parentID], &req)
		}
	}

	if root == nil {
		logger.Log.Error().Int64("taskID", taskID).Msg("No requirement roots found!")
		return nil, fmt.Errorf("no root requirement found")
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

	return convert(root), nil
}
