package db_test

import (
	"context"
	"testing"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	dbConn, err := sqlx.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	require.NotNil(t, dbConn)
	return dbConn
}

func TestRequirementSkeletons(t *testing.T) {
	ctx := context.Background()
	dbConn := setupTestDB(t)
	defer dbConn.Close()

	repo, err := db.InitRequirementSkeletons(dbConn)
	require.NoError(t, err)

	// fake user context
	user := &models.UserContext{
		ID:   1,
		Role: models.RoleAdmin,
	}

	// test Create
	skeleton := &models.RequirementSkeleton{TaskID: 123}
	created, err := repo.Create(ctx, skeleton)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.Equal(t, int64(123), created.TaskID)
	require.NotZero(t, created.ID)

	// test query by ID
	query := repo.Get(user).WithIDs(created.ID)
	results, err := query.Send(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, created.ID, results[0].ID)
	require.Equal(t, created.TaskID, results[0].TaskID)

	// test query by TaskID
	query = repo.Get(user).WithTaskIDs(created.TaskID)
	results, err = query.Send(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, created.TaskID, results[0].TaskID)

	// test multiple inserts and IN query
	s2, err := repo.Create(ctx, &models.RequirementSkeleton{TaskID: 456})
	require.NoError(t, err)
	s3, err := repo.Create(ctx, &models.RequirementSkeleton{TaskID: 789})
	require.NoError(t, err)

	query = repo.Get(user).WithIDs(created.ID, []int64{s2.ID, s3.ID})
	results, err = query.Send(ctx)
	require.NoError(t, err)
	require.Len(t, results, 3)

	// test non-existing filter
	query = repo.Get(user).WithIDs(int64(99999))
	results, err = query.Send(ctx)
	require.NoError(t, err)
	require.Len(t, results, 0)
}
