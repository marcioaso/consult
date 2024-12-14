package bybit

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// randomPercentageVariation retorna um valor com variação de até 0.01% do valor original
func randomPercentageVariation(base float64, maxVariation float64) float64 {
	variation := base * maxVariation
	return base + (rand.Float64()*2-1)*variation
}

// generateData gera uma lista de BybitResponse aleatórios com variação limitada
func generateData(count int, initialPrice float64, maxVariation float64, minPrice float64) []BybitResponse {
	rand.Seed(uint64(time.Now().UnixNano()))
	data := make([]BybitResponse, count)

	// Inicializa o preço base
	lastClose := initialPrice

	for i := 0; i < count; i++ {
		open := randomPercentageVariation(lastClose, maxVariation) // Open depende do Close da vela anterior
		close := randomPercentageVariation(open, maxVariation)     // Close varia em relação ao Open
		high := randomPercentageVariation(open, maxVariation)      // High é baseado no Open
		if high < open {
			high = open
		}
		low := randomPercentageVariation(open, maxVariation) // Low é baseado no Open
		if low > open {
			low = open
		}

		// Evita que o preço caia abaixo do limite mínimo
		if close < minPrice {
			close = minPrice
		}

		volume := randomPercentageVariation(1000, maxVariation) // Volume aleatório com base em 1000

		data[i] = BybitResponse{
			T:  time.Now().Unix() + int64(i), // Timestamp sequencial
			V:  fmt.Sprintf("%.2f", volume),  // Volume
			O:  fmt.Sprintf("%.2f", open),    // Open
			C:  fmt.Sprintf("%.2f", close),   // Close
			H:  fmt.Sprintf("%.2f", high),    // High
			L:  fmt.Sprintf("%.2f", low),     // Low
			S:  "OK",                         // Status
			SN: fmt.Sprintf("SN-%04d", i+1),  // Serial aleatório
		}

		// Atualiza o valor de lastClose para a próxima vela
		lastClose = close
	}

	return data
}

func GenerateData(count int, initialPrice, maxVariation, minPrice float64) {
	// Gera 10 candles começando com um preço inicial de 100.0
	data := generateData(count, initialPrice, maxVariation, minPrice)

	// Imprime os dados em formato JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
