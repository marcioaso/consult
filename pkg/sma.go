package pkg

import "errors"

func CalculateLastSMA(data []float64, period int) (float64, error) {
	if period <= 0 {
		return 0.0, errors.New("period must be bigger than zero")
	}
	if len(data) < period {
		return 0.0, errors.New("not enough data")
	}
	sum := 0.0
	for j := 0; j < period; j++ {
		sum += data[j]
	}
	return sum / float64(period), nil
}

func CalculateSMA(data []float64, period int) ([]float64, error) {
	if period <= 0 {
		return nil, errors.New("period must be bigger than zero")
	}
	if len(data) < period {
		return nil, errors.New("not enough data")
	}
	sma := make([]float64, 0, len(data)-period+1)

	for i := 0; i <= len(data)-period; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += data[i+j]
		}
		sma = append(sma, sum/float64(period))
	}

	return sma, nil
}
