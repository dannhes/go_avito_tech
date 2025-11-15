package postgres_test

// TODO make a db
import (
	"context"
	"go_avito_tech/internal/domain"
	"go_avito_tech/internal/repository/postgres"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

const testDatabaseURL = "postgres://postgres:password@localhost:5432/test_db?sslmode=disable"

func setupTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testDatabaseURL)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}

	// Очистка таблицы users перед тестами
	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate users table: %v", err)
	}

	return pool
}

func TestUserRepository_Save(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	user := domain.User{
		ID:       "user1",
		Username: "testuser",
		TeamName: "teamA",
		IsActive: true,
	}

	err := repo.Save(ctx, user)
	assert.Nil(t, err)
}

func TestUserRepository_FindByID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	user := domain.User{
		ID:       "user1",
		Username: "testuser",
		TeamName: "teamA",
		IsActive: true,
	}

	err := repo.Save(ctx, user)
	assert.Nil(t, err)

	got, err := repo.FindByID(ctx, "user1")
	assert.Nil(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Username, got.Username)
	assert.Equal(t, user.TeamName, got.TeamName)
	assert.Equal(t, user.IsActive, got.IsActive)
}

func TestUserRepository_SetActive(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	user := domain.User{
		ID:       "user1",
		Username: "testuser",
		TeamName: "teamA",
		IsActive: true,
	}

	err := repo.Save(ctx, user)
	assert.Nil(t, err)

	err = repo.SetActive(ctx, "user1", false)
	assert.Nil(t, err)

	got, err := repo.FindByID(ctx, "user1")
	assert.Nil(t, err)
	assert.False(t, got.IsActive)
}

func TestUserRepository_ListByTeam(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	users := []domain.User{
		{ID: "user1", Username: "u1", TeamName: "teamA", IsActive: true},
		{ID: "user2", Username: "u2", TeamName: "teamA", IsActive: false},
		{ID: "user3", Username: "u3", TeamName: "teamB", IsActive: true},
	}

	for _, u := range users {
		assert.Nil(t, repo.Save(ctx, u))
	}

	listA, err := repo.ListByTeam(ctx, "teamA")
	assert.Nil(t, err)
	assert.Len(t, listA, 2)

	listB, err := repo.ListByTeam(ctx, "teamB")
	assert.Nil(t, err)
	assert.Len(t, listB, 1)
	assert.Equal(t, "user3", listB[0].ID)
}
