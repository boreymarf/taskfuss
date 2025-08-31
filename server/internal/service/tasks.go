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
	userRepo             *db.Users
	taskSkeletons        db.TaskSkeletons
	taskSnapshots        db.TaskSnapshots
	taskPeriods          *db.TaskPeriods
	taskEntries          *db.TaskEntries
	requirementSkeletons db.RequirementSkeletons
	requirementSnapshots db.RequirementSnapshots
	requirementEntries   *db.RequirementEntries
}

func InitTaskService(
	db *sql.DB,
	userRepo *db.Users,
	taskSkeletons db.TaskSkeletons,
	taskSnapshots db.TaskSnapshots,
	taskPeriods *db.TaskPeriods,
	taskEntries *db.TaskEntries,
	requirementSkeletons db.RequirementSkeletons,
	requirementSnapshots db.RequirementSnapshots,
	requirementEntries *db.RequirementEntries,

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

func (s *TaskService) CreateTask(ctx context.Context, req *dto.CreateTaskRequest, user_id int64) (*dto.TaskResponse, error) {
	revision_uuid := uuid.New()

	logger.Log.Debug().Msg("Trying to Create new task")

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create new transaction!")
		return nil, err
	}

	// TaskSkeleton
	taskSkeleton := models.TaskSkeleton{
		OwnerID: user_id,
	}
	createdTaskSkeleton, err := s.taskSkeletons.CreateTx(ctx, tx, &taskSkeleton)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create new task skeleton!")
		return nil, err
	}

	// TaskSnapshot
	taskSnapshot := models.TaskSnapshot{
		Title:       req.Title,
		Description: utils.ToNullString(req.Description),
	}
	createdTaskSnapshot, err := s.taskSnapshots.CreateTx(ctx, tx, &taskSnapshot, revision_uuid)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create new task skeleton!")
	}

	// Requirements
	createdRequirement, err := s.createRequirementTx(
		ctx,
		tx,
		req.Requirement,
		createdTaskSkeleton.ID,
		nil,
		revision_uuid,
	)
	if err != nil {
		return nil, err
	}

	// End
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	createdTask := dto.TaskResponse{
		ID:          createdTaskSkeleton.ID,
		Title:       createdTaskSnapshot.Title,
		Description: &createdTaskSnapshot.Description.String,
		Requirement: createdRequirement,
	}

	return &createdTask, nil

}

func (s *TaskService) createRequirementTx(
	ctx context.Context,
	tx *sql.Tx,
	req *dto.CreateRequirementRequest,
	task_id int64,
	parent_id *int64,
	revision_uuid uuid.UUID,
) (*dto.RequirementResponse, error) {

	requirementSkeleton := models.RequirementSkeleton{
		TaskID: task_id,
	}

	createdRequirementSkeleton, err := s.requirementSkeletons.CreateTx(ctx, tx, &requirementSkeleton)
	if err != nil {
		return nil, err
	}

	requirementSnapshot := models.RequirementSnapshot{
		RevisionUUID: revision_uuid,
		SkeletonID:   createdRequirementSkeleton.ID,
		ParentID:     utils.ToNullInt64(parent_id),
		Title:        req.Title,
		Type:         req.Type,
		DataType:     utils.ToNullString(req.DataType),
		Operator:     utils.ToNullString(req.Operator),
		TargetValue:  req.TargetValue,
		SortOrder:    req.SortOrder,
	}

	createdRequirementSnapshot, err := s.requirementSnapshots.CreateTx(ctx, tx, &requirementSnapshot, revision_uuid)
	if err != nil {
		return nil, err
	}

	var requirementOperands []dto.RequirementResponse

	if req.Type == "condition" {
		for _, operand := range req.Operands {
			createdOp, err := s.createRequirementTx(ctx, tx, &operand, task_id, &createdRequirementSkeleton.ID, revision_uuid)
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
