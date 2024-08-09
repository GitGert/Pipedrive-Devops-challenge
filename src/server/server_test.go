package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	utils "github.com/GitGert/Pipedrive-Devops-challenge/src/utils"
)

func TestGetDeals(t *testing.T) {
	utils.LoadEnvFile("./../../.env")
	API_TOKEN = os.Getenv("API_TOKEN")
	COMPANY_DOMAIN = os.Getenv("COMPANY_DOMAIN")

	t.Setenv("COMPANY_DOMAIN", COMPANY_DOMAIN)
	t.Setenv("API_TOKEN", API_TOKEN)

	server := httptest.NewServer(http.HandlerFunc(getDeals))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", resp.StatusCode)
	}
}
