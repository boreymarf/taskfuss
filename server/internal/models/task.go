package models

import "time"

type Task struct {
	ID          int64            `json:"id"`
	OwnerID     int64            `json:"owner_id"`
	Name        string           `json:"name"`
	Requirement *TaskRequirement `json:"requirement"`
	Description *string          `json:"description,omitempty"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
}

type TaskRequirement struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}
