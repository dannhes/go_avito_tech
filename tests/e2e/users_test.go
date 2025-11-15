package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestUserCreateAndFetch(t *testing.T) {
	resp := doPost(t, "/users", map[string]any{
		"id":        "u123",
		"username":  "Roma",
		"team_name": "backend",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	respNew, err := http.Get(baseURL + "/users/u123")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer respNew.Body.Close()
	if respNew.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", respNew.StatusCode, http.StatusOK)
	}
	var temp struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	json.NewDecoder(respNew.Body).Decode(&temp)
	if temp.ID != "u123" {
		t.Errorf("wrong id: %s", temp.ID)
	}
	if temp.Username != "Roma" {
		t.Errorf("wrong username: %s", temp.Username)
	}
}
