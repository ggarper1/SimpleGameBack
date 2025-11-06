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
	pieceSegment, err := NewSegment(*point, p.Position)
	if err != nil {
		panic("Either two pieces where generated with sane position or your passing same piece twice as argument")
	}

	segmentAngle := pieceSegment.Angle()

	if segmentAngle < p.Angle-fovSemiAngle || segmentAngle > p.Angle+fovSemiAngle {
		return false
	}

	for _, segment := range m.Player1Segments {
		doesIntersect, _ := pieceSegment.IntersectionSegment(segment)

		if doesIntersect {
			return false
		}
	}

	for _, segment := range m.Player2Segments {
		doesIntersect, _ := pieceSegment.IntersectionSegment(segment)

		if doesIntersect {
			return false
		}
	}

	return true
}
