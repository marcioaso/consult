package bybit

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/marcioaso/consult/pkg"
)

var periods = []int{25, 50, 100}

type SMAItem struct {
	Value         float64 `json:"value"`
	Angle         float64 `json:"angle"`
	PreviousAngle float64 `json:"previous_angle"`
}

type SMAData struct {
	FAST  SMAItem `json:"fast"`
	SLOW  SMAItem `json:"slow"`
	HEAVY SMAItem `json:"heavy"`
}

type Directions struct {
	Close string `json:"close"`
	Fast  string `json:"fast_sma"`
	Slow  string `json:"slow_sma"`
	Heavy string `json:"heavy_sma"`
}

type KLineData struct {
	T  int64  `json:"t,omitempty"`
	V  string `json:"v,omitempty"`
	O  string `json:"o,omitempty"`
	C  string `json:"c,omitempty"`
	H  string `json:"h,omitempty"`
	L  string `json:"l,omitempty"`
	S  string `json:"s,omitempty"`
	SN string `json:"sn,omitempty"`

	Datetime       string  `json:"datetime"`
	Timestamp      int64   `json:"timestamp"`
	Symbol         string  `json:"symbol"`
	SymbolInternal string  `json:"symbol_internal"`
	Volume         float64 `json:"volume"`
	Close          float64 `json:"close"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`

	SMAS       SMAData    `json:"sma"`
	Directions Directions `json:"directions"`
}

func (kl *KLineData) ToData() KLineData {
	kd := KLineData{
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

	return kd
}

func ParseKLineData(data []byte) ([]KLineData, error) {
	serializer := struct {
		Result []KLineData `json:"result"`
	}{}
	err := json.Unmarshal(data, &serializer)
	if err != nil {
		return nil, err
	}
	rawKlineData := make([]KLineData, 0)
	allCloses := make([]float64, 0)

	for _, each := range serializer.Result {
		item := each.ToData()
		rawKlineData = append(rawKlineData, item)
		allCloses = append(allCloses, item.Close)
	}

	smas1, _ := pkg.CalculateSMA(allCloses, periods[0])
	smas2, _ := pkg.CalculateSMA(allCloses, periods[1])
	smas3, _ := pkg.CalculateSMA(allCloses, periods[2])

	periodsLimit := len(smas3)
	originalRawKLinedata := len(rawKlineData)
	originalSma1Len := len(smas1)
	originalSma2Len := len(smas2)

	tailKLineData := rawKlineData[originalRawKLinedata-periodsLimit : originalRawKLinedata]

	tailSmas1 := smas1[originalSma1Len-periodsLimit : originalSma1Len]
	tailSmas2 := smas2[originalSma2Len-periodsLimit : originalSma2Len]

	for i := 0; i < periodsLimit; i++ {
		smas := SMAData{
			FAST: SMAItem{
				Value: tailSmas1[i],
			},
			SLOW: SMAItem{
				Value: tailSmas2[i],
			},
			HEAVY: SMAItem{
				Value: smas3[i],
			},
		}
		if i > 0 {
			previous := tailKLineData[i-1]

			current := tailKLineData[i]

			priceDirection := "down"
			if current.Close > previous.Close {
				priceDirection = "up"
			}
			tailKLineData[i].Directions.Close = priceDirection

			var timeTick = float64(current.Timestamp-previous.Timestamp) / 60000

			tailKLineData[i].Directions.Fast = enhanceSMAData(
				&smas.FAST,
				previous.SMAS.FAST,
				timeTick,
			)

			tailKLineData[i].Directions.Slow = enhanceSMAData(
				&smas.SLOW,
				previous.SMAS.SLOW,
				timeTick,
			)

			tailKLineData[i].Directions.Heavy = enhanceSMAData(
				&smas.HEAVY,
				previous.SMAS.HEAVY,
				timeTick,
			)

		}
		tailKLineData[i].SMAS = smas
	}

	return tailKLineData[2:periodsLimit], nil
}

func enhanceSMAData(data *SMAItem, previous SMAItem, timeTick float64) string {
	angle := pkg.GetAngle(
		0,
		previous.Value,
		timeTick,
		data.Value,
	)
	data.Angle = angle
	data.PreviousAngle = previous.Angle

	if angle > 0 && angle > previous.Angle {
		return "up"
	}
	return "down"
}
