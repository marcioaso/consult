package bybit

import "github.com/marcioaso/consult/app/model"

var lastRecomendedKline model.KLineData
var currentRecomendation model.ActionRecomendation

func generateRecomendation(recomendation *model.ActionRecomendation, tailKLineData []model.KLineAnalysisData) model.ActionRecomendation {
	recomendation.Note = ""

	currentKLine := tailKLineData[len(tailKLineData)-1]
	firstKLine := tailKLineData[0]

	if firstKLine.KLine.Open < currentKLine.KLine.Close {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "seems to go down"
		return *recomendation
	}

	if currentRecomendation.Close < lastRecomendedKline.Low {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "currentRecomendation.Close < lastRecomendedKline.Low"
		return *recomendation
	}

	calculateBySMAS(currentKLine, firstKLine, recomendation)

	return *recomendation
}

func calculateBySMAS(currentKLine, firstKLine model.KLineAnalysisData, recomendation *model.ActionRecomendation) model.ActionRecomendation {
	currentSmas := currentKLine.SMAS

	if currentKLine.SMAS.HEAVY.Direction == "down" {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "Heavy down"
		return *recomendation
	}
	if currentKLine.SMAS.FAST.Angle < 0 {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "currentKLine.SMAS.FAST.Angle"
		return *recomendation
	}
	if currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "currentKLine.SMAS.FAST.Value < currentKLine.SMAS.SLOW.Value"
		return *recomendation
	}

	if currentKLine.SMAS.HEAVY.Value < firstKLine.SMAS.HEAVY.Angle {
		recomendation.Type = "sell"
		recomendation.Certainty = 100
		recomendation.Note = "< firstKLine.SMAS.HEAVY.Angle"
		return *recomendation
	}

	if currentSmas.FAST.Angle > 0 &&
		currentSmas.SLOW.Angle > 0 &&
		currentSmas.HEAVY.Value > 0 &&
		currentSmas.FAST.Angle < 45 && currentSmas.FAST.Angle > 3 {
		recomendation.Type = "buy"
		recomendation.Note = "tendency is up"

	}

	return *recomendation

}
