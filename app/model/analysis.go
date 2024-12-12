package model

import (
	"math"

	"github.com/marcioaso/consult/config"
	"github.com/marcioaso/consult/pkg"
	"github.com/marcioaso/consult/utils"
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
		close := item.KLine.Close
		if close == math.Trunc(close) {
			close += +0.00000000001
		}
		allCloses = append(allCloses, close)
	}
	previousSMA := previousItem.SMAS
	previousEMA := previousItem.EMAS
	previousRSI := previousItem.RSI

	sma1 := utils.CalculateLastSMA(allCloses, config.SMAS[0])
	sma2 := utils.CalculateLastSMA(allCloses, config.SMAS[1])
	sma3 := utils.CalculateLastSMA(allCloses, config.SMAS[2])

	ema1 := utils.CalculateLastEMA(allCloses, config.EMAS[0])
	ema2 := utils.CalculateLastEMA(allCloses, config.EMAS[1])

	rsi := utils.CalculateLastRSI(allCloses, config.RSI)

	item.SMAS = AverageData{
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
	item.EMAS = AverageData{
		FAST: AverageItem{
			Value:         ema1,
			Angle:         pkg.GetAngle(0, previousEMA.FAST.Value, item.Tick, ema1),
			PreviousAngle: previousEMA.FAST.Angle,
		},
		SLOW: AverageItem{
			Value:         ema2,
			Angle:         pkg.GetAngle(0, previousEMA.SLOW.Value, item.Tick, ema2),
			PreviousAngle: previousEMA.SLOW.Angle,
		},
	}
	item.RSI = AverageData{
		FAST: AverageItem{
			Value:         rsi,
			Angle:         pkg.GetAngle(0, previousRSI.FAST.Value, item.Tick, rsi),
			PreviousAngle: previousRSI.FAST.Angle,
		},
	}
}
