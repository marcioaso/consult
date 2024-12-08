package pkg

import (
	"fmt"
	"time"
)

func Elapsed(process string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("Time taken by %s is %v\n", process, time.Since(start))
	}
}
