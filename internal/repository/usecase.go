package repository

import (
	"context"
	"go_avito_tech/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	FindByID(ctx context.Context, id string) (domain.User, error)
	SetActive(ctx context.Context, id string, active bool) error
	ListByTeam(ctx context.Context, teamName string) ([]domain.User, error)
	//TODO деактивация
}

type TeamRepository interface {
	Save(ctx context.Context, team string) error
	FindByName(ctx context.Context, name string) (domain.Team, error)
}

type PullRequestRepository interface {
	Save(ctx context.Context, pr domain.PullRequest) error
	FindByID(ctx context.Context, id string) (domain.PullRequest, error)
	FindByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error)

	AssignReviewers(ctx context.Context, prID string) ([]string, error)
	ReassignReviewer(ctx context.Context, prID string, oldUserID string) (string, error)
	Merge(ctx context.Context, prID string) (domain.PullRequest, error)
}

type StatsRepository interface {
	GetStats(ctx context.Context) (domain.Stats, error)
}
