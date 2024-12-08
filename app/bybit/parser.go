package bybit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/marcioaso/consult/pkg"
)

var smaConf = []int{25, 50, 100}

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

var invested = 100.00

func ParseKLineData(data []byte) (*KLineResponse, error) {
	var bought = 0.0
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

	smas1, _ := pkg.CalculateSMA(allCloses, smaConf[0])
	smas2, _ := pkg.CalculateSMA(allCloses, smaConf[1])
	smas3, _ := pkg.CalculateSMA(allCloses, smaConf[2])

	response := &KLineResponse{
		Meta: KLineMeta{
			SMAS: KLineSMAConfig{
				Fast:  smaConf[0],
				Slow:  smaConf[1],
				Heavy: smaConf[2],
			},
		},
		Resolutions: KLineResolutions{
			Advices: make([]ActionRecomendation, 0),
		},
	}

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
			recomendation := ActionRecomendation{
				Datetime: current.Datetime,
				Type:     "hold",
				Candles:  i,
				Close:    current.Close,
			}

			tailKLineData[i].Directions.Recomendation = recomendation
			recomendationHistory := 10
			if i > recomendationHistory {
				recomendation = generateRecomendation(&recomendation, tailKLineData[i-recomendationHistory:i])

				if recomendation.Type != "hold" && recomendation.Type != currentRecomendation.Type {
					response.Resolutions.Advices = append(response.Resolutions.Advices, recomendation)

					lastRecomendedKline = tailKLineData[i]
					if recomendation.Type == "buy" {
						bought = invested / recomendation.Close
						invested = 0
					} else if recomendation.Type == "sell" && bought > 0 {
						amount := bought * recomendation.Close
						fmt.Println("operation: ", amount)
						invested = amount
						bought = 0
					}
					currentRecomendation = recomendation
				}
			}
		}
		tailKLineData[i].SMAS = smas
	}

	klineData := tailKLineData[2:periodsLimit]
	response.Data = klineData
	response.Resolutions.ResultCount = len(klineData)
	fmt.Println("**** invested", invested)

	return response, nil
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

	if angle > 0 {
		return "up"
	}
	return "down"
}
