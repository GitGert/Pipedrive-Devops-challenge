package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	models "github.com/GitGert/Pipedrive-Devops-challenge/src/models"
	utils "github.com/GitGert/Pipedrive-Devops-challenge/src/utils"
)

var API_TOKEN string
var COMPANY_DOMAIN string

func InitServer() {
	utils.LoadEnvFile(".env")
	API_TOKEN = os.Getenv("API_TOKEN")
	COMPANY_DOMAIN = os.Getenv("COMPANY_DOMAIN")

	r := mux.NewRouter()

	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/post_deals", postDeals).Methods("GET") //TODO: CHANGE TO POST
	r.HandleFunc("/put_deals", putDeals).Methods("GET")   //TODO: CHANGE TO PUT
	r.HandleFunc("/metrics", getMetrics).Methods("GET")   //TODO: CHANGE TO PUT

	fmt.Println("server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// TODO: consider adding this kind of modularity into your code:       const response = await api.addDeal(data);
func getDeals(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	requestURL := "https://" + COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals?limit=20&api_token=" + API_TOKEN

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		// log.Printf("Failed to create request: %v", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("Failed to read response body: %v", err)
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
	//TODO: UNCOMMENT:
	// if r.Method != http.MethodPost {
	// 	httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
	// 	return
	// }

	dealData := models.DealData{
		Title:             "Deal of the century",
		Value:             10000,
		Currency:          "USD",
		UserID:            nil,
		PersonID:          nil,
		OrgID:             1,
		StageID:           1,
		Status:            "open",
		ExpectedCloseDate: "2022-02-11",
		Probability:       60,
		LostReason:        nil,
		VisibleTo:         1,
		AddTime:           "2021-02-11",
	}

	reqBody, err := json.Marshal(dealData)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %s", err.Error())
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	reqBodyBytes := bytes.NewBuffer(reqBody)

	// requestURL := "https://" + COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals?api_token=" + API_TOKEN

	req, err := http.NewRequest(http.MethodPost, "requestURL", reqBodyBytes)
	if err != nil {
		// log.Printf("Failed to create request: %v", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// log.Printf("Request failed: %v\n", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("Failed to read response body: %v", err)
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
	// utils.Log_request(r, "got call to put_deals")

	// TODO: unncomnet
	// if r.Method != http.MethodPut {
	// 	httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
	// 	return
	// }

	dealID := "1" //TODO: CHANGE THIS

	// data, err := json.Marshal(map[string]interface{}{
	// "user_id": userID_Gert,
	// })

	newTitle := "THIS TITLE VALUE HAS BEEN CHANGED BY THE API"

	data := models.PatchRequstData{
		// UserID:   userID_Gert,
		Title:    newTitle,
		Currency: "EUR",
		Value:    2000,
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %s", err.Error())
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}
	reqBodyBytes := bytes.NewBuffer(reqBody)

	requestURL := "https://" + COMPANY_DOMAIN + ".pipedrive.com/api/v2/deals/" + dealID + "?api_token=" + API_TOKEN

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, requestURL, reqBodyBytes)
	if err != nil {
		// log.Printf("Failed to create request: %v", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		// log.Printf("Request failed: %v", err)
		httpErrorHandler(w, "Internal server error", http.StatusInternalServerError, r)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Println("chaged deal: " + dealID + " title to: " + newTitle)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("failed to read response body %v", err)
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
	// utils.Log_request(r, "got call to metrics")
	if r.Method != http.MethodGet {
		httpErrorHandler(w, "Invalid request method", http.StatusMethodNotAllowed, r)
		return
	}

	dealID := "1"

	getDealsLatency := timeEndpoint("https://"+COMPANY_DOMAIN+".pipedrive.com/api/v1/deals?limit=20&api_token="+API_TOKEN, http.MethodGet)
	postDealsLatency := timeEndpoint("https://"+COMPANY_DOMAIN+".pipedrive.com/api/v1/deals?api_token="+API_TOKEN, http.MethodGet)           //TODO: change to post
	putDealsLatency := timeEndpoint("https://"+COMPANY_DOMAIN+".pipedrive.com/api/v2/deals/"+dealID+"?api_token="+API_TOKEN, http.MethodGet) //TODO: change to put

	requestMetrics := models.RequestMetrics{
		GetDeals:  getDealsLatency,
		PostDeals: postDealsLatency,
		PutDeals:  putDealsLatency,
	}

	jsonData, err := json.Marshal(requestMetrics)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %s", err.Error())
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
