package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/GitGert/Pipedrive-Devops-challenge/src/constants"
	"github.com/GitGert/Pipedrive-Devops-challenge/src/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetDeals() (*http.Response, error) {

	requestURL := "https://" + constants.COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals?limit=20&api_token=" + constants.API_TOKEN

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	return client.Do(req)
}

func (s *Service) AddDeal(dealData models.PostDeal) (*http.Response, error) {

	reqBody, err := json.Marshal(dealData)
	if err != nil {
		return nil, err
	}
	reqBodyBytes := bytes.NewBuffer(reqBody)

	requestURL := "https://" + constants.COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals?api_token=" + constants.API_TOKEN

	req, err := http.NewRequest(http.MethodPost, requestURL, reqBodyBytes)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	return client.Do(req)
}

func (s *Service) ModifyDeal(data models.PatchDeal, dealID string) (*http.Response, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %s", err.Error())
		return nil, err
	}
	reqBodyBytes := bytes.NewBuffer(reqBody)

	requestURL := "https://" + constants.COMPANY_DOMAIN + ".pipedrive.com/api/v2/deals/" + dealID + "?api_token=" + constants.API_TOKEN

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, requestURL, reqBodyBytes)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
