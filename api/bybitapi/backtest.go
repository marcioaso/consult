package bybitapi

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/app/bybit"
	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/pkg"
	"github.com/marcioaso/consult/utils"
)

func BacktestHandler(c echo.Context) error {
	bybit.LastRecommendedKline = model.KLineAnalysisData{}
	bybit.CurrentRecommendation = model.ActionRecommendation{}

	defer pkg.Elapsed("backtest")()
	defaultInterval := "15m"
	defaultLimit := "150"
	defaultTo := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond)) // current time in milliseconds

	ticker := c.QueryParam("symbol")
	if ticker == "" {
		return c.String(http.StatusBadRequest, "Missing 'symbol' query parameter")
	}

	interval := c.QueryParam("interval")
	if interval == "" {
		interval = defaultInterval
	}
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = defaultLimit
	}
	numLimit, err := strconv.Atoi(limit)
	if err != nil {
		numLimit, _ = strconv.Atoi(defaultLimit)
	}

	to := c.QueryParam("to")
	if to == "" {
		to = defaultTo
	}
	numTo, err := strconv.Atoi(to)
	if err != nil {
		numTo, _ = strconv.Atoi(defaultTo)
	}
	readableTo := time.Unix(int64(numTo)/1000, 0).Format(time.RFC3339)

	klines, err := bybit.GetKLine(ticker, interval, to, numLimit)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %v", err))
	}
	klines.Meta.Params = model.KLineParams{
		Symbol:   ticker,
		Interval: interval,
		Limit:    numLimit,
		To:       readableTo,
	}

	if len(klines.Data) > 0 {
		klines.Meta.Params.From = klines.Data[0].Datetime
	}

	analysis := make([]model.KLineAnalysisData, 0)
	for i, kline := range klines.Data {
		tick := 0.0
		if i > 0 {
			previousTimestamp := analysis[i-1].KLine.Timestamp
			currentTimestamp := kline.Timestamp
			tick = (float64(currentTimestamp) - float64(previousTimestamp)) / 60000
		}
		aKline := kline.ToAnalysis(tick)
		analysis = append(analysis, aKline)
	}
	mtx := utils.GetRanges(len(klines.Data))
	gap := len(analysis) - len(mtx)

	recommendations := make([]model.ActionRecommendation, 0)

	lastStopLoss := 0.0
	for i, slice := range mtx {
		j := i + gap
		start := slice[0]
		end := slice[1]
		previousItem := model.KLineAnalysisData{}
		if i > 0 {
			previousItem = analysis[j-1]
		}
		history := analysis[start:end]
		analysis[j].History = history
		pkg.EnhanceAverageData(&analysis[j], previousItem)

		if i > 1 {
			stopLoss := analysis[i-1].KLine.Low
			risk := math.Abs(analysis[j-1].KLine.Low - analysis[j].KLine.Close)
			recommendation := model.ActionRecommendation{
				Datetime: analysis[j].KLine.Datetime,
				StopLoss: stopLoss,
				Risk:     analysis[j].KLine.Close + risk,
			}
			if lastStopLoss < stopLoss {
				recommendation.StopLoss = stopLoss
				lastStopLoss = stopLoss
			}

			bybit.GenerateRecommendation(&recommendation, history)
			if !(bybit.CurrentRecommendation.Type == "" && recommendation.Type == "sell") &&
				recommendation.Type != "wait" &&
				bybit.CurrentRecommendation.Type != recommendation.Type {
				recommendations = append(recommendations, recommendation)
				bybit.CurrentRecommendation = recommendation
				bybit.LastRecommendedKline = analysis[j]
				if recommendation.Type == "sell" {
					lastStopLoss = 0.0
				}
			}
		}
	}

	response := model.KLineBacktestResponse{}
	profit := 0.0
	totalTransactions := 0.0

	potential := model.KLinePotential{
		Initial:   500,
		Final:     500,
		Variation: 0,
	}

	for i, r := range recommendations {
		if r.Datetime != "2024-12-13T17:51:00-04:00" && r.Type == "buy" {
			potential.Price = r.Close
			potential.Bought = potential.Final / r.Close
			potential.Final = 0
		} else if r.Type == "sell" && potential.Bought > 0 {
			potential.Price = r.Close
			v := potential.Bought * r.Close
			potential.Variation = v - potential.Initial
			recommendations[i].ProfitLoss = potential.Variation
			potential.Final = v
			potential.Bought = 0
			if potential.Variation > 0 {
				profit++
			}
			totalTransactions++
		}
	}
	// if potential.Bought > 0 {
	// 	v := potential.Bought * potential.Price
	// 	potential.Variation = v - potential.Initial
	// 	potential.Final = v
	// 	potential.Bought = 0
	// }

	if totalTransactions > 0 {
		potential.Efficiency = (profit / totalTransactions) * 100
	}
	response.Potential = potential
	response.Recommendations = recommendations

	return c.JSON(http.StatusOK, response)

}
