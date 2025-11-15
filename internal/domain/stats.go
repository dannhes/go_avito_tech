package domain

type Stats struct {
	TotalUsers     int         `json:"total_users"`
	ActiveUsers    int         `json:"active_users"`
	TotalPRs       int         `json:"total_prs"`
	OpenPRs        int         `json:"open_prs"`
	ClosedPRs      int         `json:"closed_prs"`
	ReviewsPerUser map[int]int `json:"reviews_per_user"`
	PRsPerAuthor   map[int]int `json:"prs_per_author"`
	TotalTeams     int         `json:"total_teams"`
}
