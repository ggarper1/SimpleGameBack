package objects

import "math"

type Point struct {
	X float64
	Y float64
}

func NewPointFromPolar(angle float64, rho float64) Point {
	return Point{math.Cos(angle) * rho, math.Sin(angle) * rho}
}

func (point Point) DistanceTo(other Point) float64 {
	dx := other.X - point.X
	dy := other.Y - point.Y

	return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
}
