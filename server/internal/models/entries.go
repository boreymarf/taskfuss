package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskEntry struct {
	ID        int64     `db:"id"`
	TaskID    int64     `db:"task_id"`
	EntryDate time.Time `db:"entry_date"`
	Completed bool      `db:"completed"`
}

type RequirementEntry struct {
	ID            int64     `db:"id"`
	RevisionUUID  uuid.UUID `db:"revision_uuid"`
	RequirementID int64     `db:"requirement_id"`
	EntryDate     time.Time `db:"entry_date"`
	Value         string    `db:"value"`
}

type Node struct {
	Content any
	ID      int64
	Parent  int64
}

type SignedRequirementEntry struct {
	Entry         RequirementEntry
	RequirementID int64
	ParentID      int64
	RevisionUUID  uuid.UUID
}
