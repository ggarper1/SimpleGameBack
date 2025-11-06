package objects

import (
	"errors"
	"math"
)

type Line struct {
	P1 Point `json:"p1"`
	P2 Point `json:"p2"`
}

func NewLine(p1 Point, p2 Point) (Line, error) {
	if p1.X == p2.X && p1.Y == p2.Y {
		return Line{}, errors.New("cannot create a line with just one point")
	}
	return Line{p1, p2}, nil
}

func NewLineFromAngle(p Point, angle float64) Line {
	return Line{p, Point{p.X + math.Cos(angle), p.Y + math.Sin(angle)}}
}

func (line Line) Angle() float64 {
	dx := line.P1.X - line.P2.X
	dy := line.P1.Y - line.P2.Y
	return math.Mod(math.Atan2(dy, dx)+math.Pi*2, math.Pi*2)
}

func (line Line) IntersectionLine(other Line) (bool, Point) {
	// Formula from: https://en.wikipedia.org/wiki/Line–line_intersection
	denominator := (line.P1.X-line.P2.X)*(other.P1.Y-other.P2.Y) - (line.P1.Y-line.P2.Y)*(other.P1.X-other.P2.X)

	if denominator == 0 {
		return false, Point{}
	}

	xNumerator := (line.P1.X*line.P2.Y-line.P1.Y*line.P2.X)*(other.P1.X-other.P2.X) - (line.P1.X-line.P2.X)*(other.P1.X*other.P2.Y-other.P1.Y*other.P2.X)
	yNumerator := (line.P1.X*line.P2.Y-line.P1.Y*line.P2.X)*(other.P1.Y-other.P2.Y) - (line.P1.Y-line.P2.Y)*(other.P1.X*other.P2.Y-other.P1.Y*other.P2.X)

	return true, Point{xNumerator / denominator, yNumerator / denominator}
}

func (line Line) IntersectionSegment(segment Segment) (bool, Point) {
	// Formula from: https://en.wikipedia.org/wiki/Line–line_intersection
	denominator := (line.P1.X-line.P2.X)*(segment.P1.Y-segment.P2.Y) - (line.P1.Y-line.P2.Y)*(segment.P1.X-segment.P2.X)

	if denominator == 0 {
		return false, Point{}
	}

	u := -((line.P1.X-line.P2.X)*(line.P1.Y-segment.P1.Y) - (line.P1.Y-line.P2.Y)*(line.P1.X-segment.P1.X)) / denominator

	if u < 0 || u > 1 {
		return false, Point{}
	}

	x := segment.P1.X + u*(segment.P2.X-segment.P1.X)
	y := segment.P1.Y + u*(segment.P2.Y-segment.P1.Y)

	return true, Point{x, y}
}

func (line Line) ShortestDistanceTo(point Point) float64 {
	dx := line.P1.X - line.P2.X
	dy := line.P1.Y - line.P2.Y

	return math.Abs(dy*point.X-dx*point.Y+line.P2.X*line.P1.Y-line.P1.X*line.P2.Y) / math.Sqrt(math.Pow(dy, 2)+math.Pow(dx, 2))
}
