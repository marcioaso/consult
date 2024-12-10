package bybit

import (
	"fmt"

	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/utils"
)

func GetKLine(ticker, interval, to string, limit int) (*model.KLineResponse, error) {
	url := getUrl(
		fmt.Sprintf(
			"/spot/api/quote/v2/klines?symbol=%s&interval=%s&limit=%v&to=%s",
			ticker,
			interval,
			limit+1,
			to,
		),
	)

	result, err := utils.Request(url, defaultHeaders)
	if err != nil {
		return nil, err
	}

	response, err := ParseData(result)

	if err != nil {
		return nil, err
	}

	return response, nil
}
