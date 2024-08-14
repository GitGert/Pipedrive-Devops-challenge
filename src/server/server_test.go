package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/GitGert/Pipedrive-Devops-challenge/src/models"
	"gotest.tools/v3/assert"
)

func TestGetDeals(t *testing.T) {

	mockService := NewMockAPI()

	mockGetDealsHandler := HandlerWithService(mockService, getDeals)

	server := httptest.NewServer(http.HandlerFunc(mockGetDealsHandler))

	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAddDeal(t *testing.T) {

	mockService := NewMockAPI()

	mockGetDealsHandler := HandlerWithService(mockService, postDeals)

	server := httptest.NewServer(http.HandlerFunc(mockGetDealsHandler))

	defer server.Close()

	payload := []byte(`{
		"title": "Example Deal",
		"value": "1000",
		"currency": "EUR",
		"status": "open"
	}`)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestModifyDeal(t *testing.T) {

	mockService := NewMockAPI()

	mockGetDealsHandler := HandlerWithService(mockService, putDeals)

	server := httptest.NewServer(http.HandlerFunc(mockGetDealsHandler))

	payload := []byte(`{
		"title": "Deal Title",
		"value": 1000.00,
		"currency": "EUR",
		"is_deleted": false,
		"status": "won"
	}`)

	req, err := http.NewRequest(http.MethodPut, server.URL+"/?dealId=1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

type Mock_api struct {
}

func NewMockAPI() *Mock_api {
	return &Mock_api{}
}

// TODO: change this
type Deal struct {
	ID   int
	Name string
}

func (m *Mock_api) GetDeals() (*http.Response, error) {

	deals := []Deal{
		{ID: 1, Name: "Deal One"},
		{ID: 2, Name: "Deal Two"},
	}

	jsonData, err := json.Marshal(deals)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
	}

	return resp, nil
}

func (m *Mock_api) AddDeal(dealData models.PostDeal) (*http.Response, error) {
	deals := []Deal{
		{ID: 1, Name: "Deal One"},
		{ID: 2, Name: "Deal Two"},
	}

	jsonData, err := json.Marshal(deals)
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
	}

	return resp, nil
}

func (m *Mock_api) ModifyDeal(data models.PatchDeal, dealID string) (*http.Response, error) {
	deals := []Deal{
		{ID: 1, Name: "Deal One"},
		{ID: 2, Name: "Deal Two"},
	}

	jsonData, err := json.Marshal(deals)
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
	}

	return resp, nil
}
