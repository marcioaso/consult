package model

type KLineAnalysisData struct {
	KLine KLineData `json:"kline"`
	Tick  float64   `json:"tick"`

	Recomendations []ActionRecomendation `json:"recomendations"`

	SMAS AverageData `json:"sma"`
	EMAS AverageData `json:"ema"`
	RSI  AverageData `json:"rsi"`

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
