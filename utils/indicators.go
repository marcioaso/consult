package utils

import (
	"math"

	trade_indicators "github.com/go-whale/trade-indicators"
	"github.com/marcioaso/consult/app/model"
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

func CalculateLastStochastic(data []model.KLineAnalysisData, period int) float64 {
	highs := make([]float64, 0)
	lows := make([]float64, 0)
	closes := make([]float64, 0)
	for _, aData := range data {
		kline := aData.KLine
		high := kline.High
		if high == math.Trunc(high) {
			high += +0.00000001
		}
		low := kline.Low
		if low == math.Trunc(low) {
			low += +0.00000001
		}
		close := kline.Close
		if close == math.Trunc(close) {
			close += +0.00000001
		}
		highs = append(highs, high)
		lows = append(lows, low)
		closes = append(closes, close)
	}
	k, _, _ := trade_indicators.CalculateStochastic(highs, lows, closes, period)
	return k[len(k)-1]
}

func CalculateBreakoutProbabilities(data []model.KLineAnalysisData, perc float64, nbr int) []model.BreakoutLevel {
	candles := make([]model.KLineData, 0)
	for _, candle := range data {
		candles = append(candles, candle.KLine)
	}
	totalGreen := 0
	totalRed := 0
	breakoutScores := make([][4]model.BreakoutItem, nbr) // [greenHigh, greenLow, redHigh, redLow]

	for i := 1; i < len(candles); i++ {
		step := candles[i-1].Close * (perc / 100)

		green := candles[i].Close > candles[i].Open
		if green {
			totalGreen++
		} else {
			totalRed++
		}

		for j := 0; j < nbr; j++ {
			highLevel := candles[i-1].Close + float64(j+1)*step
			lowLevel := candles[i-1].Close - float64(j+1)*step

			if green && candles[i].High >= highLevel {
				breakoutScores[j][0].Probability++ // Green High
				breakoutScores[j][0].Price = highLevel
			}
			if green && candles[i].Low <= lowLevel {
				breakoutScores[j][1].Probability++ // Green Low
				breakoutScores[j][1].Price = lowLevel
			}
			if !green && candles[i].High >= highLevel {
				breakoutScores[j][2].Probability++ // Red High
				breakoutScores[j][2].Price = highLevel
			}
			if !green && candles[i].Low <= lowLevel {
				breakoutScores[j][3].Probability++ // Red Low
				breakoutScores[j][3].Price = lowLevel
			}
		}
	}
	breakouts := make([]model.BreakoutLevel, 0)
	// Print probabilities safely
	for j, scores := range breakoutScores {
		breakout := model.BreakoutLevel{
			Level: j + 1,
		}
		if totalGreen > 0 {
			breakout.GreenHigh = model.BreakoutItem{
				Probability: (scores[0].Probability / float64(totalGreen)) * 100,
				Price:       scores[0].Price,
			}
			breakout.GreenLow = model.BreakoutItem{
				Probability: (scores[1].Probability / float64(totalGreen) * 100),
				Price:       scores[1].Price,
			}
		}
		if totalRed > 0 {
			breakout.RedHigh = model.BreakoutItem{
				Probability: (scores[2].Probability / float64(totalRed)) * 100,
				Price:       scores[2].Price,
			}
			breakout.RedLow = model.BreakoutItem{
				Probability: (scores[3].Probability / float64(totalRed)) * 100,
				Price:       scores[3].Price,
			}
		}
		breakouts = append(breakouts, breakout)
	}
	return breakouts
}
