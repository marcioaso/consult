package kline

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/app/bybit"
)

func KlineHandler(c echo.Context) error {
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
	to := c.QueryParam("to")
	if to == "" {
		to = defaultTo
	}

	kline, err := bybit.GetKLine(ticker, interval, limit, to)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %v", err))
	}

	return c.JSON(http.StatusOK, kline)

}
