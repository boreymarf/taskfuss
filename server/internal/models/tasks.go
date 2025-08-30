package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type TaskSkeleton struct {
	ID      int64  `db:"id"`
	OwnerID int64  `db:"owner_id"`
	Status  string `db:"status"` // Archived, active
}

type TaskSnapshot struct {
	RevisionUUID uuid.UUID      `db:"revision_uuid"`
	TaskID       int64          `db:"task_id"`
	Title        string         `db:"title"`
	Description  sql.NullString `db:"description"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	IsCurrent    bool           `db:"is_current"`
}

type TaskPeriod struct {
	ID        int64        `db:"id"`
	TaskID    int64        `db:"task_id"`
	StartDate sql.NullTime `db:"start_date"`
	EndDate   sql.NullTime `db:"end_date"`
}

type TaskEntry struct {
	ID        int64     `db:"id"`
	TaskID    int64     `db:"task_id"`
	EntryDate time.Time `db:"entry_date"`
	Completed bool      `db:"completed"`
}

type RequirementSkeleton struct {
	ID     int64 `db:"id"`
	TaskID int64 `db:"task_id"`
}

type RequirementSnapshot struct {
	RevisionUUID uuid.UUID      `db:"revision_uuid"`
	SkeletonID   int64          `db:"skeleton_id"`
	ParentID     sql.NullInt64  `db:"parent_id"`
	Title        string         `db:"title"`
	Type         string         `db:"type"`         // atom or condition
	DataType     sql.NullString `db:"data_type"`    // int, time, bool, none, etc.
	Operator     sql.NullString `db:"operator"`     // or, not, and, ==, >=, <=, !=, >, < and etc.
	TargetValue  sql.NullString `db:"target_value"` // any value that needs to be parsed using DataType field
	SortOrder    int            `db:"sort_order"`
}

type RequirementEntry struct {
	ID            int64     `db:"id"`
	RequirementID int64     `db:"requirement_id"`
	EntryDate     time.Time `db:"entry_date"`
	Value         string    `db:"value"` // any value that needs to be parsed using DataType field of RequirementSnapshot
}
