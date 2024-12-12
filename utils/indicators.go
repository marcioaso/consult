package utils

import (
	trade_indicators "github.com/go-whale/trade-indicators"
)

func CalculateLastSMA(data []float64, period int) float64 {
	smas, _ := trade_indicators.CalculateSMA(data, period)
	return smas[len(smas)-1]
}

func CalculateLastEMA(data []float64, period int) float64 {
	mme, _ := trade_indicators.CalculateEMA(data, period)
	return mme[len(mme)-1]
}

func CalculateLastRSI(prices []float64, period int) float64 {
	rsis, _ := trade_indicators.CalculateRSI(prices, period)
	return rsis[len(rsis)-1]
}
