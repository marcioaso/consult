package utils

import (
	"fmt"
	"slices"
)

func GetRanges(max int, smas, emas []int) [][]int {
	starts := []int{slices.Max(smas), slices.Max(emas)}
	start := slices.Max(starts)
	gap := 25

	arr := make([][]int, 0)
	for i := start + gap; i < max; i++ {
		arr = append(arr, []int{i - gap, i})
	}
	fmt.Println(start, max, arr)

	return [][]int{}
}
