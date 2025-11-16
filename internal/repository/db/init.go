package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func GetPool() *pgxpool.Pool {
	return pool
}

func InitDB(ctx context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/pr_service?sslmode=disable"
	}
	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("cannot create db pool: %w", err)
	}
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
	_, err = pool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("cannot create tables: %w", err)
	}
	return nil
}

func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}
