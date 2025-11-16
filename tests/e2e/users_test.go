package e2e

import (
	"net/http"
	"testing"
)

func TestUserCreateAndFetch(t *testing.T) {
	resp := doPost(t, "/team/add", map[string]any{
		"team_name": "itmo",
		"members": []map[string]any{
			{"user_id": "112", "username": "Roma", "is_active": true},
		},
	})
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("failed to add user: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	resp = doPost(t, "/users/setIsActive", map[string]any{
		"user_id":   "112",
		"is_active": false,
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("failed to set active: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
}
