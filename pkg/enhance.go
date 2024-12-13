package pkg

import (
	"math"

	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/config"
	"github.com/marcioaso/consult/utils"
)

func EnhanceAverageData(item *model.KLineAnalysisData, previousItem model.KLineAnalysisData) {
	allCloses := []float64{}
	for _, item := range item.History {
		close := item.KLine.Close
		if close == math.Trunc(close) {
			close += +0.00000000001
		}
		allCloses = append(allCloses, close)
	}

	item.Stochastic = utils.CalculateLastStochastic(item.History, config.Stochastic)

	previousSMA := previousItem.SMAS
	previousEMA := previousItem.EMAS
	previousRSI := previousItem.RSI

	sma1 := utils.CalculateLastSMA(allCloses, config.SMAS[0])
	sma2 := utils.CalculateLastSMA(allCloses, config.SMAS[1])
	sma3 := utils.CalculateLastSMA(allCloses, config.SMAS[2])

	ema1 := 0.0
	if len(item.History) > config.EMAS[0] {
		ema1 = utils.CalculateLastEMA(allCloses, config.EMAS[0])
	}
	ema2 := 0.0
	if len(item.History) > config.EMAS[1] {
		ema2 = utils.CalculateLastEMA(allCloses, config.EMAS[1])
	}
	ema3 := 0.0
	if len(item.History) > config.EMAS[2] {
		ema3 = utils.CalculateLastEMA(allCloses, config.EMAS[2])
	}

	rsi := utils.CalculateLastRSI(allCloses, config.RSI)

	item.SMAS = model.AverageData{
		FAST: model.AverageItem{
			Value:         sma1,
			Angle:         GetAngle(0, previousSMA.FAST.Value, item.Tick, sma1),
			PreviousAngle: previousSMA.FAST.Angle,
		},
		SLOW: model.AverageItem{
			Value:         sma2,
			Angle:         GetAngle(0, previousSMA.SLOW.Value, item.Tick, sma2),
			PreviousAngle: previousSMA.SLOW.Angle,
		},
		HEAVY: model.AverageItem{
			Value:         sma3,
			Angle:         GetAngle(0, previousSMA.HEAVY.Value, item.Tick, sma3),
			PreviousAngle: previousSMA.HEAVY.Angle,
		},
	}
	item.EMAS = model.AverageData{
		FAST: model.AverageItem{
			Value:         ema1,
			Angle:         GetAngle(0, previousEMA.FAST.Value, item.Tick, ema1),
			PreviousAngle: previousEMA.FAST.Angle,
		},
		SLOW: model.AverageItem{
			Value:         ema2,
			Angle:         GetAngle(0, previousEMA.SLOW.Value, item.Tick, ema2),
			PreviousAngle: previousEMA.SLOW.Angle,
		},
		HEAVY: model.AverageItem{
			Value:         ema3,
			Angle:         GetAngle(0, previousEMA.HEAVY.Value, item.Tick, ema3),
			PreviousAngle: previousEMA.SLOW.Angle,
		},
	}
	item.RSI = model.AverageData{
		FAST: model.AverageItem{
			Value:         rsi,
			Angle:         GetAngle(0, previousRSI.FAST.Value, item.Tick, rsi),
			PreviousAngle: previousRSI.FAST.Angle,
		},
	}

	breakouts := utils.CalculateBreakoutProbabilities(item.History, config.BREAKOUT_PERCENT, config.BREAKOUT_LAYERS)
	item.Breakouts = breakouts
}
