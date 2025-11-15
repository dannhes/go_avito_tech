package e2e

import (
	"encoding/json"
	"go_avito_tech/internal/domain"
	"net/http"
	"testing"
)

func TestPrCreateAndFetch(t *testing.T) {
	resp := doPost(t, "/pull_requests", map[string]any{
		"id":                  "u123",
		"name":                "Roma",
		"author_id":           "123",
		"status":              domain.StatusOpen,
		"need_more_reviewers": true,
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	respNew, err := http.Get(baseURL + "/pull_requests/u123")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer respNew.Body.Close()
	if respNew.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", respNew.StatusCode, http.StatusOK)
	}
	var temp struct {
		ID                string                   `json:"id"`
		Username          string                   `json:"name"`
		AuthorId          string                   `json:"author_id"`
		Status            domain.PullRequestStatus `json:"status"`
		NeedMoreReviewers bool                     `json:"need_more_reviewers"`
	}
	json.NewDecoder(respNew.Body).Decode(&temp)
	if temp.ID != "u123" {
		t.Errorf("wrong id: %s", temp.ID)
	}
	if temp.Username != "Roma" {
		t.Errorf("wrong username: %s", temp.Username)
	}
	if temp.AuthorId != "123" {
		t.Errorf("wrong author id: %s", temp.AuthorId)
	}
	if temp.Status != domain.StatusOpen {
		t.Errorf("wrong status: %s", temp.Status)
	}
	if temp.NeedMoreReviewers != true {
		t.Errorf("wrong need more reviewers: %v", temp.NeedMoreReviewers)
	}
}
