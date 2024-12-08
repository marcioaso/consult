package bybit

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcioaso/consult/utils"
)

func (c *KLineData) Convert() KLineData {
	data := KLineData{
		Timestamp:      c.T,
		Symbol:         c.S,
		SymbolInternal: c.SN,
	}
	vol, _ := strconv.ParseFloat(c.V, 64)
	data.Volume = vol

	open, _ := strconv.ParseFloat(c.O, 64)
	data.Open = open

	close, _ := strconv.ParseFloat(c.C, 64)
	data.Close = close

	high, _ := strconv.ParseFloat(c.H, 64)
	data.High = high

	low, _ := strconv.ParseFloat(c.H, 64)
	data.Low = low

	return data
}

func GetKLine(ticker, interval, limit, to string) ([]KLineData, error) {
	url := getUrl(
		fmt.Sprintf(
			"/spot/api/quote/v2/klines?symbol=%s&interval=%s&limit=%s&to=%s",
			ticker,
			interval,
			limit,
			to,
		),
	)

	requestData := struct {
		Result []KLineData `json:"result"`
	}{}

	result, err := utils.Request(url, defaultHeaders)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, &requestData)
	if err != nil {
		return nil, err
	}

	response := make([]KLineData, 0)

	for _, item := range requestData.Result {
		response = append(response, item.Convert())
	}
	return response, nil
}
