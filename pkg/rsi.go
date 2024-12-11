package pkg

import "fmt"

func CalculateRSI(closes []float64, period int) (float64, error) {
	if len(closes) < period {
		return 0, fmt.Errorf("not enough data to calculate RSI")
	}

	var gains, losses float64

	// Initial average gain/loss calculation
	for i := 1; i <= period; i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change // losses are positive values
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// Calculate RSI
	for i := period + 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			avgGain = ((avgGain * float64(period-1)) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) - change) / float64(period)
		}
	}

	// Avoid division by zero
	if avgLoss == 0 {
		return 100, nil
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi, nil
}
