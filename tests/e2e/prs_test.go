package e2e

import (
	"encoding/json"
	"fmt"
	"go_avito_tech/api/gen"
	"net/http"
	"testing"
)

func TestPrCreateAndMerge(t *testing.T) {
	resp := doPost(t, "/team/add", map[string]any{
		"team_name": "dev",
		"members": []map[string]any{
			{"user_id": "222", "username": "Alice", "is_active": true},
		},
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create team/author: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	resp = doPost(t, "/pullRequest/create", map[string]any{
		"pull_request_id":   "522",
		"pull_request_name": "Feature X",
		"author_id":         "222",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create PR: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var pr gen.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		t.Errorf("failed to decode PR response: %v", err)
	}
	fmt.Println(pr.PullRequestId)
	if pr.PullRequestId != "522" {
		t.Errorf("wrong PR id: %s", pr.PullRequestId)
	}
	if pr.Status != gen.PullRequestStatusOPEN {
		t.Errorf("wrong PR status: %s", pr.Status)
	}
	resp = doPost(t, "/pullRequest/merge", map[string]any{
		"pull_request_id": "522",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("failed to merge PR: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var mergedPr gen.PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&mergedPr); err != nil {
		t.Errorf("failed to decode merged PR: %v", err)
	}
	if mergedPr.Status != gen.PullRequestStatusMERGED {
		t.Errorf("PR not merged, status: %s", mergedPr.Status)
	}
}
