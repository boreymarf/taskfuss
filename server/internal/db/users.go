package db

import (
	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	core.Repository[models.User]
	core.Creator[models.User]
	core.Getter[models.User]
	core.Updater[models.User]
	core.Deleter[models.User]
}

type usersRepo struct {
	*core.BaseRepo[models.User]
}

const UsersSQL = `
	CREATE TABLE IF NOT EXISTS users (
		id 				INTEGER NOT NULL PRIMARY KEY,
		name 					VARCHAR(255) NOT NULL,
		email 				VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at 		DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at 		DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

func InitUsers(db *sqlx.DB) (Users, error) {
	repo, err := core.InitRepo[models.User](db, "users", UsersSQL)
	if err != nil {
		return nil, err
	}
	return &usersRepo{BaseRepo: repo}, nil
}

