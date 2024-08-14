package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/GitGert/Pipedrive-Devops-challenge/src/api"
	models "github.com/GitGert/Pipedrive-Devops-challenge/src/models"
	utils "github.com/GitGert/Pipedrive-Devops-challenge/src/utils"
)

func InitServer() {

	r := mux.NewRouter()

	pipedriveApi := api.NewService()

	getDealsFunc := HandlerWithService(pipedriveApi, getDeals)
	postDealsFunc := HandlerWithService(pipedriveApi, postDeals)
	putDealsFunc := HandlerWithService(pipedriveApi, putDeals)

	r.HandleFunc("/deals", getDealsFunc).Methods("GET")
	r.HandleFunc("/deals", postDealsFunc).Methods("POST")
	r.HandleFunc("/deals", putDealsFunc).Methods("PUT")
	r.HandleFunc("/metrics", getMetrics).Methods("GET")

	fmt.Println("server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}

type ApiService interface {
	GetDeals() (*http.Response, error)
	AddDeal(dealData models.PostDeal) (*http.Response, error)
	ModifyDeal(data models.PatchDeal, dealID string) (*http.Response, error)
}

type handlerFunc func(w http.ResponseWriter, r *http.Request, service ApiService)

// HandlerWithService wraps the original handler and injects the service
func HandlerWithService(service ApiService, fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, service)
	}
}

func getDeals(w http.ResponseWriter, r *http.Request, service ApiService) {
	if r.Method != http.MethodGet {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	// pipedriveApi := api.NewService()
	// resp, err := pipedriveApi.GetDeals()

	resp, err := service.GetDeals()
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	utils.Log_request(r, "statusCode: 200")
}

func postDeals(w http.ResponseWriter, r *http.Request, service ApiService) {
	if r.Method != http.MethodPost {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	var dealData models.PostDeal
	err := json.NewDecoder(r.Body).Decode(&dealData)
	if err != nil {
		fmt.Println(err)
		httpErrorHandler(w, "Bad request", http.StatusBadRequest, r)
		return
	}

	// pipedriveApi := api.NewService()
	resp, err := service.AddDeal(dealData)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	utils.Log_request(r, "statusCode: 200")
}

func putDeals(w http.ResponseWriter, r *http.Request, service ApiService) {
	if r.Method != http.MethodPut {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	dealID := r.URL.Query().Get("dealId")
	if dealID == "" {
		fmt.Println("empty")
		httpErrorHandler(w, "Bad request", http.StatusBadRequest, r)
		return
	}

	var data models.PatchDeal
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		httpErrorHandler(w, "Bad request", http.StatusBadRequest, r)
		return
	}

	// pipedriveApi := api.NewService()
	resp, err := service.ModifyDeal(data, dealID)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	log.Println(string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	utils.Log_request(r, "statusCode: 200")
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	dealID := "1"

	getDealsLatency, err := timeEndpoint("http://localhost:8080/deals", http.MethodGet, nil)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	postData := models.PostDeal{
		Title:    "testing Post endpoint latency",
		Value:    "3000",
		Currency: "EUR",
	}
	postDataJSON, err := json.Marshal(postData)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Err: %s", err)
	}
	postDealsLatency, err := timeEndpoint("http://localhost:8080/deals", http.MethodPost, bytes.NewBuffer(postDataJSON))
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	putData := models.PatchDeal{
		Title:    "testing PUT endpoint latency",
		Value:    5000,
		Currency: "USD",
	}
	putDataJSON, err := json.Marshal(putData)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Err: %s", err)
	}

	putDealsLatency, err := timeEndpoint("http://localhost:8080/deals?dealId="+dealID, http.MethodPut, bytes.NewBuffer(putDataJSON))
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	requestMetrics := models.RequestMetrics{
		GetDeals:  getDealsLatency,
		PostDeals: postDealsLatency,
		PutDeals:  putDealsLatency,
	}

	jsonData, err := json.Marshal(requestMetrics)
	if err != nil {
		httpErrorHandler(w, err.Error(), http.StatusInternalServerError, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	utils.Log_request(r, "statusCode: 200")
}

func timeEndpoint(url string, HttpMethod string, requestBody *bytes.Buffer) (string, error) {
	start := time.Now()
	fmt.Println(requestBody)
	if requestBody == nil {
		requestBody = &bytes.Buffer{}
	}

	client := &http.Client{}
	req, err := http.NewRequest(HttpMethod, url, requestBody)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		elapsed := float64(time.Since(start).Milliseconds()) / 1000
		returnString := fmt.Sprint(elapsed)
		return returnString, nil
	} else {
		return "", fmt.Errorf("error with api call")
	}
}

func httpErrorHandler(w http.ResponseWriter, message string, statusMethod int, r *http.Request) {
	utils.Log_event(utils.MakeRed(utils.GetUrl(r)), utils.MakeRed("statusCode: "+strconv.Itoa(statusMethod)))
	http.Error(w, message, statusMethod)
}
