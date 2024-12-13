package model

type ActionRecommendation struct {
	Datetime   string  `json:"datetime"`
	Type       string  `json:"type"`
	Certainty  float64 `json:"certainty"`
	Close      float64 `json:"close"`
	Risk       float64 `json:"risk"`
	StopLoss   float64 `json:"stop_loss"`
	ProfitLoss float64 `json:"profit_loss"`
	Note       string  `json:"note"`
}

type Directions struct {
	Recommendation ActionRecommendation `json:"recommendation"`
}
