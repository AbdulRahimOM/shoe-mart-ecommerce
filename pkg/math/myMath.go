package myMath

import (
	"math"
)

func RoundFloat32(value float32, decimalPlaces float64) float32 {
	divisor := math.Pow(10, decimalPlaces)
	return float32(math.Round(float64(value)*divisor) / divisor)
}
