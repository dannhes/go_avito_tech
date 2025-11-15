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
	ErrUpdateFailed = errors.New("failed to update entity")
)

type UserRepository struct {
	pool *pgxpool.Pool
}

const (
	saveUserQuery  = `INSERT INTO users (id, username, team_name, is_active, created_at) VALUES ($1, $2, $3, $4, $5)`
	findUserQuery  = `SELECT id, username, team_name, is_active FROM users WHERE id = $1`
	setActiveQuery = `UPDATE users SET is_active = $1 WHERE id = $2`
	listUsersQuery = `SELECT id, username, team_name, is_active FROM users WHERE team_name = $1`
)

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Save(ctx context.Context, user domain.User) error {
	_, err := r.pool.Exec(ctx, saveUserQuery, user.ID, user.Username, user.TeamName, user.IsActive, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSaveFailed, err)
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	row := r.pool.QueryRow(ctx, findUserQuery, id)
	if err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		return domain.User{}, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}
	return u, nil
}

func (r *UserRepository) SetActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, setActiveQuery, active, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}
	return nil
}

func (r *UserRepository) ListByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	ans := make([]domain.User, 0)
	rows, err := r.pool.Query(ctx, listUsersQuery, teamName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		er := rows.Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive)
		if er != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanFailed, err)
		}
		ans = append(ans, user)
	}
	return ans, nil
}
