package bybit

import (
	"github.com/marcioaso/consult/app/model"
)

var LastRecommendedKline = model.KLineAnalysisData{}
var PreviousRecommendation = model.ActionRecommendation{
	Multiplier: 1,
}

func GenerateRecommendations(tailKLineData []model.KLineAnalysisData) []model.ActionRecommendation {
	changeRecommendations := make([]model.ActionRecommendation, 0)
	allTimeHigh := 0.0

	for _, each := range tailKLineData {
		if each.KLine.Close > allTimeHigh {
			allTimeHigh = each.KLine.Close
		}
	}

	currentKLine := tailKLineData[len(tailKLineData)-1]
	previousKLine := tailKLineData[len(tailKLineData)-2]

	_, sell := Position.GetPositions(currentKLine.KLine.Close)
	if len(sell) > 0 {
		for _, sold := range sell {
			kline := currentKLine.KLine
			s := model.ActionRecommendation{
				Type:        "sell",
				Close:       kline.Close,
				BoughtPrice: sold.Close,
				Qty:         sold.Qty,
				Note:        "peaked all time high",
				ProfitLoss:  (sold.Qty * kline.Close) - (sold.Qty * sold.Close),
				SellAt:      sold.SellAt,
				ItemId:      kline.Id,
				Timestamp:   kline.Timestamp,
				Datetime:    kline.Datetime,
			}

			changeRecommendations = append(changeRecommendations, s)
		}
	}

	previousIsLower := previousKLine.SMAS.SLOW.Value < previousKLine.EMAS.FAST.Value
	currentIsHigher := currentKLine.SMAS.SLOW.Value > currentKLine.EMAS.FAST.Value

	if previousIsLower && currentIsHigher &&
		currentKLine.SMAS.SLOW.Angle > 0 {
		recommendation := model.ActionRecommendation{
			Type:      "buy",
			Note:      "lowest price",
			Close:     currentKLine.KLine.Close,
			Timestamp: currentKLine.KLine.Timestamp,
			ItemId:    currentKLine.KLine.Id,
			SellAt:    allTimeHigh,
		}
		if recommendation.Multiplier == 0 {
			recommendation.Multiplier = 1
		}
		changeRecommendations = append(changeRecommendations, recommendation)
	}

	return changeRecommendations
}
