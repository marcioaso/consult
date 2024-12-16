package bybitapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/app/bybit"
	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/config"
	"github.com/marcioaso/consult/db"
	"github.com/marcioaso/consult/pkg"
)

func BacktestHandler(c echo.Context) error {
	bybit.LastRecommendedKline = model.KLineAnalysisData{}
	bybit.PreviousRecommendation = model.ActionRecommendation{}

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

	db.Db.Connect()
	rawData := db.Db.GetKLineForBacktest()
	klineData := make([]model.KLineData, 0)
	for _, d := range rawData {
		klineData = append(klineData, d.ToData())
	}
	klines := model.KLineResponse{
		Data: klineData,
	}
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
	firstClose := 0.0
	analysis := make([]model.KLineAnalysisData, 0)
	for i, kline := range klines.Data {
		tick := 0.0
		if i > 0 {
			previousTimestamp := analysis[i-1].KLine.Timestamp
			currentTimestamp := kline.Timestamp
			tick = (float64(currentTimestamp) - float64(previousTimestamp)) / 60000
		} else {
			firstClose = kline.Close
		}
		aKline := kline.ToAnalysis(tick)
		analysis = append(analysis, aKline)
	}

	recommendations := make([]model.ActionRecommendation, 0)
	response := model.KLineBacktestResponse{}
	profits := 0.0
	losses := make([]model.ActionRecommendation, 0)
	totalTransactions := 0.0

	initialInvestiment := 5000.00

	bybit.Position = bybit.PositionBalance{
		Close:      0,
		Initial:    initialInvestiment,
		Final:      initialInvestiment,
		Multiplier: 1,
	}

	potential := model.KLinePotential{
		Initial:   initialInvestiment,
		Final:     initialInvestiment,
		Variation: 0,
	}

	emaRange := config.EMAS[0] * 2
	lastClose := 0.0
	for i := emaRange + 1; i < len(analysis); i++ {
		previousItem := analysis[i-1]
		history := analysis[i-emaRange : i]
		analysis[i].History = history
		pkg.EnhanceAverageData(&analysis[i], previousItem)
		recs := bybit.GenerateRecommendations(history)
		for _, each := range recs {
			each.Datetime = analysis[i].KLine.Datetime
			if each.Type == "buy" {
				bybit.Position.BuyPosition(each.Close, &each)
				recommendations = append(recommendations, each)
			} else if each.Type == "sell" {
				sell, profit := bybit.Position.SellPositions(each.Close)
				if len(sell) > 0 {
					if profit > 0 {
						profits++
					} else {
						losses = append(losses, each)
					}
					potential.Variation += profit
					totalTransactions++
					lastClose = analysis[i].KLine.Close
				}
				recommendations = append(recommendations, each)
			}
		}
		positions, amount := bybit.Position.CurrentHold(analysis[i].KLine.Close)
		fmt.Printf("t:%v, Positions: %v,  Amount: %v\n", analysis[i].KLine.Datetime, positions, amount)
	}
	if totalTransactions > 0 {
		potential.Efficiency = (profits / totalTransactions) * 100
	}
	_, hold := bybit.Position.CurrentHold(lastClose)
	potential.Final = potential.Initial + potential.Variation + hold
	potential.Variation = potential.Final - potential.Initial
	potential.WithoutBot = (initialInvestiment / firstClose) * lastClose

	response.Potential = potential
	response.Recommendations = recommendations
	response.Potential.Losses = losses

	return c.JSON(http.StatusOK, response)

}
