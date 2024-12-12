package utils

import (
	"fmt"

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

func calculateBreakoutProbabilities(data []model.KLineAnalysisData, perc float64, nbr int) {
	candles := make([]model.KLineData, 0)
	for _, candle := range data {
		candles = append(candles, candle.KLine)
	}
	totalGreen := 0
	totalRed := 0
	type LevelData struct {
		Probability float64
		Price       float64
	}
	breakoutScores := make([][4]LevelData, nbr) // [greenHigh, greenLow, redHigh, redLow]

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

	// Print probabilities safely
	for j, scores := range breakoutScores {
		fmt.Printf("Level %d:\n", j+1)
		if totalGreen > 0 {
			fmt.Printf("  Green High: %.2f%% (Price: %.2f)\n", (scores[0].Probability/float64(totalGreen))*100, scores[0].Price)
			fmt.Printf("  Green Low: %.2f%% (Price: %.2f)\n", (scores[1].Probability/float64(totalGreen))*100, scores[1].Price)
		} else {
			fmt.Println("  Green High: N/A")
			fmt.Println("  Green Low: N/A")
		}

		if totalRed > 0 {
			fmt.Printf("  Red High: %.2f%% (Price: %.2f)\n", (scores[2].Probability/float64(totalRed))*100, scores[2].Price)
			fmt.Printf("  Red Low: %.2f%% (Price: %.2f)\n", (scores[3].Probability/float64(totalRed))*100, scores[3].Price)
		} else {
			fmt.Println("  Red High: N/A")
			fmt.Println("  Red Low: N/A")
		}
	}
}
