package objects

import (
	"math"
)

const (
	fovSemiAngle = 5 * math.Pi / 180
)

type Piece struct {
	Position Point
	Angle    float64
}

func (p *Piece) DoesHit(m *Map, point *Point) bool {
	line, err := NewLine(p.Position, *point)
	if err != nil {
		panic("Either two pieces where generated with sane position or your passing same piece twice as argument")
	}

	lineAngle := line.Angle()

	if lineAngle < p.Angle-fovSemiAngle || lineAngle > p.Angle+fovSemiAngle {
		return false
	}

	for _, segment := range m.Player1Segments {
		doesIntersect, _ := line.IntersectionSegment(segment)

		if doesIntersect {
			return false
		}
	}

	for _, segment := range m.Player2Segments {
		doesIntersect, _ := line.IntersectionSegment(segment)

		if doesIntersect {
			return false
		}
	}

	return true
}
