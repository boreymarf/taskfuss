package core

import (
	"context"

	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository[T any] interface {
	WithTx(tx *sqlx.Tx) Repository[T]
}

// Optional extensions.
type Creator[T any] interface {
	Create(ctx context.Context, obj *T) error
}

type Getter[T any] interface {
	Get(ctx context.Context, uc *models.UserContext) *GetQuery[T]
}

type Updater[T any] interface {
	Update(ctx context.Context, uc *models.UserContext) *UpdateQuery[T]
}

type Deleter[T any] interface {
	Delete(ctx context.Context, uc *models.UserContext) *DeleteQuery[T]
}

type AccessChecker interface {
	CheckAccess(ctx context.Context, uc *models.UserContext, id any) (bool, error)
}
