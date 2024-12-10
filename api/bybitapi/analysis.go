package bybitapi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/app/bybit"
	"github.com/marcioaso/consult/app/model"
)

func AnalysisHandler(c echo.Context) error {
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

	kline, err := bybit.GetKLine(ticker, interval, to, numLimit)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %v", err))
	}

	kline.Meta.Params = model.KLineParams{
		Symbol:   ticker,
		Interval: interval,
		Limit:    numLimit,
		To:       readableTo,
	}

	if len(kline.Data) > 0 {
		kline.Meta.Params.From = kline.Data[0].Datetime
	}

	return c.JSON(http.StatusOK, kline)

}
