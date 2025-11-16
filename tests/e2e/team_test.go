package e2e

import (
	"encoding/json"
	"fmt"
	"go_avito_tech/internal/domain"
	"net/http"
	"testing"
)

func TestTeamCreateAndFetch(t *testing.T) {
	resp := doPost(t, "/team/add", map[string]any{
		"team_name": "ctitmo",
		"members": []map[string]any{
			{"user_id": "43", "username": "Roma", "is_active": true},
			{"user_id": "42", "username": "Anna", "is_active": true},
		},
	})
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("resp.StatusCode = %d; want %d", resp.StatusCode, http.StatusCreated)
	}
	defer resp.Body.Close()
	respNew, err := http.Get(baseURL + "/team/get?team_name=ctitmo")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer respNew.Body.Close()
	if respNew.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", respNew.StatusCode, http.StatusOK)
	}
	var temp struct {
		TeamName string        `json:"name"`
		Members  []domain.User `json:"members"`
	}
	err = json.NewDecoder(respNew.Body).Decode(&temp)
	if err != nil {
		return
	}
	fmt.Println("team created:", temp.TeamName)
	if temp.TeamName != "ctitmo" {
		t.Errorf("wrong team_name: %s", temp.TeamName)
	}
}
