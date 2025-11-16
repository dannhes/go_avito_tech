package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_avito_tech/internal/domain"
	"time"
)

var (
	ErrTeamNotFound = errors.New("team not found")
	ErrUserNotFound = errors.New("user not found")
	ErrScanFailed   = errors.New("failed to scan row")
	ErrSaveFailed   = errors.New("failed to save entity")
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

const (
	saveTeamQuery   = `INSERT INTO teams (name, created_at) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`
	getMembersQuery = `SELECT id, username, team_name, is_active FROM users WHERE team_name = $1`
	getTeamQuery    = `SELECT name, created_at FROM teams WHERE name = $1`
)

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) Save(ctx context.Context, team string) error {
	_, err := r.pool.Exec(ctx, saveTeamQuery, team, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSaveFailed, err)
	}
	return nil
}

func (r *TeamRepository) FindByName(ctx context.Context, name string) (domain.Team, error) {
	row := r.pool.QueryRow(ctx, getTeamQuery, name)
	var createdAt time.Time
	err := row.Scan(&name, &createdAt)
	if err != nil {
		return domain.Team{}, fmt.Errorf("%w: %v", ErrTeamNotFound, err)
	}
	users := make([]domain.User, 0)
	rows, err := r.pool.Query(ctx, getMembersQuery, name)
	if err != nil {
		return domain.Team{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}
	defer rows.Close()
	for rows.Next() {
		user := domain.User{}
		if er := rows.Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive); er != nil {
			return domain.Team{}, fmt.Errorf("%w: %v", ErrScanFailed, er)
		}
		users = append(users, user)
	}
	team := domain.Team{Name: name, Members: users}
	return team, nil
}
