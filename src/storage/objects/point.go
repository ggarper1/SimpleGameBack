package objects

import "math"

type Point struct {
	X float64
	Y float64
}

func (point Point) DistanceTo(other Point) float64 {
	dx := other.X - point.X
	dy := other.Y - point.Y

	return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
}
