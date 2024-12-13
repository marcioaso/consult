package bybit

import "github.com/marcioaso/consult/app/model"

var lastRecomendedKline model.KLineData
var currentRecommendation model.ActionRecommendation

func generateRecommendation(recommendation *model.ActionRecommendation, tailKLineData []model.KLineAnalysisData) model.ActionRecommendation {
	recommendation.Note = ""

	currentKLine := tailKLineData[len(tailKLineData)-1]
	firstKLine := tailKLineData[0]

	if firstKLine.KLine.Open < currentKLine.KLine.Close {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "seems to go down"
		return *recommendation
	}

	if currentRecommendation.Close < lastRecomendedKline.Low {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "currentRecommendation.Close < lastRecomendedKline.Low"
		return *recommendation
	}

	calculateBySMAS(currentKLine, firstKLine, recommendation)

	return *recommendation
}

func calculateBySMAS(currentKLine, firstKLine model.KLineAnalysisData, recommendation *model.ActionRecommendation) model.ActionRecommendation {
	currentSmas := currentKLine.SMAS

	if currentKLine.SMAS.HEAVY.Direction == "down" {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "Heavy down"
		return *recommendation
	}
	if currentKLine.SMAS.FAST.Angle < 0 {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "currentKLine.SMAS.FAST.Angle"
		return *recommendation
	}
	if currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value"
		return *recommendation
	}

	if currentKLine.SMAS.HEAVY.Value < firstKLine.SMAS.HEAVY.Angle {
		recommendation.Type = "sell"
		recommendation.Certainty = 100
		recommendation.Note = "< firstKLine.SMAS.HEAVY.Angle"
		return *recommendation
	}

	if currentSmas.FAST.Angle > 0 &&
		currentSmas.SLOW.Angle > 0 &&
		currentSmas.HEAVY.Value > 0 &&
		currentSmas.FAST.Angle < 45 && currentSmas.FAST.Angle > 3 {
		recommendation.Type = "buy"
		recommendation.Note = "tendency is up"

	}

	return *recommendation

}
