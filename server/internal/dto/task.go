package dto

import "time"

type Task struct {
	Name        string           `json:"name"`
	Requirement *TaskRequirement `json:"task_requiremenet"`
	Description *string          `json:"description,omitempty"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
}

type TaskRequirement struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type TaskCreateRequest struct {
	Task Task `json:"task"`
}
