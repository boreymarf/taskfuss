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
	requirementSkeletons *db.RequirementSkeletons
	requirementSnapshots *db.RequirementSnapshots
	requirementEntries   *db.RequirementEntries
}

func InitTaskService(
	db *sql.DB,
	userRepo *db.Users,
	taskSkeletons db.TaskSkeletons,
	taskSnapshots db.TaskSnapshots,
	taskPeriods *db.TaskPeriods,
	taskEntries *db.TaskEntries,
	requirementSkeletons *db.RequirementSkeletons,
	requirementSnapshots *db.RequirementSnapshots,
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

	logger.Log.Debug().Msg("Trying to Create new task")

	if !(req.Requirement.Type == "atom" || req.Requirement.Type == "condition") {

		logger.Log.Debug().Str("req.Requirement.Type", req.Requirement.Type).Msg("Incorrect requirement type")
		// TODO: Add error later
		return nil, nil

	}

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

	// Revision UUID
	revision_uuid := uuid.New()

	// TaskSnapshot
	taskSnapshot := models.TaskSnapshot{
		Title:       req.Title,
		Description: utils.ToNullString(req.Description),
	}
	createdTaskSnapshot, err := s.taskSnapshots.CreateTx(ctx, tx, &taskSnapshot, revision_uuid)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create new task skeleton!")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	createdTask := dto.TaskResponse{
		ID:          createdTaskSkeleton.ID,
		Title:       createdTaskSnapshot.Title,
		Description: &createdTaskSnapshot.Description.String,
	}

	return &createdTask, nil

}
