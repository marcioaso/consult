package model

type BreakoutItem struct {
	Probability float64 `json:"probability"`
	Price       float64 `json:"price"`
}
type BreakoutLevel struct {
	Level     int          `json:"level"`
	GreenHigh BreakoutItem `json:"green_high"`
	GreenLow  BreakoutItem `json:"green_low"`
	RedHigh   BreakoutItem `json:"red_high"`
	RedLow    BreakoutItem `json:"red_low"`
}
type KLineAnalysisData struct {
	KLine KLineData `json:"kline"`
	Tick  float64   `json:"tick"`

	StopLoss float64 `json:"stop_loss"`

	SMAS      AverageData     `json:"sma"`
	EMAS      AverageData     `json:"ema"`
	RSI       AverageData     `json:"rsi"`
	Breakouts []BreakoutLevel `json:"breakouts"`

	InsideBar bool `json:"inside_bar"`

	History []KLineAnalysisData `json:"-"`
}

type KLinePotential struct {
	Initial    float64 `json:"initial"`
	Final      float64 `json:"final"`
	Bought     float64 `json:"bought"`
	Price      float64 `json:"-"`
	Variation  float64 `json:"variation"`
	Percentage float64 `json:"percentage"`
}

type KLineBacktestResponse struct {
	Potential       KLinePotential         `json:"potential,omitempty"`
	Recommendations []ActionRecommendation `json:"recommendations"`
}

type KLineAnalysisResponse struct {
	Recommendations []ActionRecommendation `json:"recommendations"`
	Data            []KLineAnalysisData    `json:"data"`
}

func (r KLineData) ToAnalysis(tick float64) KLineAnalysisData {
	n := KLineAnalysisData{
		KLine: r,
		Tick:  tick,
	}

	return n
}
