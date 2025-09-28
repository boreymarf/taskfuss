package models

import "time"

type TasksFilter struct {
	Active   *string
	Archived *string
}

type EntriesFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Dates     []time.Time
}
