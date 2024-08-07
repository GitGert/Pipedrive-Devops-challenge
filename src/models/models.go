package models

// TOOL USED:
// https://transform.tools/json-to-go
type DealData struct {
	Title             string      `json:"title"`
	Value             float64     `json:"value"`
	Currency          string      `json:"currency"`
	UserID            interface{} `json:"user_id"`
	PersonID          interface{} `json:"person_id"`
	OrgID             int         `json:"org_id"`
	StageID           int         `json:"stage_id"`
	Status            string      `json:"status"`
	ExpectedCloseDate string      `json:"expected_close_date"`
	Probability       int         `json:"probability"`
	LostReason        interface{} `json:"lost_reason"`
	VisibleTo         int         `json:"visible_to"`
	AddTime           string      `json:"add_time"`
}

type RequestMetrics struct {
	GetDeals  string `json:"get_deals"`
	PostDeals string `json:"post_deals"`
	PutDeals  string `json:"put_deals"`
}

type PatchRequstData struct {
	Title    string  `json:"title"`
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}
