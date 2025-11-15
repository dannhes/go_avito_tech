package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestTeamCreateAndFetch(t *testing.T) {
	resp := doPost(t, "/teams", map[string]any{
		"name": "Roma",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", resp.StatusCode, http.StatusOK)
	}
	defer resp.Body.Close()
	respNew, err := http.Get(baseURL + "/teams/Roma")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer respNew.Body.Close()
	if respNew.StatusCode != http.StatusOK {
		t.Errorf("resp.StatusCode = %d; want %d", respNew.StatusCode, http.StatusOK)
	}
	var temp struct {
		Name string `json:"name"`
	}
	json.NewDecoder(respNew.Body).Decode(&temp)
	if temp.Name != "Roma" {
		t.Errorf("wrong id: %s", temp.Name)
	}
}
