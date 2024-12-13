package pkg

import (
	"math"
)

func GetAngle(x1, y1, x2, y2 float64) float64 {
	deltaX := x2 - x1
	deltaY := y2 - y1

	angle := math.Atan2(deltaY, deltaX)
	angleDegrees := angle * (180 / math.Pi)

	return angleDegrees
}
