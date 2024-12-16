package bybit

import (
	"encoding/json"

	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/pkg"
)

func ParseData(data []byte) (*model.KLineResponse, error) {
	serializer := struct {
		Result []model.BybitResponse `json:"result"`
	}{}

	err := json.Unmarshal(data, &serializer)
	if err != nil {
		return nil, err
	}
	rawKlineData := make([]model.KLineData, 0)

	for i, each := range serializer.Result {
		item := each.ToData()
		if i > 0 {
			previous := rawKlineData[i-1]
			tick := (float64(item.Timestamp) - float64(previous.Timestamp)) / 10000
			item.Angle = pkg.GetAngle(0, previous.Close, tick, item.Close)
		}
		rawKlineData = append(rawKlineData, item)
	}

	response := &model.KLineResponse{
		Data: rawKlineData,
	}
	response.Resolutions.ResultCount = len(rawKlineData)

	return response, nil
}
