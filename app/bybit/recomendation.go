package bybit

import "fmt"

var lastRecomendedKline KLineData
var currentRecomendation ActionRecomendation

func generateRecomendation(recomendation *ActionRecomendation, tailKLineData []KLineData) ActionRecomendation {
	recomendation.Note = ""

	currentKLine := tailKLineData[len(tailKLineData)-1]
	firstKLine := tailKLineData[0]

	if firstKLine.Open < currentKLine.Close {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "seems to be down"
		return *recomendation
	}

	if currentKLine.Directions.Heavy == "down" &&
		currentKLine.Directions.Slow == "down" {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "Slow Heavy down"
		return *recomendation
	}
	if currentKLine.SMAS.SLOW.Angle < 0 {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "currentKLine.SMAS.SLOW.Angle"
		return *recomendation
	}
	if currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value"
		return *recomendation
	}

	if currentRecomendation.Close < lastRecomendedKline.Low {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "currentRecomendation.Close < lastRecomendedKline.Low"
		return *recomendation
	}

	if currentKLine.SMAS.HEAVY.Value < firstKLine.SMAS.HEAVY.Angle {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "< firstKLine.SMAS.HEAVY.Angle"
		return *recomendation
	}

	for i, cursor := range tailKLineData {
		recomendation.Candles = i + 1
		if cursor.Directions.Fast == "up" &&
			cursor.Directions.Heavy == "up" &&
			cursor.Directions.Slow == "up" {
			calculateCertainty("up", tailKLineData, recomendation)
		}
	}
	return *recomendation
}

func calculateCertainty(movement string, klineData []KLineData, recomendation *ActionRecomendation) {
	if len(klineData) < 5 {
		recomendation.Certainty = 1.0
		recomendation.Note = ""
		return
	}

	candles := len(klineData)

	recomendation.Candles = candles

	totals := 0.0
	sums := 0.0
	weights := []float64{5, 6, 4}

	for j := 0; j < candles; j++ {
		cursor := klineData[j]
		totals += float64(j+1) * (weights[0] + weights[1] + weights[2])

		sum := 0.0
		if cursor.Directions.Fast == movement {
			sum += weights[0]
		}
		if cursor.Directions.Slow == movement {
			sum += weights[1]
		}
		if cursor.Directions.Heavy == movement {
			sum += weights[2]
		}
		sums += float64(j+1) * sum
	}
	certainty := sums / totals
	recomendation.Certainty = certainty * 100
	if movement == "up" && certainty > 0.6 {
		recomendation.Type = "buy"
		recomendation.Note = fmt.Sprintf("***** %v %v", sums, certainty)
	} else {
		recomendation.Type = "sell"
		recomendation.Note = fmt.Sprintf("***** %v %v", sums, certainty)
	}
}
