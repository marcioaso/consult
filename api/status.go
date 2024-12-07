package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StatusHandler responde com o status do servidor
func StatusHandler(c echo.Context) error {
	response := map[string]bool{"online": true}
	return c.JSON(http.StatusOK, response)
}
