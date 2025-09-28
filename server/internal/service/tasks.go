package service

import (
	"context"
	"database/sql"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/google/uuid"
)

type TaskService struct {
	db                   *sql.DB
	userRepo             db.Users
	taskSkeletons        db.TaskSkeletons
	taskSnapshots        db.TaskSnapshots
	taskPeriods          *db.TaskPeriods
	taskEntries          *db.TaskEntries
	requirementSkeletons db.RequirementSkeletons
	requirementSnapshots db.RequirementSnapshots
	requirementEntries   db.RequirementEntries
}

func InitTaskService(
	db *sql.DB,
	userRepo db.Users,
	taskSkeletons db.TaskSkeletons,
	taskSnapshots db.TaskSnapshots,
	taskPeriods *db.TaskPeriods,
	taskEntries *db.TaskEntries,
	requirementSkeletons db.RequirementSkeletons,
	requirementSnapshots db.RequirementSnapshots,
	requirementEntries db.RequirementEntries,

) (*TaskService, error) {
	return &TaskService{
		db:                   db,
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

func (s *TaskService) CreateTask(ctx context.Context, req *dto.CreateTaskRequest, userID int64) (*dto.TaskResponse, error) {
	revisionUUID := uuid.New()

	logger.Log.Debug().Msg("Trying to Create new task")

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// TaskSkeleton
	taskSkeleton := models.TaskSkeleton{
		OwnerID: userID,
	}
	createdTaskSkeleton, err := s.taskSkeletons.WithTx(tx).Create(ctx, &taskSkeleton)
	if err != nil {
		return nil, err
	}

	// TaskSnapshot
	taskSnapshot := models.TaskSnapshot{
		Title:       req.Title,
		Description: utils.ToNullString(req.Description),
	}
	createdTaskSnapshot, err := s.taskSnapshots.CreateTx(ctx, tx, &taskSnapshot, revisionUUID, createdTaskSkeleton.ID)
	if err != nil {
		return nil, err
	}
	err = s.taskSnapshots.SetCurrentRevisionTx(ctx, tx, createdTaskSkeleton.ID, revisionUUID)
	if err != nil {
		return nil, err
	}

	// Requirements
	createdRequirement, err := s.createRequirementTx(
		ctx,
		tx,
		req.Requirement,
		createdTaskSkeleton.ID,
		nil,
		revisionUUID,
	)
	if err != nil {
		return nil, err
	}

	// End
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	createdTask := dto.TaskResponse{
		ID:           createdTaskSkeleton.ID,
		Title:        createdTaskSnapshot.Title,
		RevisionUUID: revisionUUID,
		Description:  &createdTaskSnapshot.Description.String,
		Requirement:  createdRequirement,
	}

	return &createdTask, nil

}

func (s *TaskService) createRequirementTx(
	ctx context.Context,
	tx *sql.Tx,
	req *dto.CreateRequirementRequest,
	taskID int64,
	parentID *int64,
	revisionUUID uuid.UUID,
) (*dto.RequirementResponse, error) {

	requirementSkeleton := models.RequirementSkeleton{
		TaskID: taskID,
	}

	createdRequirementSkeleton, err := s.requirementSkeletons.CreateTx(ctx, tx, &requirementSkeleton)
	if err != nil {
		return nil, err
	}

	requirementSnapshot := models.RequirementSnapshot{
		RevisionUUID: revisionUUID,
		SkeletonID:   createdRequirementSkeleton.ID,
		ParentID:     utils.ToNullInt64(parentID),
		Title:        req.Title,
		Type:         req.Type,
		DataType:     utils.ToNullString(req.DataType),
		Operator:     utils.ToNullString(req.Operator),
		TargetValue:  req.TargetValue,
		SortOrder:    req.SortOrder,
	}

	createdRequirementSnapshot, err := s.requirementSnapshots.CreateTx(ctx, tx, &requirementSnapshot, revisionUUID)
	if err != nil {
		return nil, err
	}

	var requirementOperands []dto.RequirementResponse

	if req.Type == "condition" {
		for _, operand := range req.Operands {
			createdOp, err := s.createRequirementTx(ctx, tx, &operand, taskID, &createdRequirementSkeleton.ID, revisionUUID)
			if err != nil {
				return nil, err
			}
			requirementOperands = append(requirementOperands, *createdOp)
		}
	}

	createdRequirement := dto.RequirementResponse{
		ID:          createdRequirementSkeleton.ID,
		Title:       createdRequirementSnapshot.Title,
		Type:        createdRequirementSnapshot.Type,
		DataType:    req.DataType,
		Operator:    &createdRequirementSnapshot.Operator.String,
		TargetValue: createdRequirementSnapshot.TargetValue,
		Operands:    requirementOperands,
		SortOrder:   createdRequirementSnapshot.SortOrder,
	}

	return &createdRequirement, nil
}

type GetAllTasksQueryParams struct {
	DetailLevel   string
	ShowActive    bool
	ShowArchived  bool
	ShowCompleted bool
}

func (s *TaskService) GetAllTasks(ctx context.Context, params *GetAllTasksQueryParams, userID int64) (*[]dto.TaskResponse, error) {

	taskSkeletons, err := s.taskSkeletons.GetAll(ctx, params.ShowActive, params.ShowArchived)
	if err != nil {
		return nil, err
	}

	taskSkeletonsIDs := make([]int64, 0, len(taskSkeletons))
	for _, ts := range taskSkeletons {
		taskSkeletonsIDs = append(taskSkeletonsIDs, ts.ID)
	}

	taskSnapshots, err := s.taskSnapshots.GetAllLatest(ctx, taskSkeletonsIDs)
	if err != nil {
		return nil, err
	}

	requirementSkeletons, err := s.requirementSkeletons.GetByTaskIDs(taskSkeletonsIDs)
	if err != nil {
		return nil, err
	}

	// Group requirement skeletons by task id
	reqSkeletonsByTaskID := make(map[int64][]*models.RequirementSkeleton)
	for _, rs := range requirementSkeletons {
		reqSkeletonsByTaskID[rs.TaskID] = append(reqSkeletonsByTaskID[rs.TaskID], rs)
	}

	taskRevisions := make(map[int64]uuid.UUID)
	for _, ts := range taskSkeletons {
		for _, sn := range taskSnapshots {
			if sn.SkeletonID == ts.ID {
				taskRevisions[ts.ID] = sn.RevisionUUID
				break
			}
		}
	}

	// Fetch requirement snapshots using composite keys
	reqSnapshotMap := make(map[int64]*models.RequirementSnapshot)
	for taskID, reqSkeletons := range reqSkeletonsByTaskID {
		revisionUUID, exists := taskRevisions[taskID]
		if !exists {
			continue // Skip if no revision UUID found for this task
		}

		for _, rs := range reqSkeletons {
			snapshot, err := s.requirementSnapshots.GetByCompositeKey(ctx, revisionUUID, rs.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					continue // Skip if snapshot not found
				}
				return nil, err
			}
			reqSnapshotMap[rs.ID] = snapshot
		}
	}

	tasks := make([]dto.TaskResponse, 0, len(taskSkeletons))
	for _, ts := range taskSkeletons {
		// Находим снапшот для текущей задачи
		var taskSnapshot *models.TaskSnapshot
		for _, sn := range taskSnapshots {
			if sn.SkeletonID == ts.ID {
				taskSnapshot = sn
				break
			}
		}

		if taskSnapshot == nil {
			continue // Пропускаем задачи без снапшота
		}

		// Формируем ответ для задачи
		taskResponse := dto.TaskResponse{
			ID:           ts.ID,
			Title:        taskSnapshot.Title,
			RevisionUUID: taskSnapshot.RevisionUUID,
			Description:  &taskSnapshot.Description.String,
		}

		// Добавляем требования если они есть
		if reqSkeletons, exists := reqSkeletonsByTaskID[ts.ID]; exists && len(reqSkeletons) > 0 {
			// Collect snapshots for these skeletons
			reqSnapshots := make([]*models.RequirementSnapshot, 0, len(reqSkeletons))
			for _, rs := range reqSkeletons {
				if snapshot, ok := reqSnapshotMap[rs.ID]; ok {
					reqSnapshots = append(reqSnapshots, snapshot)
				}
			}

			// Build requirements tree
			requirements, err := s.buildRequirementsTree(reqSkeletons, reqSnapshots)
			if err != nil {
				return nil, err
			}

			if len(requirements) > 0 {
				taskResponse.Requirement = &requirements[0]
			}
		}

		tasks = append(tasks, taskResponse)
	}

	return &tasks, nil
}

func (s *TaskService) GetTask(ctx context.Context, taskID int64, userID int64) (*dto.TaskResponse, error) {
	// Get specific task skeleton
	taskSkeleton, err := s.taskSkeletons.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Get latest snapshot for this task
	taskSnapshot, err := s.taskSnapshots.GetLatest(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Get requirement skeletons for this task
	requirementSkeletons, err := s.requirementSkeletons.GetByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// Get requirement snapshots
	reqSnapshots := make([]*models.RequirementSnapshot, 0, len(requirementSkeletons))
	if len(requirementSkeletons) > 0 {
		for _, rs := range requirementSkeletons {
			snapshot, err := s.requirementSnapshots.GetByCompositeKey(ctx, taskSnapshot.RevisionUUID, rs.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					continue
				}
				return nil, err
			}
			reqSnapshots = append(reqSnapshots, snapshot)
		}
	}

	// Build response
	taskResponse := dto.TaskResponse{
		ID:           taskSkeleton.ID,
		Title:        taskSnapshot.Title,
		RevisionUUID: taskSnapshot.RevisionUUID,
		Description:  &taskSnapshot.Description.String,
	}

	// Build requirements tree if they exist
	if len(requirementSkeletons) > 0 {
		requirementResponses, err := s.buildRequirementsTree(requirementSkeletons, reqSnapshots)
		if err != nil {
			return nil, err
		}
		if len(requirementResponses) > 0 {
			taskResponse.Requirement = &requirementResponses[0]
		}
	}

	return &taskResponse, nil
}

func (s *TaskService) buildRequirementsTree(
	skeletons []*models.RequirementSkeleton,
	snapshots []*models.RequirementSnapshot,
) ([]dto.RequirementResponse, error) {
	reqMap := make(map[int64]*dto.RequirementResponse)
	var rootRequirements []*dto.RequirementResponse

	for i, rs := range skeletons {
		if i >= len(snapshots) || snapshots[i] == nil {
			continue // Skip if no snapshot
		}
		snapshot := snapshots[i]
		if snapshot == nil {
			continue
		}

		req := &dto.RequirementResponse{
			ID:          rs.ID,
			Title:       snapshot.Title,
			Type:        snapshot.Type,
			DataType:    &snapshot.DataType.String,
			Operator:    &snapshot.Operator.String,
			TargetValue: snapshot.TargetValue,
			SortOrder:   snapshot.SortOrder,
		}

		reqMap[rs.ID] = req

		if !snapshot.ParentID.Valid {
			rootRequirements = append(rootRequirements, req)
		}
	}

	// Затем добавляем операнды к родительским требованиям
	for i, rs := range skeletons {
		if i >= len(snapshots) {
			continue
		}

		snapshot := snapshots[i]
		if snapshot == nil {
			continue
		}

		if snapshot.ParentID.Valid {
			if parent, exists := reqMap[snapshot.ParentID.Int64]; exists {
				if child, exists := reqMap[rs.ID]; exists {
					parent.Operands = append(parent.Operands, *child)
				}
			}
		}
	}

	// Конвертируем корневые требования в slice
	result := make([]dto.RequirementResponse, 0, len(rootRequirements))
	for _, req := range rootRequirements {
		result = append(result, *req)
	}

	return result, nil
}
