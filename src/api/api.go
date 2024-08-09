package api

import (
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Fetch makes an HTTP GET request to the specified URL.
func (s *Service) Fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{} // Use a default client for simplicity
	return client.Do(req)
}

// func (s *)
