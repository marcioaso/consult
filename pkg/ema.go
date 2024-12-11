package pkg

import "fmt"

var EmaConf = []int{25, 50}

func CalculateLastEMA(data []float64, period int) (float64, error) {
	if len(data) < period {
		return 0.0, fmt.Errorf("data length must be at least equal to the period")
	}

	mme := make([]float64, len(data))
	multiplier := 2.0 / float64(period+1)

	var sum float64
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	mme[period-1] = sum / float64(period)

	for i := period; i < len(data); i++ {
		mme[i] = (data[i]-mme[i-1])*multiplier + mme[i-1]
	}

	for i := 0; i < period-1; i++ {
		mme[i] = 0
	}

	return mme[len(mme)-1], nil
}

func CalculateEMA(data []float64, period int) ([]float64, error) {
	if len(data) < period {
		return nil, fmt.Errorf("data length must be at least equal to the period")
	}

	mme := make([]float64, len(data))
	multiplier := 2.0 / float64(period+1)

	var sum float64
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	mme[period-1] = sum / float64(period)

	for i := period; i < len(data); i++ {
		mme[i] = (data[i]-mme[i-1])*multiplier + mme[i-1]
	}

	for i := 0; i < period-1; i++ {
		mme[i] = 0
	}

	return mme, nil
}
