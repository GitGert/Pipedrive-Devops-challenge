package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	utils "github.com/GitGert/Pipedrive-Devops-challenge/src/utils"
)

// This test currently fails because I am not mocking the .evn file imports
// also I need to mock the request.
func TestGetDeals(t *testing.T) {
	utils.LoadEnvFile("./../../.env")

	server := httptest.NewServer(http.HandlerFunc(getDeals))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", resp.StatusCode)
	}
}
