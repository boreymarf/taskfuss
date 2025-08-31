package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title       string                    `json:"title" binding:"required"`
	Description *string                   `json:"description,omitempty"`
	StartDate   *time.Time                `json:"start_date,omitempty"`
	EndDate     *time.Time                `json:"end_date,omitempty"`
	Requirement *CreateRequirementRequest `json:"requirement" binding:"required"`
}

type CreateRequirementRequest struct {
	Title       string                     `json:"title" binding:"required"`
	Type        string                     `json:"type" binding:"required"`         // atom or condition
	DataType    *string                    `json:"data_type,omitempty"`             // int, time, bool, none, etc.
	Operator    *string                    `json:"operator,omitempty"`              // or, not, and, ==, >=, <=, !=, >, < and etc.
	TargetValue string                     `json:"target_value" binding:"required"` // any value that needs to be parsed using DataType field
	Operands    []CreateRequirementRequest `json:"operands,omitempty"`
	SortOrder   int                        `json:"sort_order"`
}

type TaskResponse struct {
	ID           int64                `json:"id"`
	Title        string               `json:"title"`
	RevisionUUID uuid.UUID            `json:"revision_uuid"`
	Requirement  *RequirementResponse `json:"requirement"`
	Description  *string              `json:"description,omitempty"`
	CreatedAt    *time.Time           `json:"created_at,omitempty"`
	UpdatedAt    *time.Time           `json:"updated_at,omitempty"`
	StartDate    *time.Time           `json:"start_date,omitempty"`
	EndDate      *time.Time           `json:"end_date,omitempty"`
}

type RequirementResponse struct {
	ID          int64                 `json:"id"`
	Title       string                `json:"title"`
	Type        string                `json:"type"`
	DataType    *string               `json:"data_type,omitempty"`
	Operator    *string               `json:"operator,omitempty"`
	TargetValue string                `json:"target_value"`
	Operands    []RequirementResponse `json:"operands,omitempty"`
	SortOrder   int                   `json:"sort_order"`
}

type CreateTaskResponse struct {
	Task TaskResponse `json:"task"`
}

type GetTaskByIDResponse struct {
	Task TaskResponse `json:"task"`
}

type GetAllTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
