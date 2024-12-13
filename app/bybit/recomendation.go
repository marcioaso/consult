package bybit

import (
	"strings"

	"github.com/marcioaso/consult/app/model"
)

var LastRecommendedKline = model.KLineAnalysisData{}
var CurrentRecommendation = model.ActionRecommendation{}

func GenerateRecommendation(recommendation *model.ActionRecommendation, tailKLineData []model.KLineAnalysisData) model.ActionRecommendation {
	recommendation.Note = ""
	currentKLine := tailKLineData[len(tailKLineData)-1]

	if currentKLine.KLine.Close < CurrentRecommendation.StopLoss {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "stop loss"
		recommendation.Close = currentKLine.KLine.Close
		return *recommendation
	}

	if CurrentRecommendation.Type == "buy" && currentKLine.KLine.Close > CurrentRecommendation.Risk {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "by risk"
		recommendation.Close = currentKLine.KLine.Close
		return *recommendation
	}
	calculateRecommendation(tailKLineData, recommendation)
	return *recommendation
}

var deviation = 0.0

func calculateRecommendation(tailKLineData []model.KLineAnalysisData, recommendation *model.ActionRecommendation) model.ActionRecommendation {
	deviation = tailKLineData[0].KLine.Close * 0.001 // 0.1% range
	currentKLine := tailKLineData[len(tailKLineData)-1]
	fiveCandles := tailKLineData[len(tailKLineData)-5]
	recommendation.Type = "wait"
	recommendation.Close = currentKLine.KLine.Close

	certainties := 0
	notes := []string{}

	// // selling conditions
	// if currentKLine.EMAS.HEAVY.Angle < 2 {
	// 	certainties += 5
	// 	notes = append(notes, "EMAS.HEAVY is down")
	// }
	// if currentKLine.SMAS.FAST.Value < currentKLine.SMAS.HEAVY.Value {
	// 	certainties += 10
	// 	notes = append(notes, "currentKLine.SMAS.FAST.Value above currentKLine.HEAVY.SLOW.Value")
	// }

	// if currentKLine.SMAS.HEAVY.Value < fiveCandles.SMAS.HEAVY.Angle {
	// 	certainties += 3
	// 	notes = append(notes, "currentKLine.SMAS.HEAVY.Value < fiveCandles.SMAS.HEAVY.Angle")
	// }

	// if certainties > 5 {
	// 	recommendation.Type = "sell"
	// 	recommendation.Certainty = float64(certainties)
	// 	recommendation.Note = strings.Join(notes, ", ")
	// 	return *recommendation
	// }
	// certainties = 0
	// notes = []string{}

	// buying conditions

	if currentKLine.RSI.FAST.Value < 31 {
		notes = append(notes, "currentKLine.RSI.FAST.Value < 30")
		certainties += 4
	}
	if currentKLine.SMAS.FAST.Angle > 0 &&
		currentKLine.SMAS.SLOW.Angle > 0 &&
		currentKLine.SMAS.HEAVY.Angle > 0 {
		certainties += 2
		notes = append(notes, "all SMAS increased")
	}
	if currentKLine.EMAS.FAST.Angle > 0 &&
		currentKLine.EMAS.SLOW.Angle > 0 &&
		currentKLine.EMAS.FAST.Angle > currentKLine.EMAS.SLOW.Angle {
		certainties += 3
		notes = append(notes, "EMAS are positive")
	}
	if currentKLine.EMAS.HEAVY.Angle > 15 && currentKLine.EMAS.HEAVY.Angle > currentKLine.EMAS.HEAVY.PreviousAngle {
		certainties += 5
		notes = append(notes, "EMAS heavy > 15")
	}
	if currentKLine.InsideBar {
		certainties += 2
		notes = append(notes, "Is an inside bar")
	}
	if withinRange(currentKLine.EMAS.FAST.Value, currentKLine.EMAS.SLOW.Value) &&
		withinRange(currentKLine.EMAS.FAST.Value, currentKLine.SMAS.FAST.Value) &&
		withinRange(currentKLine.SMAS.FAST.Value, currentKLine.SMAS.SLOW.Value) &&
		withinRange(currentKLine.SMAS.SLOW.Value, currentKLine.SMAS.HEAVY.Value) &&
		currentKLine.SMAS.SLOW.Value > fiveCandles.SMAS.SLOW.Value {
		certainties += 7
		notes = append(notes, "All aligned and are growing")
	}

	if certainties > 9 {
		recommendation.Type = "buy"
		recommendation.Certainty = float64(certainties)
		recommendation.Note = strings.Join(notes, ", ")
	}

	return *recommendation
}

func withinRange(input1, input2 float64) bool {
	return (input1-deviation <= input2 && input2 <= input1+deviation)
}
