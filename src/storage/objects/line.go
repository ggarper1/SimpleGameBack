package objects

import "math"

type Line struct {
	p1 Point
	p2 Point
}

func NewLine(p1 Point, p2 Point) (bool, Line) {
	if p1.X == p2.X && p1.Y == p2.Y {
		return false, Line{}
	}
	return true, Line{p1, p2}
}

func (line Line) angle() float64 {
	dx := line.p1.X - line.p2.X
	dy := line.p1.Y - line.p2.Y
	return math.Atan2(dy, dx)
}

func (line Line) IntersectionLine(other Line) (bool, Point) {
	// Formula from: https://en.wikipedia.org/wiki/Line–line_intersection
	denominator := (line.p1.X-line.p2.X)*(other.p1.Y-other.p2.Y) - (line.p1.Y-line.p2.Y)*(other.p1.X-other.p2.X)

	if denominator == 0 {
		return false, Point{}
	}

	xNumerator := (line.p1.X*line.p2.Y-line.p1.Y*line.p2.Y)*(other.p1.X-other.p2.X) - (line.p1.X-line.p2.X)*(other.p1.X*other.p2.Y-other.p1.Y*other.p2.X)
	yNumerator := (line.p1.X*line.p2.Y-line.p1.Y*line.p2.X)*(other.p1.Y-other.p2.Y) - (line.p1.Y-line.p2.Y)*(other.p1.X*other.p2.Y-other.p1.Y*other.p2.X)

	return true, Point{xNumerator / denominator, yNumerator / denominator}
}

func (line Line) IntersectionSegment(segment Segment) (bool, Point) {
	// Formula from: https://en.wikipedia.org/wiki/Line–line_intersection
	denominator := (line.p1.X-line.p2.X)*(segment.p1.Y-segment.p2.Y) - (line.p1.Y-line.p2.Y)*(segment.p1.X-segment.p2.X)

	if denominator == 0 {
		return false, Point{}
	}

	u := ((line.p1.X-line.p2.X)*(line.p1.Y-segment.p1.Y) - (line.p1.Y-line.p2.Y)*(line.p1.X-segment.p1.X)) / ((line.p1.X-line.p2.X)*(segment.p1.Y-segment.p2.Y) - (line.p1.Y-line.p2.Y)*(segment.p1.X-segment.p2.X))
	if u < 0 || u > 1 {
		return false, Point{}
	}

	xNumerator := (line.p1.X*line.p2.Y-line.p1.Y*line.p2.Y)*(segment.p1.X-segment.p2.X) - (line.p1.X-line.p2.X)*(segment.p1.X*segment.p2.Y-segment.p1.Y*segment.p2.X)
	yNumerator := (line.p1.X*line.p2.Y-line.p1.Y*line.p2.X)*(segment.p1.Y-segment.p2.Y) - (line.p1.Y-line.p2.Y)*(segment.p1.X*segment.p2.Y-segment.p1.Y*segment.p2.X)

	x := xNumerator / denominator
	y := yNumerator / denominator

	return true, Point{x, y}
}

func (line Line) ShortestDistanceTo(point Point) float64 {
	dx := line.p1.X - line.p2.X
	dy := line.p1.Y - line.p2.Y

	return math.Abs(dy*point.X-dx*point.Y+line.p2.X*line.p1.Y-line.p1.X*line.p2.Y) / math.Sqrt(math.Pow(dy, 2)+math.Pow(dx, 2))
}
