package dto

import "time"

type Task struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Requirement *Requirement `json:"requirement,omitempty"`
	Description *string      `json:"description,omitempty"` // Nullable
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty"`
	StartDate   *time.Time   `json:"start_date,omitempty"`
	EndDate     *time.Time   `json:"end_date,omitempty"` // Nullable
}

type Requirement struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	DataType    *string       `json:"data_type,omitempty"`
	Operator    *string       `json:"operator,omitempty"`
	TargetValue *string       `json:"target_value,omitempty"`
	Value       *string       `json:"value,omitempty"`
	Operands    []Requirement `json:"operands,omitempty"`
	SortOrder   int           `json:"sort_order"`
}

type TaskAddRequest struct {
	Task Task `json:"task"`
}

type GetTaskByIDResponse struct {
	Task Task `json:"task"`
}

type GetAllTasksResponse struct {
	Tasks []Task `json:"tasks"`
}
