package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/utils"
)

type EntriesService struct {
	db                   *sql.DB
	taskSkeletons        db.TaskSkeletons
	taskSnapshots        db.TaskSnapshots
	requirementSkeletons db.RequirementSkeletons
	requirementSnapshots db.RequirementSnapshots
	requirementEntries   db.RequirementEntries
}

func InitEntriesService(
	db *sql.DB,
	taskSkeletons db.TaskSkeletons,
	taskSnapshots db.TaskSnapshots,
	requirementSkeletons db.RequirementSkeletons,
	requirementSnapshots db.RequirementSnapshots,
	requirementEntries db.RequirementEntries,
) (*EntriesService, error) {
	return &EntriesService{
		db:                   db,
		taskSkeletons:        taskSkeletons,
		taskSnapshots:        taskSnapshots,
		requirementSkeletons: requirementSkeletons,
		requirementSnapshots: requirementSnapshots,
		requirementEntries:   requirementEntries,
	}, nil
}

// TODO:
// 1) Проверка пользователя и его доступ к задаче и требованию
// 2) Проверка типа входного данного с типом требования
func (s *EntriesService) UpsertRequirementEntry(ctx context.Context, entry *models.RequirementEntry, userID int64) (*models.RequirementEntry, error) {

	// Checking if user can access the requirement
	requirementSkeleton, err := s.requirementSkeletons.GetByID(entry.RequirementID)
	if err != nil {
		return nil, fmt.Errorf("requirement not found: %w", err)
	}

	taskSkeleton, err := s.taskSkeletons.GetByID(requirementSkeleton.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	if taskSkeleton.OwnerID != userID {
		return nil, fmt.Errorf("access denied: user doesn't own this task")
	}

	// Checking if the requirement is atom
	taskSnapshot, err := s.taskSnapshots.GetEarliestFromDate(ctx, taskSkeleton.ID, entry.EntryDate)
	if err != nil {
		return nil, err
	}

	if taskSnapshot == nil {
		return nil, fmt.Errorf("service failed to get the earliest task snapshot!")
	}

	requirementSnapshot, err := s.requirementSnapshots.GetByCompositeKey(ctx, taskSnapshot.RevisionUUID, requirementSkeleton.ID)
	if err != nil {
		return nil, err
	}

	if requirementSnapshot.Type == "condition" {
		logger.Log.Error().Int64("requirementSkeleton.ID", requirementSkeleton.ID).Msg("Can't update entry for condition type requirement directly")
		return nil, fmt.Errorf("can't update entry for condition type requirement directly!")
	}

	// Inserting RevisionUUID
	entry.RevisionUUID = taskSnapshot.RevisionUUID

	// Checking if value type is correct
	if err := utils.ValidateValueByDataType(entry.Value, requirementSnapshot.DataType.String); err != nil {
		return nil, fmt.Errorf("invalid value for requirement: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	createdEntry, err := s.requirementEntries.UpsertTx(ctx, tx, entry)
	if err != nil {
		return nil, err
	}

	if requirementSnapshot.ParentID.Valid {
		parentSnapshotRequirement, err := s.requirementSnapshots.GetByCompositeKey(ctx, requirementSnapshot.RevisionUUID, requirementSnapshot.ParentID.Int64)
		if err != nil {
			return nil, err
		}

		_, err = s.updateConditionRequirement(ctx, tx, entry, parentSnapshotRequirement)
		if err != nil {
			return nil, err
		}
	}

	// End
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// TODO: I DON'T KNOW HOW I SUPPOSED TO RETURN THIS CRAP LMAO
	return createdEntry, nil
}

// Updates the *passed* requirementSnapshot and it's parent by chain
func (s *EntriesService) updateConditionRequirement(
	ctx context.Context,
	tx *sql.Tx,
	entry *models.RequirementEntry,
	requirementSnapshot *models.RequirementSnapshot,
) (*models.RequirementEntry, error) {

	if requirementSnapshot.Type != "condition" {
		logger.Log.Error().Str("requirementSnapshot.Type", requirementSnapshot.Type).Msg("Unexpected not condition at updateConditionRequirement!")
		return nil, fmt.Errorf("unexpected not condition at updateConditionRequirement!")
	}

	children, err := s.requirementSnapshots.GetChildren(ctx, requirementSnapshot.RevisionUUID, requirementSnapshot.SkeletonID)
	if err != nil {
		return nil, err
	}

	childrenIDs := make([]int64, len(children))
	for i, child := range children {
		childrenIDs[i] = child.SkeletonID
	}

	childrenEntries, err := s.requirementEntries.GetByRequirementIDsTx(ctx, tx, childrenIDs, requirementSnapshot.RevisionUUID)
	if err != nil {
		return nil, err
	}

	childrenResults := make([]any, len(children))
	for i, child := range children {
		childEntry, exists := childrenEntries[child.SkeletonID]
		if !exists {
			// Return false by default if no entry exists
			childrenResults[i] = false
			continue
		}

		result, err := utils.EvaluateAtomRequirement(&child, &childEntry)
		if err != nil {
			return nil, err
		}
		childrenResults[i] = result
	}

	parentResult, err := utils.EvaluateCondition(requirementSnapshot.Operator.String, childrenResults)
	if err != nil {
		return nil, err
	}

	parentEntry := models.RequirementEntry{
		RevisionUUID:  requirementSnapshot.RevisionUUID,
		RequirementID: requirementSnapshot.SkeletonID,
		EntryDate:     entry.EntryDate,
		Value:         strconv.FormatBool(parentResult),
	}

	s.requirementEntries.UpsertTx(ctx, tx, &parentEntry)

	if requirementSnapshot.ParentID.Valid {
		parentSnapshotRequirement, err := s.requirementSnapshots.GetByCompositeKey(ctx, requirementSnapshot.RevisionUUID, requirementSnapshot.ParentID.Int64)
		if err != nil {
			return nil, err
		}

		_, err = s.updateConditionRequirement(ctx, tx, entry, parentSnapshotRequirement)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
