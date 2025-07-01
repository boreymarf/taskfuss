package models

import "time"

type Task struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Requirement *TaskRequirement `json:"task_requiremenet"`
	Description *string          `json:"description,omitempty"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
}

type TaskRequirement struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}
