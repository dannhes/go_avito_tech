package e2e

import (
	"context"
	"fmt"
	"go_avito_tech/api/gen"
	"go_avito_tech/internal/gateways/http"
	"go_avito_tech/internal/repository/postgres"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var (
	pool   *pgxpool.Pool
	server *echo.Echo
)

//nolint:errcheck
func TestMain(m *testing.M) {
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL not set, fallback to localhost:5433")
		dbURL = "postgres://test:test@localhost:5433/testdb?sslmode=disable"
	}
	var err error
	for i := 0; i < 15; i++ {
		pool, err = pgxpool.New(ctx, dbURL)
		if err == nil {
			break
		}
		fmt.Println("DB not ready, retrying:", err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}
	if err := initDB(ctx, pool); err != nil {
		log.Fatalf("cannot init DB: %v", err)
	}
	teamsRepo := postgres.NewTeamRepository(pool)
	usersRepo := postgres.NewUserRepository(pool)
	prRepo := postgres.NewPullRequestRepository(pool)
	statsRepo := postgres.NewPgStatsRepository(pool)
	useCases := http.UseCases{
		Users:  usersRepo,
		Teams:  teamsRepo,
		PullRs: prRepo,
		Stats:  statsRepo,
	}
	server = echo.New()
	h := http.NewHandler(useCases)
	gen.RegisterHandlers(server, h)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	go func() {
		if err := server.Start(":" + port); err != nil {
			log.Fatalf("server start failed: %v", err)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	code := m.Run()
	server.Shutdown(ctx)
	pool.Close()
	os.Exit(code)
}

func initDB(ctx context.Context, pool *pgxpool.Pool) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teams (
		name TEXT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT now()
	);
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT now()
	);
	CREATE TABLE IF NOT EXISTS pull_requests (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')) DEFAULT 'OPEN',
		need_more_reviewers BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT now(),
		merged_at TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS pull_request_reviewers (
		pr_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		PRIMARY KEY (pr_id, user_id)
	);
	`
	_, err := pool.Exec(ctx, schema)
	return err
}
