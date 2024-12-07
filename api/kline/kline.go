package kline

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/utils"
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

	url := fmt.Sprintf("https://api2-2.bybit.com/spot/api/quote/v2/klines?symbol=%s&interval=%s&limit=%s&to=%s", ticker, interval, limit, to)

	headers := map[string]string{
		"Accept":             "application/json",
		"Content-Type":       "application/json",
		"Sec-CH-UA":          `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
		"Sec-CH-UA-Platform": `"Windows"`,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	}
	response, err := utils.Request(url, headers)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %v", err))
	}

	// Return the response from the Bybit API
	return c.JSON(http.StatusOK, response)
}
