package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_avito_tech/internal/domain"
)

type PgStatsRepository struct {
	pool *pgxpool.Pool
}

func NewPgStatsRepository(pool *pgxpool.Pool) *PgStatsRepository {
	return &PgStatsRepository{pool: pool}
}

const (
	countUsersQuery = `SELECT 
    		COUNT(*),
			SUM(CASE WHEN is_active THEN 1 ELSE 0 END) FROM users;`
	countPrQuery = `SELECT
			COUNT(*),
            SUM(CASE WHEN status = 'open' THEN 1 ELSE 0 END),
            SUM(CASE WHEN status = 'merged' THEN 1 ELSE 0 END)
        	FROM pull_requests;`
	countRevsPerUserQuery = `SELECT 
    		reviewer_id, COUNT(*) 
        	FROM reviews GROUP BY reviewer_id;`
	countPrsPerUserQuery = `SELECT 
    		author_id, COUNT(*) 
        	FROM pull_requests GROUP BY author_id;`
	countTeamsQuery = `SELECT
			COUNT(*) FROM teams;`
)

func (r *PgStatsRepository) GetStats(ctx context.Context) (domain.Stats, error) {
	stats := domain.Stats{
		ReviewsPerUser: make(map[int]int),
		PRsPerAuthor:   make(map[int]int),
	}
	row := r.pool.QueryRow(ctx, countUsersQuery)
	row.Scan(&stats.TotalUsers, &stats.ActiveUsers)
	row = r.pool.QueryRow(ctx, countPrQuery)
	row.Scan(&stats.TotalPRs, &stats.OpenPRs, &stats.ClosedPRs)
	rows, _ := r.pool.Query(ctx, countRevsPerUserQuery)
	for rows.Next() {
		var id, count int
		rows.Scan(&id, &count)
		stats.ReviewsPerUser[id] = count
	}
	rows2, _ := r.pool.Query(ctx, countPrsPerUserQuery)
	for rows2.Next() {
		var id, count int
		rows2.Scan(&id, &count)
		stats.PRsPerAuthor[id] = count
	}
	row = r.pool.QueryRow(ctx, countTeamsQuery)
	row.Scan(&stats.TotalTeams)
	return stats, nil
}
