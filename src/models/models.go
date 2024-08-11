package models

// POST DEAL: /v1/deals
type PostDeal struct {
	Title             string  `json:"title"`
	Value             string  `json:"value,omitempty"`
	Label             []int   `json:"label,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	UserID            int     `json:"user_id,omitempty"`
	PersonID          int     `json:"person_id,omitempty"`
	OrgID             int     `json:"org_id,omitempty"`
	PipelineID        int     `json:"pipeline_id,omitempty"`
	StageID           int     `json:"stage_id,omitempty"`
	Status            string  `json:"status,omitempty"`
	OriginID          string  `json:"origin_id,omitempty"`
	Channel           int     `json:"channel,omitempty"`
	ChannelID         string  `json:"channel_id,omitempty"`
	AddTime           string  `json:"add_time,omitempty"`
	WonTime           string  `json:"won_time,omitempty"`
	LostTime          string  `json:"lost_time,omitempty"`
	CloseTime         string  `json:"close_time,omitempty"`
	ExpectedCloseDate string  `json:"expected_close_date,omitempty"`
	Probability       float64 `json:"probability,omitempty"`
	LostReason        string  `json:"lost_reason,omitempty"`
	VisibleTo         int     `json:"visible_to,omitempty"`
}

type RequestMetrics struct {
	GetDeals  string `json:"get_deals"`
	PostDeals string `json:"post_deals"`
	PutDeals  string `json:"put_deals"`
}

// PATCH DEAL: /api/v2/deals/{id}
type PatchDeal struct {
	Title             string  `json:"title,omitempty"`
	OwnerID           int     `json:"owner_id,omitempty"`
	PersonID          int     `json:"person_id,omitempty"`
	OrgID             int     `json:"org_id,omitempty"`
	PipelineID        int     `json:"pipeline_id,omitempty"`
	StageID           int     `json:"stage_id,omitempty"`
	Value             float64 `json:"value,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	AddTime           string  `json:"add_time,omitempty"`
	UpdateTime        string  `json:"update_time,omitempty"`
	StageChangeTime   string  `json:"stage_change_time,omitempty"`
	IsDeleted         bool    `json:"is_deleted,omitempty"`
	Status            string  `json:"status,omitempty"`
	Probability       float64 `json:"probability,omitempty"`
	LostReason        string  `json:"lost_reason,omitempty"`
	VisibleTo         int     `json:"visible_to,omitempty"`
	CloseTime         string  `json:"close_time,omitempty"`
	WonTime           string  `json:"won_time,omitempty"`
	LostTime          string  `json:"lost_time,omitempty"`
	ExpectedCloseDate string  `json:"expected_close_date,omitempty"`
	LabelIDs          []int   `json:"label_ids,omitempty"`
}
