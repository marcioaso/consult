package model

import (
	"github.com/marcioaso/consult/config"
	"github.com/marcioaso/consult/pkg"
)

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

func (a *KLineAnalysisData) Analyze(previousItem KLineAnalysisData) {
	enhanceAverageData(a, previousItem)
}
func (r KLineData) ToAnalysis(tick float64) KLineAnalysisData {
	n := KLineAnalysisData{
		KLine: r,
		Tick:  tick,
	}

	return n
}

func enhanceAverageData(item *KLineAnalysisData, previousItem KLineAnalysisData) {
	allCloses := []float64{}
	for _, item := range item.History {
		allCloses = append(allCloses, item.KLine.Close)
	}

	sma1, _ := pkg.CalculateLastSMA(allCloses, config.SmaConf[0])
	sma2, _ := pkg.CalculateLastSMA(allCloses, config.SmaConf[1])
	sma3, _ := pkg.CalculateLastSMA(allCloses, config.SmaConf[2])

	previousSMA := previousItem.SMAS

	smas := AverageData{
		FAST: AverageItem{
			Value:         sma1,
			Angle:         pkg.GetAngle(0, previousSMA.FAST.Value, item.Tick, sma1),
			PreviousAngle: previousSMA.FAST.Angle,
		},
		SLOW: AverageItem{
			Value:         sma2,
			Angle:         pkg.GetAngle(0, previousSMA.SLOW.Value, item.Tick, sma2),
			PreviousAngle: previousSMA.SLOW.Angle,
		},
		HEAVY: AverageItem{
			Value:         sma3,
			Angle:         pkg.GetAngle(0, previousSMA.HEAVY.Value, item.Tick, sma3),
			PreviousAngle: previousSMA.HEAVY.Angle,
		},
	}
	item.SMAS = smas
}
