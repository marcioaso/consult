package bybitapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/app/bybit"
	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/db"
	"github.com/marcioaso/consult/utils"
)

func PopulateHandler(c echo.Context) error {
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

	url := fmt.Sprintf(
		"https://api2.bybit.com/spot/api/quote/v2/klines?symbol=%s&interval=%s&limit=%v&to=%s",
		ticker,
		interval,
		numLimit+1,
		to,
	)
	result, err := utils.Request(url, bybit.DefaultHeaders)
	if err != nil {
		return err
	}
	serializer := struct {
		Result []model.BybitResponse `json:"result"`
	}{}

	err = json.Unmarshal(result, &serializer)
	if err != nil {
		return err
	}
	from := serializer.Result[0].T

	klines := model.KLineResponse{}

	klines.Meta.Params = model.KLineParams{
		Symbol:   ticker,
		Interval: interval,
		Limit:    numLimit,
		To:       readableTo,
		From:     fmt.Sprintf("%v", from),
	}

	db.Db.Connect()
	db.Db.InsertKLines(serializer.Result)

	return c.JSON(http.StatusOK, klines)

}
