package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/GitGert/Pipedrive-Devops-challenge/src/api"
	"github.com/GitGert/Pipedrive-Devops-challenge/src/constants"
	models "github.com/GitGert/Pipedrive-Devops-challenge/src/models"
	utils "github.com/GitGert/Pipedrive-Devops-challenge/src/utils"
)

func InitServer() {

	r := mux.NewRouter()

	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/post_deals", postDeals).Methods("POST") //TODO: CHANGE TO POST
	r.HandleFunc("/put_deals", putDeals).Methods("PUT")    //TODO: CHANGE TO PUT
	r.HandleFunc("/metrics", getMetrics).Methods("GET")

	fmt.Println("server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// TODO: consider adding this kind of modularity into your code:       const response = await api.addDeal(data);
func getDeals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	pipedriveApi := api.NewService()
	resp, err := pipedriveApi.GetDeals()
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

func postDeals(w http.ResponseWriter, r *http.Request) {
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

	pipedriveApi := api.NewService()
	resp, err := pipedriveApi.AddDeal(dealData)
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

func putDeals(w http.ResponseWriter, r *http.Request) {
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

	pipedriveApi := api.NewService()
	resp, err := pipedriveApi.ModifyDeal(data, dealID)
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

	getDealsLatency := timeEndpoint("https://"+constants.COMPANY_DOMAIN+".pipedrive.com/api/v1/deals?limit=20&api_token="+constants.API_TOKEN, http.MethodGet)
	postDealsLatency := timeEndpoint("https://"+constants.COMPANY_DOMAIN+".pipedrive.com/api/v1/deals?api_token="+constants.API_TOKEN, http.MethodGet)           //TODO: change to post
	putDealsLatency := timeEndpoint("https://"+constants.COMPANY_DOMAIN+".pipedrive.com/api/v2/deals/"+dealID+"?api_token="+constants.API_TOKEN, http.MethodGet) //TODO: change to put

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

// TODO ADD ERROR HANDLING WHEN CALLING THIS FUNCTION.
func timeEndpoint(url string, HttpMethod string) string {
	start := time.Now()

	client := &http.Client{}
	req, err := http.NewRequest(HttpMethod, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return ""
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		elapsed := float64(time.Since(start).Milliseconds()) / 1000
		returnString := fmt.Sprint(elapsed)
		return returnString
	} else {
		fmt.Println("something went wrong")
		return ""
	}
}

func httpErrorHandler(w http.ResponseWriter, message string, statusMethod int, r *http.Request) {
	utils.Log_event(utils.MakeRed(utils.GetUrl(r)), utils.MakeRed("statusCode: "+strconv.Itoa(statusMethod)))
	http.Error(w, message, statusMethod)
}
