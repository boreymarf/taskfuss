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
func (s *EntriesService) UpsertRequirementEntry(ctx context.Context, entry *models.RequirementEntry, userID int64) ([]models.Node, error) {
	var nodes []models.Node

	requirementSkeleton, err := s.requirementSkeletons.GetByID(entry.RequirementID)
	if err != nil {
		return nil, fmt.Errorf("requirement not found: %w", err)
	}

	taskSkeleton, err := s.taskSkeletons.GetByID(requirementSkeleton.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

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

	// Checking if user can access the requirement
	if taskSkeleton.OwnerID != userID {
		return nil, fmt.Errorf("access denied: user doesn't own this task")
	}

	// Checking if the requirement is atom
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

	node := models.Node{
		Content: createdEntry,
		ID:      requirementSkeleton.ID,
	}

	if requirementSnapshot.ParentID.Valid {
		node.Parent = requirementSnapshot.ParentID.Int64
	} else {
		node.Parent = 0
	}

	nodes = append(nodes, node)

	if requirementSnapshot.ParentID.Valid {
		parentSnapshotRequirement, err := s.requirementSnapshots.GetByCompositeKey(ctx, requirementSnapshot.RevisionUUID, requirementSnapshot.ParentID.Int64)
		if err != nil {
			return nil, err
		}

		nodes, err = s.updateConditionChain(ctx, tx, entry, parentSnapshotRequirement, nodes)
		if err != nil {
			return nil, err
		}
	}

	// End
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nodes, nil
}

// updateConditionChain updates a chain of requirement entries starting from the given entry and
// moving upwards through its parent entries until it encounters a requirement without a parent
// (where RequirementSnapshot.ParentID is nil). The function uses the provided context and
// transaction to ensure the changes are applied within the transaction scope.
//
// Parameters:
// - ctx: The context for the operation, typically used for cancellation or deadlines.
// - tx: The transaction that wraps the changes made to the database, ensuring atomicity of updates.
// - entry: The initial requirement entry to start the update chain. This entry itself has not been updated yet.
// - requirementSnapshot: A snapshot of the requirement that the entry belongs to, providing data about its parent-child relationship.
// - nodes: A slice of nodes that holds all updated *models.RequirementEntry from the current entry upwards to the root entry.
//
// The function recursively calls itself to update the parent entries, ensuring the entire chain of entries is updated
// until a requirement is reached that does not have a parent (i.e., where ParentID is nil).
//
// Returns:
//   - A slice of updated nodes, containing all the updated *models.RequirementEntry from the current entry
//     up to the root entry, reflecting the changes made during the recursive updates.
func (s *EntriesService) updateConditionChain(
	ctx context.Context,
	tx *sql.Tx,
	entry *models.RequirementEntry,
	requirementSnapshot *models.RequirementSnapshot,
	nodes []models.Node,
) ([]models.Node, error) {

	logger.Log.Debug().Msg("WORKING ON IT BOSS")

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

	createdEntry, err := s.requirementEntries.UpsertTx(ctx, tx, &parentEntry)
	if err != nil {
		return nil, err
	}

	node := models.Node{
		Content: createdEntry,
		ID:      requirementSnapshot.SkeletonID,
	}

	if requirementSnapshot.ParentID.Valid {
		node.Parent = requirementSnapshot.ParentID.Int64
	} else {
		node.Parent = 0
	}

	nodes = append(nodes, node)

	if requirementSnapshot.ParentID.Valid {
		parentSnapshotRequirement, err := s.requirementSnapshots.GetByCompositeKey(ctx, requirementSnapshot.RevisionUUID, requirementSnapshot.ParentID.Int64)
		if err != nil {
			return nil, err
		}

		nodes, err = s.updateConditionChain(ctx, tx, entry, parentSnapshotRequirement, nodes)
		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}
