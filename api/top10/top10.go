package top10

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/utils"
)

type Top10Receiver struct {
	AlgoType  string `json:"algoType"`
	BaseCoin  string `json:"baseCoin"`
	BuType    string `json:"buType"`
	QuoteCoin string `json:"quoteCoin"`
	SymbolId  string `json:"symbolId"`
}

// Top10Handler responde ao endpoint /top10
func Top10Handler(c echo.Context) error {
	// Definindo os cabeçalhos necessários
	headers := map[string]string{
		"Accept":             "application/json",
		"Content-Type":       "application/json",
		"Sec-CH-UA":          `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`,
		"Sec-CH-UA-Platform": `"Windows"`,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	}

	// Faz a requisição usando a função genérica
	result, err := utils.Request("https://api2.bybit.com/de/cht/api/searchSymbolRecommend/list", headers)
	if err != nil {
		log.Println("Error fetching data:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch data from external API",
		})
	}

	// Retorna a resposta da API externa como JSON
	return c.JSON(http.StatusOK, result)
}
