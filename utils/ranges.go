package utils

import (
	"slices"

	"github.com/marcioaso/consult/config"
)

func GetRanges(max int) [][]int {
	emas := slices.Max(config.EMAS) * 2
	starts := []int{slices.Max(config.SMAS), emas}
	start := slices.Max(starts)

	arr := make([][]int, 0)
	for i := start; i < max; i++ {
		arr = append(arr, []int{i - start, i})
	}
	return arr
}
