package objects_tests

import (
	"math"
)

// Structs
type point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type line struct {
	P1 point `json:"p1"`
	P2 point `json:"p2"`
}

type segment struct {
	P1 point `json:"p1"`
	P2 point `json:"p2"`
}

var threshold float64 = 1e-4

func CloseEnough(f1 float64, f2 float64) bool {
	return math.Abs(f1-f2) < threshold
}
