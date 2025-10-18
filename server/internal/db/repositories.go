package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Users                Users
	TaskSkeletons        TaskSkeletons
	TaskSnapshots        TaskSnapshots
	TaskEntries          TaskEntries
	TaskPeriods          TaskPeriods
	RequirementSkeletons RequirementSkeletons
	RequirementSnapshots RequirementSnapshots
	RequirementEntries   RequirementEntries
}

func InitializeRepositories(database *sqlx.DB) (*Repositories, error) {
	repos := &Repositories{}

	var err error
	repos.Users, err = InitUsers(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize user repository: %w", err)
	}

	repos.TaskSkeletons, err = InitTaskSkeletons(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize task skeletons repository: %w", err)
	}

	repos.TaskSnapshots, err = InitTaskSnapshots(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize task snapshots repository: %w", err)
	}

	repos.TaskEntries, err = InitTaskEntries(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize task entry repository: %w", err)
	}

	repos.TaskPeriods, err = InitTaskPeriods(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize task periods repository: %w", err)
	}

	repos.RequirementSkeletons, err = InitRequirementSkeletons(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize requirement skeleton repository: %w", err)
	}

	repos.RequirementSnapshots, err = InitRequirementSnapshots(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize requirement snapshots repository: %w", err)
	}

	repos.RequirementEntries, err = InitRequirementEntries(database)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize requirement entry repository: %w", err)
	}

	return repos, nil
}
