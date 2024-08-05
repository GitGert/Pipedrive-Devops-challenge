package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	// Define routes
	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/post_deals", postDeals).Methods("GET") //TODO: CHANGE TO POST
	r.HandleFunc("/put_deals", putDeals).Methods("GET")   //TODO: CHANGE TO PUT

	// Start the server
	fmt.Println("server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// // Pipedrive company domain
// $company_domain = 'YOUR_COMPANY_DOMAIN';
// https://pipedrive.com/api/v1/deals?limit=10
// URL for Deal listing
// $url = 'https://'+ COMPANY_DOMAIN+'.pipedrive.com/api/v1/deals?limit=10&api_token=' + API_TOKEN;

// TODO: consider adding this kind of modularity into your code:       const response = await api.addDeal(data);
func getDeals(w http.ResponseWriter, r *http.Request) {
	var API_TOKEN = os.Getenv("API_TOKEN")
	var COMPANY_DOMAIN = os.Getenv("COMPANY_DOMAIN")
	w.Header().Set("Content-Type", "application/json")

	// json.NewEncoder(w).Encode(todos)

	// requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	// requestURL := "https://developers.pipedrive.com/docs/api/v1/Deals)"
	requestURL := "https://" + COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals?limit=10&api_token=" + API_TOKEN
	// requestURL := "https://pipedrive.com/api/v1/deals?limit=10"
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))

	var apiResponse ApiResponse

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	// fmt.Println("\n\n\n\n")

	fmt.Println(apiResponse.Data)

	// value := json.Unmarshal(res)
	// fmt.Println(res)
	// fmt.Printf("client: status code: %d\n", res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Use the appropriate status code

	json.NewEncoder(w).Encode(apiResponse.Data)
}

func postDeals(w http.ResponseWriter, r *http.Request) {

	var API_TOKEN = os.Getenv("API_TOKEN")
	var COMPANY_DOMAIN = os.Getenv("COMPANY_DOMAIN")

	dealData := DealData{
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

	postBody, err := json.Marshal(dealData)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %s", err.Error())
	}
	responseBody := bytes.NewBuffer(postBody)

	requestURL := "https://" + COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals?api_token=" + API_TOKEN
	fmt.Println(requestURL)

	resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)

	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	fmt.Println(sb)
	if resp.StatusCode == 200 {
		// json.NewEncoder(w).Encode(resp.Body)
		w.Header().Set("Content-Type", "application/json") // Set the Content-Type header
		w.Write(body)                                      // Write the response body
	}

}

func putDeals(w http.ResponseWriter, r *http.Request) {
	var API_TOKEN = os.Getenv("API_TOKEN")
	var COMPANY_DOMAIN = os.Getenv("COMPANY_DOMAIN")

	userID_Gert := "21814964"
	dealID := "1" //TODO: CHANGE THIS
	data, err := json.Marshal(map[string]interface{}{
		"user_id": userID_Gert,
	})

	requestURL := "https://" + COMPANY_DOMAIN + ".pipedrive.com/api/v1/deals/" + dealID + "?api_token=" + API_TOKEN
	fmt.Println(requestURL)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(data))
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
	// resp, err := http.Pu("https://postman-echo.com/post", "application/json", responseBody)

	// json.NewEncoder(w).Encode(todos)

	w.Header().Set("Content-Type", "application/json")
}

