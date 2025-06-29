package models

import "time"

type Task struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type TaskRequirement struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Completed bool   `json:"completed"`
}
