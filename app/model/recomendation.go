package model

type ActionRecommendation struct {
	ItemId      int64   `json:"item_id,omitempty"`
	Datetime    string  `json:"datetime"`
	Type        string  `json:"type"`
	Certainty   float64 `json:"certainty"`
	Close       float64 `json:"close"`
	BoughtPrice float64 `json:"bought_price,omitempty"`
	ProfitLoss  float64 `json:"profit_loss"`
	Note        string  `json:"note"`
	Timestamp   int64   `json:"timestamp"`
	Multiplier  float64 `json:"multiplier"`
	SellAt      float64 `json:"sell_at"`
	Liquidated  bool    `json:"liquidated"`
	Qty         float64 `json:"qty"`
}

type Directions struct {
	Recommendation ActionRecommendation `json:"recommendation"`
}