// TOOL USED:
// https://transform.tools/json-to-go
type ApiResponse struct {
	Success bool `json:"success"`
	Data    []struct {
		ID            int `json:"id"`
		CreatorUserID struct {
			ID         int         `json:"id"`
			Name       string      `json:"name"`
			Email      string      `json:"email"`
			HasPic     int         `json:"has_pic"`
			PicHash    interface{} `json:"pic_hash"`
			ActiveFlag bool        `json:"active_flag"`
			Value      int         `json:"value"`
		} `json:"creator_user_id"`
		UserID struct {
			ID         int         `json:"id"`
			Name       string      `json:"name"`
			Email      string      `json:"email"`
			HasPic     int         `json:"has_pic"`
			PicHash    interface{} `json:"pic_hash"`
			ActiveFlag bool        `json:"active_flag"`
			Value      int         `json:"value"`
		} `json:"user_id"`
		PersonID struct {
			ActiveFlag bool   `json:"active_flag"`
			Name       string `json:"name"`
			Email      []struct {
				Label   string `json:"label"`
				Value   string `json:"value"`
				Primary bool   `json:"primary"`
			} `json:"email"`
			Phone []struct {
				Label   string `json:"label"`
				Value   string `json:"value"`
				Primary bool   `json:"primary"`
			} `json:"phone"`
			OwnerID int `json:"owner_id"`
			Value   int `json:"value"`
		} `json:"person_id"`
		OrgID struct {
			Name        string        `json:"name"`
			PeopleCount int           `json:"people_count"`
			OwnerID     int           `json:"owner_id"`
			Address     interface{}   `json:"address"`
			ActiveFlag  bool          `json:"active_flag"`
			CcEmail     string        `json:"cc_email"`
			LabelIds    []interface{} `json:"label_ids"`
			OwnerName   string        `json:"owner_name"`
			Value       int           `json:"value"`
		} `json:"org_id"`
		StageID                int         `json:"stage_id"`
		Title                  string      `json:"title"`
		Value                  int         `json:"value"`
		Acv                    interface{} `json:"acv"`
		Mrr                    interface{} `json:"mrr"`
		Arr                    interface{} `json:"arr"`
		Currency               string      `json:"currency"`
		AddTime                string      `json:"add_time"`
		UpdateTime             string      `json:"update_time"`
		StageChangeTime        interface{} `json:"stage_change_time"`
		Active                 bool        `json:"active"`
		Deleted                bool        `json:"deleted"`
		Status                 string      `json:"status"`
		Probability            interface{} `json:"probability"`
		NextActivityDate       string      `json:"next_activity_date"`
		NextActivityTime       interface{} `json:"next_activity_time"`
		NextActivityID         int         `json:"next_activity_id"`
		LastActivityID         int         `json:"last_activity_id"`
		LastActivityDate       string      `json:"last_activity_date"`
		LostReason             interface{} `json:"lost_reason"`
		VisibleTo              string      `json:"visible_to"`
		CloseTime              interface{} `json:"close_time"`
		PipelineID             int         `json:"pipeline_id"`
		WonTime                interface{} `json:"won_time"`
		FirstWonTime           interface{} `json:"first_won_time"`
		LostTime               interface{} `json:"lost_time"`
		ProductsCount          int         `json:"products_count"`
		FilesCount             int         `json:"files_count"`
		NotesCount             int         `json:"notes_count"`
		FollowersCount         int         `json:"followers_count"`
		EmailMessagesCount     int         `json:"email_messages_count"`
		ActivitiesCount        int         `json:"activities_count"`
		DoneActivitiesCount    int         `json:"done_activities_count"`
		UndoneActivitiesCount  int         `json:"undone_activities_count"`
		ParticipantsCount      int         `json:"participants_count"`
		ExpectedCloseDate      string      `json:"expected_close_date"`
		LastIncomingMailTime   interface{} `json:"last_incoming_mail_time"`
		LastOutgoingMailTime   interface{} `json:"last_outgoing_mail_time"`
		Label                  interface{} `json:"label"`
		LocalWonDate           interface{} `json:"local_won_date"`
		LocalLostDate          interface{} `json:"local_lost_date"`
		LocalCloseDate         interface{} `json:"local_close_date"`
		Origin                 interface{} `json:"origin"`
		OriginID               interface{} `json:"origin_id"`
		Channel                interface{} `json:"channel"`
		ChannelID              interface{} `json:"channel_id"`
		StageOrderNr           int         `json:"stage_order_nr"`
		PersonName             string      `json:"person_name"`
		OrgName                string      `json:"org_name"`
		NextActivitySubject    string      `json:"next_activity_subject"`
		NextActivityType       string      `json:"next_activity_type"`
		NextActivityDuration   interface{} `json:"next_activity_duration"`
		NextActivityNote       interface{} `json:"next_activity_note"`
		FormattedValue         string      `json:"formatted_value"`
		WeightedValue          int         `json:"weighted_value"`
		FormattedWeightedValue string      `json:"formatted_weighted_value"`
		WeightedValueCurrency  string      `json:"weighted_value_currency"`
		RottenTime             interface{} `json:"rotten_time"`
		AcvCurrency            interface{} `json:"acv_currency"`
		MrrCurrency            interface{} `json:"mrr_currency"`
		ArrCurrency            interface{} `json:"arr_currency"`
		OwnerName              string      `json:"owner_name"`
		CcEmail                string      `json:"cc_email"`
		OrgHidden              bool        `json:"org_hidden"`
		PersonHidden           bool        `json:"person_hidden"`
	} `json:"data"`
	AdditionalData struct {
		Pagination struct {
			Start                 int  `json:"start"`
			Limit                 int  `json:"limit"`
			MoreItemsInCollection bool `json:"more_items_in_collection"`
			NextStart             int  `json:"next_start"`
		} `json:"pagination"`
	} `json:"additional_data"`
	RelatedObjects struct {
		User struct {
			Num21814964 struct {
				ID         int         `json:"id"`
				Name       string      `json:"name"`
				Email      string      `json:"email"`
				HasPic     int         `json:"has_pic"`
				PicHash    interface{} `json:"pic_hash"`
				ActiveFlag bool        `json:"active_flag"`
			} `json:"21814964"`
		} `json:"user"`
		Organization struct {
			Num1 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"1"`
			Num2 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"2"`
			Num3 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"3"`
			Num4 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"4"`
			Num5 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"5"`
			Num8 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"8"`
			Num9 struct {
				ID          int           `json:"id"`
				Name        string        `json:"name"`
				PeopleCount int           `json:"people_count"`
				OwnerID     int           `json:"owner_id"`
				Address     interface{}   `json:"address"`
				ActiveFlag  bool          `json:"active_flag"`
				CcEmail     string        `json:"cc_email"`
				LabelIds    []interface{} `json:"label_ids"`
				OwnerName   string        `json:"owner_name"`
			} `json:"9"`
		} `json:"organization"`
		Pipeline struct {
			Num1 struct {
				ID              int         `json:"id"`
				Name            string      `json:"name"`
				URLTitle        string      `json:"url_title"`
				OrderNr         int         `json:"order_nr"`
				Active          bool        `json:"active"`
				DealProbability bool        `json:"deal_probability"`
				AddTime         string      `json:"add_time"`
				UpdateTime      interface{} `json:"update_time"`
			} `json:"1"`
		} `json:"pipeline"`
		Person struct {
			Num1 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"1"`
			Num2 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"2"`
			Num3 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"3"`
			Num4 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"4"`
			Num5 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"5"`
			Num6 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"6"`
			Num7 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"7"`
			Num8 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"8"`
			Num9 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"9"`
			Num10 struct {
				ActiveFlag bool   `json:"active_flag"`
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Email      []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"email"`
				Phone []struct {
					Label   string `json:"label"`
					Value   string `json:"value"`
					Primary bool   `json:"primary"`
				} `json:"phone"`
				OwnerID int `json:"owner_id"`
			} `json:"10"`
		} `json:"person"`
		Stage struct {
			Num1 struct {
				ID                      int         `json:"id"`
				OrderNr                 int         `json:"order_nr"`
				Name                    string      `json:"name"`
				ActiveFlag              bool        `json:"active_flag"`
				DealProbability         int         `json:"deal_probability"`
				PipelineID              int         `json:"pipeline_id"`
				RottenFlag              bool        `json:"rotten_flag"`
				RottenDays              interface{} `json:"rotten_days"`
				AddTime                 string      `json:"add_time"`
				UpdateTime              interface{} `json:"update_time"`
				PipelineName            string      `json:"pipeline_name"`
				PipelineDealProbability bool        `json:"pipeline_deal_probability"`
			} `json:"1"`
			Num2 struct {
				ID                      int         `json:"id"`
				OrderNr                 int         `json:"order_nr"`
				Name                    string      `json:"name"`
				ActiveFlag              bool        `json:"active_flag"`
				DealProbability         int         `json:"deal_probability"`
				PipelineID              int         `json:"pipeline_id"`
				RottenFlag              bool        `json:"rotten_flag"`
				RottenDays              interface{} `json:"rotten_days"`
				AddTime                 string      `json:"add_time"`
				UpdateTime              interface{} `json:"update_time"`
				PipelineName            string      `json:"pipeline_name"`
				PipelineDealProbability bool        `json:"pipeline_deal_probability"`
			} `json:"2"`
			Num3 struct {
				ID                      int         `json:"id"`
				OrderNr                 int         `json:"order_nr"`
				Name                    string      `json:"name"`
				ActiveFlag              bool        `json:"active_flag"`
				DealProbability         int         `json:"deal_probability"`
				PipelineID              int         `json:"pipeline_id"`
				RottenFlag              bool        `json:"rotten_flag"`
				RottenDays              interface{} `json:"rotten_days"`
				AddTime                 string      `json:"add_time"`
				UpdateTime              interface{} `json:"update_time"`
				PipelineName            string      `json:"pipeline_name"`
				PipelineDealProbability bool        `json:"pipeline_deal_probability"`
			} `json:"3"`
			Num4 struct {
				ID                      int         `json:"id"`
				OrderNr                 int         `json:"order_nr"`
				Name                    string      `json:"name"`
				ActiveFlag              bool        `json:"active_flag"`
				DealProbability         int         `json:"deal_probability"`
				PipelineID              int         `json:"pipeline_id"`
				RottenFlag              bool        `json:"rotten_flag"`
				RottenDays              interface{} `json:"rotten_days"`
				AddTime                 string      `json:"add_time"`
				UpdateTime              interface{} `json:"update_time"`
				PipelineName            string      `json:"pipeline_name"`
				PipelineDealProbability bool        `json:"pipeline_deal_probability"`
			} `json:"4"`
			Num5 struct {
				ID                      int         `json:"id"`
				OrderNr                 int         `json:"order_nr"`
				Name                    string      `json:"name"`
				ActiveFlag              bool        `json:"active_flag"`
				DealProbability         int         `json:"deal_probability"`
				PipelineID              int         `json:"pipeline_id"`
				RottenFlag              bool        `json:"rotten_flag"`
				RottenDays              interface{} `json:"rotten_days"`
				AddTime                 string      `json:"add_time"`
				UpdateTime              interface{} `json:"update_time"`
				PipelineName            string      `json:"pipeline_name"`
				PipelineDealProbability bool        `json:"pipeline_deal_probability"`
			} `json:"5"`
		} `json:"stage"`
	} `json:"related_objects"`
}

type DealData struct {
	Title             string      `json:"title"`
	Value             float64     `json:"value"`
	Currency          string      `json:"currency"`
	UserID            interface{} `json:"user_id"`   // Using interface{} to allow null values
	PersonID          interface{} `json:"person_id"` // Using interface{} to allow null values
	OrgID             int         `json:"org_id"`
	StageID           int         `json:"stage_id"`
	Status            string      `json:"status"`
	ExpectedCloseDate string      `json:"expected_close_date"`
	Probability       int         `json:"probability"`
	LostReason        interface{} `json:"lost_reason"` // Using interface{} to allow null values
	VisibleTo         int         `json:"visible_to"`
	AddTime           string      `json:"add_time"`
}
