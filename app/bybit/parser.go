package bybit

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/pkg"
)

type BybitResponse struct {
	T  int64  `json:"t,omitempty"`
	V  string `json:"v,omitempty"`
	O  string `json:"o,omitempty"`
	C  string `json:"c,omitempty"`
	H  string `json:"h,omitempty"`
	L  string `json:"l,omitempty"`
	S  string `json:"s,omitempty"`
	SN string `json:"sn,omitempty"`
}

func (kl *BybitResponse) ToData() model.KLineData {
	kd := model.KLineData{
		Timestamp:      kl.T,
		Symbol:         kl.S,
		SymbolInternal: kl.SN,
	}
	v, _ := strconv.ParseFloat(kl.V, 64)
	kd.Volume = v

	o, _ := strconv.ParseFloat(kl.O, 64)
	kd.Open = o

	c, _ := strconv.ParseFloat(kl.C, 64)
	kd.Close = c

	h, _ := strconv.ParseFloat(kl.H, 64)
	kd.High = h

	l, _ := strconv.ParseFloat(kl.L, 64)
	kd.Low = l

	dt := time.Unix(kl.T/1000, 0).Format(time.RFC3339)
	kd.Datetime = dt

	kd.CloseOpen = kd.Close - kd.Open
	kd.CloseOpenPercent = (kd.CloseOpen / kd.Open) * 100

	if kd.Close > kd.Open {
		kd.Color = "green"
	} else {
		kd.Color = "red"
	}

	return kd
}

func ParseData(data []byte) (*model.KLineResponse, error) {
	serializer := struct {
		Result []BybitResponse `json:"result"`
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
