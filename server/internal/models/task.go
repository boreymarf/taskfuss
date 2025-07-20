package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int64        `json:"id"`
	OwnerID     int64        `json:"owner_id"`
	Title       string       `json:"title"`
	Description *string      `json:"description"`
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	StartDate   sql.NullTime `json:"start_date"`
	EndDate     sql.NullTime `json:"end_date"`
	Status      string       `json:"status"`
}

type TaskEntry struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	EntryDate time.Time `json:"entry_date"`
	Completed bool      `json:"completed"` // Stored as 0/1 in DB
}

type Requirement struct {
	ID          int64   `json:"id"`
	TaskID      int64   `json:"task_id"`
	ParentID    *int64  `json:"parent_id"` // Nullable
	Title       string  `json:"title"`
	Type        string  `json:"type"`         // ENUM: 'atom','and','or','not'
	DataType    *string `json:"data_type"`    // ENUM: 'bool','int','float','duration'
	Operator    *string `json:"operator"`     // Nullable
	TargetValue *string `json:"target_value"` // Nullable
	Value       *string `json:"value"`        // Nullable
	SortOrder   int     `json:"sort_order"`   // Nullable
}

type RequirementEntry struct {
	ID            int64     `json:"id"`
	RequirementID int64     `json:"requirement_id"`
	EntryDate     time.Time `json:"entry_date"`
	Value         string    `json:"value"`
}
