package top10

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/app/bybit"
)

// Top10Handler responde ao endpoint /top10
func Top10Handler(c echo.Context) error {
	result, _ := bybit.GetTop10()
	return c.JSON(http.StatusOK, result)
}
