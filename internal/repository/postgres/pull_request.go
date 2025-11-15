package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_avito_tech/internal/domain"
	"math/rand"
	"time"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

var (
	ErrPRNotFound          = errors.New("pull request not found")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrPRMerged            = errors.New("cannot modify merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer not assigned to PR")
)

const (
	saveRequestQuery = `INSERT INTO pull_requests 
    	(id, name, author_id, status, need_more_reviewers, created_at)
    	VALUES ($1,$2,$3,$4,$5,$6)`
	findRequestQuery          = `SELECT * FROM pull_requests WHERE id = $1`
	findRequestReviewersQuery = `SELECT * FROM pull_request_reviewers WHERE id = $1`
	selectCandidateQuery      = `SELECT id FROM users WHERE team_name=(SELECT team_name 
                                 FROM users WHERE id=$1) 
		 						 AND id <> $1 AND is_active=TRUE`
	saveReviewerToPrQuery = `INSERT INTO pull_request_reviewers (pr_id, user_id)
							 VALUES ($1,$2)`
	updateReviewerQuery = `UPDATE pull_request_reviewers SET user_id=$1 WHERE pr_id=$2 AND user_id=$3`
	updatePrMerged      = `UPDATE pull_requests SET status='MERGED', merged_at = $1, where id = $2`
)

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool: pool}
}

func (r *PullRequestRepository) Save(ctx context.Context, pr domain.PullRequest) error {
	_, err := r.pool.Exec(ctx, saveRequestQuery, pr.ID, pr.Name, pr.Status, pr.NeedMoreReviewers, time.Now())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSaveFailed, err)
	}
	return nil
}

func (r *PullRequestRepository) FindByID(ctx context.Context, id string) (domain.PullRequest, error) {
	var pr domain.PullRequest
	err := r.pool.QueryRow(ctx, findRequestQuery, id).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status,
		&pr.NeedMoreReviewers)
	if err != nil {
		return pr, fmt.Errorf("%w: %v", ErrPRNotFound, err)
	}
	rows, err := r.pool.Query(ctx, findRequestReviewersQuery, id)
	if err != nil {
		return pr, fmt.Errorf("%w: %v", ErrPRMerged, err)
	}
	defer rows.Close()
	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return pr, fmt.Errorf("%w: %v", ErrReviewerNotAssigned, err)
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer)
	}
	return pr, nil
}

func (r *PullRequestRepository) FindByReviewer(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	rows, err := r.pool.Query(ctx, `SELECT * FROM pull_requests WHERE user_id=$1`, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReviewerNotAssigned, err)
	}
	defer rows.Close()
	var prsId []string
	for rows.Next() {
		var prId string
		if err := rows.Scan(&prId); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReviewerNotAssigned, err)
		}
		prsId = append(prsId, prId)
	}
	var prs []domain.PullRequest
	for _, val := range prsId {
		var pr domain.PullRequest
		pr, err = r.FindByID(ctx, val)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPRNotFound, err)
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func (r *PullRequestRepository) AssignReviewers(ctx context.Context, prID string) ([]string, error) {
	pr, err := r.FindByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == domain.StatusMerged {
		return nil, ErrPRMerged
	}
	rows, err := r.pool.Query(ctx, selectCandidateQuery, pr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviewers []string
	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReviewerNotAssigned, err)
		}
		reviewers = append(reviewers, reviewer)
	}
	rand.Shuffle(len(reviewers), func(i, j int) { reviewers[i], reviewers[j] = reviewers[j], reviewers[i] })
	n := 2
	if len(reviewers) < 2 {
		n = len(reviewers)
	}
	revs := reviewers[:n]
	for _, val := range revs {
		_, err := r.pool.Exec(ctx, saveReviewerToPrQuery, pr.ID, val)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSaveFailed, err)
		}
	}
	return revs, nil
}

func (r *PullRequestRepository) ReassignReviewer(ctx context.Context, prID string, oldUserID string) (string, error) {
	pr, err := r.FindByID(ctx, prID)
	if err != nil {
		return "", err
	}
	if pr.Status == domain.StatusMerged {
		return "", ErrPRMerged
	}
	// TODO validate
	var teamName string
	err = r.pool.QueryRow(ctx, `SELECT team_name FROM users WHERE id=$1`, oldUserID).Scan(&teamName)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoCandidate, err)
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id FROM users WHERE team_name=$1 AND id<>$2 AND is_active=TRUE`,
		teamName, oldUserID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoCandidate, err)
	}
	defer rows.Close()
	var candidates []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return "", fmt.Errorf("%w: %v", ErrNoCandidate, err)
		}
		candidates = append(candidates, uid)
	}
	if len(candidates) == 0 {
		return "", ErrNoCandidate
	}
	newReviewer := candidates[rand.Intn(len(candidates))]
	_, err = r.pool.Exec(ctx, updateReviewerQuery, newReviewer)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoCandidate, err)
	}
	return newReviewer, nil
}

func (r *PullRequestRepository) Merge(ctx context.Context, prID string) (domain.PullRequest, error) {
	pr, err := r.FindByID(ctx, prID)
	if err != nil {
		return pr, ErrPRNotFound
	}
	if pr.Status == domain.StatusMerged {
		return pr, nil
	}
	pr.Status = domain.StatusMerged
	_, err = r.pool.Exec(ctx, updatePrMerged, pr.ID, time.Now())
	if err != nil {
		return pr, fmt.Errorf("%w: %v", ErrPRMerged, err)
	}
	return pr, nil
}
