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

	if currentKLine.KLine.Close < currentKLine.StopLoss {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "stop loss"
		recommendation.Close = currentKLine.KLine.Close
		return *recommendation
	}

	calculateRecommendation(tailKLineData, recommendation)

	return *recommendation
}

func calculateRecommendation(tailKLineData []model.KLineAnalysisData, recommendation *model.ActionRecommendation) model.ActionRecommendation {
	currentKLine := tailKLineData[len(tailKLineData)-1]
	fiveCandles := tailKLineData[len(tailKLineData)-5]
	recommendation.Type = "wait"
	recommendation.Close = currentKLine.KLine.Close

	certainties := 0
	notes := []string{}

	// selling conditions
	if currentKLine.SMAS.HEAVY.Angle < 5 {
		certainties += 2
		notes = append(notes, "SMAS.HEAVY is down")
	}
	if currentKLine.SMAS.FAST.Angle < 5 {
		certainties += 2
		notes = append(notes, "SMAS.SMAS.Angle < 30")
	}
	if currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value {
		certainties += 10
		notes = append(notes, "currentKLine.SMAS.FAST.Value above currentKLine.SMAS.SLOW.Value")
	}

	if currentKLine.SMAS.HEAVY.Value < fiveCandles.SMAS.HEAVY.Angle {
		certainties += 3
		notes = append(notes, "currentKLine.SMAS.HEAVY.Value < fiveCandles.SMAS.HEAVY.Angle")
	}

	if certainties > 5 {
		recommendation.Type = "sell"
		recommendation.Certainty = float64(certainties)
		recommendation.Note = strings.Join(notes, ", ")
		return *recommendation
	}
	certainties = 0

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
		currentKLine.EMAS.SLOW.Angle > 0 {
		certainties += 3
		notes = append(notes, "EMAS are positive and growing")
	}
	if currentKLine.InsideBar {
		certainties += 2
		notes = append(notes, "Is an inside bar")
	}

	if certainties > 3 {
		recommendation.Type = "buy"
		recommendation.Certainty = float64(certainties)
		recommendation.Note = strings.Join(notes, ", ")
	}

	return *recommendation
}
