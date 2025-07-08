package dto

import "time"

type Task struct {
	Title       string          `json:"title"`
	Requirement TaskRequirement `json:"requirement"`
	Description *string         `json:"description"` // Nullable
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     *time.Time      `json:"end_date"` // Nullable
}

type TaskRequirement struct {
	Title       string             `json:"title"`
	Type        string             `json:"type"`
	DataType    *string            `json:"data_type"`
	Operator    *string            `json:"operator"`
	TargetValue *string            `json:"target_value"`
	Value       *string            `json:"value"`
	Operands    *[]TaskRequirement `json:"operands"`
	SortOrder   int                `json:"sort_order"`
}

func (t TaskRequirement) IsEmpty() bool {
	return t == TaskRequirement{}
}

type TaskAddRequest struct {
	Task Task `json:"task"`
}

type GetTaskByID struct {
	Task Task `json:"task"`
}
