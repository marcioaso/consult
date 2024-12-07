package bybit

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/marcioaso/consult/utils"
)

// type KLineRequestData struct {
// 	T  int    `json:"t,omitempty"`
// 	S  string `json:"s"`
// 	SN string `json:"sn"`
// 	C  string `json:"c"`
// 	H  string `json:"h"`
// 	L  string `json:"l"`
// 	O  string `json:"o"`
// 	V  string `json:"v"`
// }

type KLineData struct {
	// Response
	Timestamp      int     `json:"timestamp"`
	Symbol         string  `json:"symbol"`
	SymbolInternal string  `json:"symbol_internal"`
	Volume         float64 `json:"volume"`
	Open           float64 `json:"open"`
	Close          float64 `json:"close"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`

	// Request
	T  int    `json:"t,omitempty"`
	S  string `json:"s,omitempty"`
	SN string `json:"sn,omitempty"`
	C  string `json:"c,omitempty"`
	H  string `json:"h,omitempty"`
	L  string `json:"l,omitempty"`
	O  string `json:"o,omitempty"`
	V  string `json:"v,omitempty"`
}

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
