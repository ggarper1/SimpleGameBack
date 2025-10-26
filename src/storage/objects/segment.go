package objects

import (
	"errors"
	"math"
)

type Segment struct {
	P1 Point
	P2 Point
}

func NewSegment(p1 Point, p2 Point) (Segment, error) {
	if p1.X == p2.X && p1.Y == p2.Y {
		return Segment{}, errors.New("Cannot create segment with one point")
	}
	return Segment{p1, p2}, nil
}

func (segment Segment) ShortestDistanceTo(point Point) float64 {
	// Algorithm explanation:
	// 		Imagine a triangle formed by this segment's endpoints
	// 		and `point`. The shortest distance bewteen the segment
	// 		and this point will simply be the shortest distance
	// 		between the infinite line containignt the segment and
	// 		`point` unless one of the angles of the triangle at the
	// 		segment's endpoints is greater than 90 degrees.
	//
	// 		If the angle at vertice P (with P being one og the
	// 		segment's endpoints) is greater than 90 degrees the
	// 		the shortest distance of `point` to that segment is the
	// 		distance between `point` and P.
	//
	// 		We calculate if the angles are greater than 90 degrees
	// 		dot product.

	dxP1A := point.X - segment.P1.X
	dyP1A := point.Y - segment.P1.Y
	dxS := segment.P2.X - segment.P1.X
	dyS := segment.P2.Y - segment.P1.Y

	if dxP1A*dxS+dyP1A*dyS <= 0 {
		return point.DistanceTo(segment.P1)
	}

	dxP2A := point.X - segment.P2.X
	dyP2A := point.Y - segment.P2.Y
	dxS = -dxS
	dyS = -dyS

	if dxP2A*dxS+dyP2A*dyS <= 0 {
		return point.DistanceTo(segment.P2)
	}

	dx := segment.P1.X - segment.P2.X
	dy := segment.P1.Y - segment.P2.Y
	return math.Abs(dy*point.X-dx*point.Y+segment.P1.X*segment.P2.Y-segment.P2.X*segment.P1.Y) / math.Sqrt(math.Pow(dy, 2)+math.Pow(dx, 2))
}

func (segment Segment) IntersectionSegment(other Segment) (bool, Point) {
	// Formula from: https://en.wikipedia.org/wiki/Line–line_intersection
	denominator := (segment.P1.X-segment.P2.X)*(other.P1.Y-other.P2.Y) - (segment.P1.Y-segment.P2.Y)*(other.P1.X-other.P2.X)

	if denominator == 0 {
		return false, Point{}
	}

	t := ((segment.P1.X-other.P1.X)*(other.P1.Y-other.P2.Y) - (segment.P1.Y-other.P1.Y)*(other.P1.X-other.P2.X)) / denominator
	u := -((segment.P1.X-segment.P2.X)*(segment.P1.Y-other.P1.Y) - (segment.P1.Y-segment.P2.Y)*(segment.P1.X-other.P1.X)) / denominator

	if t < 0 || t > 1 ||
		u < 0 || u > 1 {
		return false, Point{}
	}

	x := segment.P1.X + t*(segment.P2.X-segment.P1.X)
	y := segment.P1.Y + t*(segment.P2.Y-segment.P1.Y)

	return true, Point{x, y}
}
