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

	Recommendations []ActionRecommendation `json:"recommendations"`

	SMAS      AverageData     `json:"sma"`
	EMAS      AverageData     `json:"ema"`
	RSI       AverageData     `json:"rsi"`
	Breakouts []BreakoutLevel `json:"breakouts"`

	InsideBar bool `json:"inside_bar"`

	History []KLineAnalysisData `json:"-"`
}

func (r KLineData) ToAnalysis(tick float64) KLineAnalysisData {
	n := KLineAnalysisData{
		KLine: r,
		Tick:  tick,
	}

	return n
}
