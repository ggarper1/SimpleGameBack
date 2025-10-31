package objects

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
)

const (
	numSegments = 4

	minSegmentSeperation = 0.01

	maxAttempts = 10

	maxSegmentLength float64 = 0.4
	minSegmentLength float64 = 0.15
)

type Map struct {
	Player1Segments []Segment `json:"player1Segments"`
	Player2Segments []Segment `json:"player2Segments"`

	Player1King Point `json:"player1King"`
	Player2King Point `json:"player2King"`
}

// Map Initializer Code
func generateValidEndPoint(start *Point) (*Point, error) {
	for range maxAttempts {
		angle := rand.Float64() * 2 * math.Pi

		var intersection1, intersection2 Point
		var doesIntersect1, doesIntersect2 bool
		line := NewLineFromAngle(*start, angle)

		if angle <= math.Pi/2 {
			rightEdge, _ := NewLine(Point{1, 1}, Point{1, 0})
			topEdge, _ := NewLine(Point{0, 1}, Point{1, 1})

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(topEdge)
		} else if angle <= math.Pi {
			leftEdge, _ := NewLine(Point{0, 0}, Point{0, 1})
			topEdge, _ := NewLine(Point{0, 1}, Point{1, 1})

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(topEdge)
		} else if angle <= 3*math.Pi/2 {
			leftEdge, _ := NewLine(Point{0, 0}, Point{0, 1})
			bottomEdge, _ := NewLine(Point{0, 0}, Point{1, 0})

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(bottomEdge)
		} else {
			rightEdge, _ := NewLine(Point{1, 1}, Point{1, 0})
			bottomEdge, _ := NewLine(Point{0, 0}, Point{1, 0})

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(bottomEdge)
		}

		maxRho := maxSegmentLength
		if doesIntersect1 {
			distance := start.DistanceTo(intersection1)
			if distance < maxRho {
				maxRho = distance
			}
		}
		if doesIntersect2 {
			distance := start.DistanceTo(intersection2)
			if distance < maxRho {
				maxRho = distance
			}
		}

		if maxRho > minSegmentLength {
			rho := rand.Float64()*(maxRho-minSegmentLength) + minSegmentLength
			return &Point{start.X + rho*math.Cos(angle), start.Y + rho*math.Sin(angle)}, nil
		}
	}
	return &Point{}, errors.New("could not generate endpoint under max number of attempts")
}

func generateRandomSegments() ([]Segment, error) {
	var segments [numSegments]Segment

	for idx := range numSegments {
		added := false

		for range maxAttempts {
			start := Point{rand.Float64(), rand.Float64()}
			end, err := generateValidEndPoint(&start)
			if err != nil {
				continue
			}

			newSegment, err := NewSegment(start, *end)
			if err != nil {
				panic("this error should always be nil, there must be a bug in generateValidEndPoint")
			}

			isValid := true
			for _, segment := range segments {
				if newSegment.ShortestDistanceToSegment(segment) < minSegmentSeperation {
					isValid = false
					break
				}
			}

			if isValid {
				segments[idx] = newSegment
				added = true
				break
			}
		}
		if !added {
			return segments[:], errors.New("could not generated segment under max number of attempts")
		}
	}

	return segments[:], nil
}

func generateKingPiece(segments []Segment) (*Point, error) {
	for range maxAttempts {
		viewPoint := Point{rand.Float64(), 1}

		for range maxAttempts {
			kingPiece := Point{rand.Float64(), rand.Float64()}

			fakeSegment, err := NewSegment(viewPoint, kingPiece)
			for err != nil {
				fakeSegment, err = NewSegment(viewPoint, kingPiece)
			}

			isValid := true
			for _, segment := range segments {
				doIntersect, _ := fakeSegment.IntersectionSegment(segment)
				if doIntersect {
					isValid = false
					break
				}
			}
			if isValid {
				return &kingPiece, nil
			}
		}
	}
	return &Point{}, errors.New("could not generate valid king piece under max attempts")
}

func NewMap() Map {
	player1Segments, err := generateRandomSegments()
	if err != nil {
		panic(fmt.Sprintf("Could not create map: %v", err))
	}

	player2Segments, err := generateRandomSegments()
	if err != nil {
		panic(fmt.Sprintf("could not create map: %v", err))
	}

	player1king, err := generateKingPiece(player1Segments[:])
	if err != nil {
		panic(fmt.Sprintf("could not create map: %v", err))
	}

	player2king, err := generateKingPiece(player2Segments[:])
	if err != nil {
		panic(fmt.Sprintf("could not create map: %v", err))
	}

	return Map{
		Player1Segments: player1Segments[:],
		Player2Segments: player2Segments[:],
		Player1King:     *player1king,
		Player2King:     *player2king,
	}
}

// Other Map functions
func (m Map) ToDTO() []byte {
	// Convert to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		panic("unable to jsonify this map")
	}

	return jsonData
}
