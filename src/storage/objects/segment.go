package objects

import "math"

type Segment struct {
	p1 Point
	p2 Point
}

func NewSegment(p1 Point, p2 Point) (bool, Segment) {
	if p1.X == p2.X && p1.Y == p2.Y {
		return false, Segment{}
	}
	return true, Segment{p1, p2}
}

func (segment Segment) IntersectionSegment(other Segment) (bool, Point) {
	// Formula from: https://en.wikipedia.org/wiki/Line–line_intersection
	denominator := (segment.p1.X-segment.p2.X)*(other.p1.Y-other.p2.Y) - (segment.p1.Y-segment.p2.Y)*(other.p1.X-other.p2.X)

	if denominator == 0 {
		return false, Point{}
	}

	tNumerator := (segment.p1.X-other.p1.X)*(other.p1.Y-other.p2.Y) - (segment.p1.Y-other.p1.Y)*(other.p1.X-other.p2.X)
	uNumerator := (segment.p1.X-segment.p2.X)*(segment.p1.Y-other.p2.Y) - (segment.p1.Y-segment.p2.Y)*(segment.p1.X-other.p1.X)

	if tNumerator/denominator < 0 || tNumerator/denominator > 1 ||
		uNumerator/denominator < 0 || uNumerator/denominator > 1 {
		return false, Point{}
	}

	xNumerator := (segment.p1.X*segment.p2.Y-segment.p1.Y*segment.p2.Y)*(other.p1.X-other.p2.X) - (segment.p1.X-segment.p2.X)*(other.p1.X*other.p2.Y-other.p1.Y*other.p2.X)
	yNumerator := (segment.p1.X*segment.p2.Y-segment.p1.Y*segment.p2.X)*(other.p1.Y-other.p2.Y) - (segment.p1.Y-segment.p2.Y)*(other.p1.X*other.p2.Y-other.p1.Y*other.p2.X)

	point := Point{xNumerator / denominator, yNumerator / denominator}

	return true, point
}

func (segment Segment) ShortestDistanceSegment(point Point) float64 {
	angle1 := math.Atan2(segment.p1.Y-point.Y, segment.p1.X-point.X)
	angle2 := math.Atan2(segment.p2.Y-point.Y, segment.p2.X-point.X)

	if angle1 > math.Pi/2 {
		return point.DistanceTo(segment.p1)
	} else if angle2 > math.Pi/2 {
		return point.DistanceTo(segment.p2)
	} else {
		dx := segment.p1.X - segment.p2.X
		dy := segment.p1.Y - segment.p2.Y

		return math.Abs(dy*point.X-dx*point.Y+segment.p2.X*segment.p1.Y-segment.p1.X*segment.p2.Y) / math.Sqrt(math.Pow(dy, 2)+math.Pow(dx, 2))
	}
}
