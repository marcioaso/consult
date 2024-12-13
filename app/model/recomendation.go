package model

type ActionRecommendation struct {
	Datetime  string  `json:"datetime"`
	Type      string  `json:"type"`
	Candles   int     `json:"candles"`
	Certainty float64 `json:"certainty"`
	Close     float64 `json:"close"`
	Note      string  `json:"note"`
}

type Directions struct {
	Recommendation ActionRecommendation `json:"recommendation"`
}
