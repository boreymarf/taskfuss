package db

import (
	"context"
	"testing"

	"github.com/boreymarf/task-fuss/server/internal/db/core"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// --------------------
// Implementation
// --------------------

// Model
type Cat struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	OwnerID int64  `db:"owner_id"`
}

const catTableSQL = `
CREATE TABLE IF NOT EXISTS cats (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	owner_id INTEGER NOT NULL
)
`

// In-memory test repo
type Cats interface {
	core.Creator[Cat]
	core.Updater[Cat]
	core.Deleter[Cat]
	core.Getter[Cat]
	core.AccessChecker
}

type catsRepoImpl struct {
	*core.BaseRepo[Cat]
}

// InitCats initializes the cats repository and creates the table
func InitCats(db *sqlx.DB) (Cats, error) {
	repo, err := core.InitRepo[Cat](db, "cats", catTableSQL)
	if err != nil {
		return nil, err
	}
	return &catsRepoImpl{BaseRepo: repo}, nil
}

// --------------------
// Tests
// --------------------

func TestCatRepo(t *testing.T) {
	dbConn, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer dbConn.Close()

	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	repo, err := InitCats(dbConn)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	ctx := context.Background()

	// Predeclare cats with owners
	cat1 := &Cat{ID: 1, Name: "Whiskers", OwnerID: 10}
	cat2 := &Cat{ID: 2, Name: "Mittens", OwnerID: 10}
	cat3 := &Cat{ID: 3, Name: "Shadow", OwnerID: 11}

	// Predeclare user contexts
	adminCtx := &models.UserContext{ID: 999, Role: models.RoleAdmin}
	userCtx := &models.UserContext{ID: 10, Role: models.RoleUser}
	guestCtx := &models.UserContext{ID: 0, Role: models.RoleGuest}

	t.Run("Create cat", func(t *testing.T) {
		err := repo.Create(ctx, cat1)
		if err != nil {
			t.Errorf("create failed: %v", err)
		}
	})

	t.Run("Get cat by ID", func(t *testing.T) {
		cat, err := repo.Get(ctx, userCtx).With("id", cat1.ID).First()
		if err != nil {
			t.Errorf("get failed: %v", err)
		}
		if cat.Name != "Whiskers" {
			t.Errorf("expected Whiskers, got %s", cat.Name)
		}
	})

	t.Run("Get all cats by OwnerID", func(t *testing.T) {
		cats, err := repo.Get(ctx, userCtx).With("owner_id", cat1.OwnerID).All()
		if err != nil {
			t.Errorf("get all failed: %v", err)
		}
		if len(cats) == 0 {
			t.Errorf("expected at least 1 cat, got %d", len(cats))
		}
	})

	t.Run("Update cat", func(t *testing.T) {
		err := repo.Update(ctx, adminCtx).With("id", 1).Set("name", "Tiger").One()
		if err != nil {
			t.Errorf("update failed: %v", err)
		}
		cat, _ := repo.Get(ctx, adminCtx).With("id", 1).First()
		if cat.Name != "Tiger" {
			t.Errorf("expected Tiger, got %s", cat.Name)
		}
	})

	t.Run("Delete cat", func(t *testing.T) {
		err := repo.Delete(ctx, adminCtx).With("id", 1).One()
		if err != nil {
			t.Errorf("delete failed: %v", err)
		}
		cats, _ := repo.Get(ctx, adminCtx).All()
		if len(cats) != 0 {
			t.Errorf("expected 0 cats, got %d", len(cats))
		}
	})

	t.Run("Access control", func(t *testing.T) {
		// Seed remaining cats
		_ = repo.Create(ctx, cat2)
		_ = repo.Create(ctx, cat3)

		t.Run("admin can access all", func(t *testing.T) {
			ok, err := repo.CheckAccess(ctx, adminCtx, cat2.ID)
			if err != nil {
				t.Errorf("access check failed: %v", err)
			}
			if !ok {
				t.Errorf("admin should have access")
			}
			ok, _ = repo.CheckAccess(ctx, adminCtx, cat3.ID)
			if !ok {
				t.Errorf("admin should have access")
			}
		})

		t.Run("user can access own cat", func(t *testing.T) {
			ok, err := repo.CheckAccess(ctx, userCtx, cat2.ID)
			if err != nil {
				t.Errorf("access check failed: %v", err)
			}
			if !ok {
				t.Errorf("user should have access to own cat")
			}
			ok, _ = repo.CheckAccess(ctx, userCtx, cat3.ID)
			if ok {
				t.Errorf("user should not have access to another cat")
			}
		})

		t.Run("guest cannot access any cat", func(t *testing.T) {
			ok, err := repo.CheckAccess(ctx, guestCtx, cat2.ID)
			if err != nil {
				t.Errorf("access check failed: %v", err)
			}
			if ok {
				t.Errorf("guest should not have access")
			}
			ok, _ = repo.CheckAccess(ctx, guestCtx, cat3.ID)
			if ok {
				t.Errorf("guest should not have access")
			}
		})
	})
}
