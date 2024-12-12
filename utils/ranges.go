package utils

import (
	"slices"

	"github.com/marcioaso/consult/config"
)

func GetRanges(max int) [][]int {
	starts := []int{slices.Max(config.SMAS), slices.Max(config.EMAS)}
	start := slices.Max(starts)

	arr := make([][]int, 0)
	for i := start; i < max; i++ {
		arr = append(arr, []int{i - start, i})
	}

	return arr
}
