package services

import (
	"errors"
	"math"
	"math/rand/v2"

	"ggarper1/SimpleGameBack/src/storage/objects"

	"github.com/twitchyliquid64/golang-asm/obj"
)

const (
	numSegments = 4

	minSegmentSeperation = 0.1

	maxAttempts = 20

	maxSegmentLength float64 = 0.4
	minSegmentLength float64 = 0.2
)

func generateValidEndPoint(start objects.Point) (objects.Point, error) {
	ctr := 0

	for ctr < maxAttempts {
		angle := rand.Float64() * 2 * math.Pi

		var intersection1, intersection2 objects.Point
		var doesIntersect1, doesIntersect2 bool
		line := objects.NewLineFromAngle(start, angle)

		if angle <= math.Pi/2 {
			rightEdge, _ := objects.NewLine(objects.Point{1, 1}, objects.Point{1, 0})
			topEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{1, 0})

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(topEdge)
		} else if angle <= math.Pi {
			leftEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{0, 1})
			topEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{1, 0})

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(topEdge)
		} else if angle <= 3*math.Pi/2 {
			leftEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{0, 1})
			bottomEdge, _ := objects.NewLine(objects.Point{0, 1}, objects.Point{1, 1})

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(bottomEdge)
		} else {
			rightEdge, _ := objects.NewLine(objects.Point{1, 1}, objects.Point{1, 0})
			bottomEdge, _ := objects.NewLine(objects.Point{0, 1}, objects.Point{1, 1})

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
			return objects.NewPointFromPolar(angle, rho), nil
		}

		ctr++
	}
	return objects.Point{}, errors.New("could not generate endpoint under max number of attempts")
}

func generateRandomSegments() ([numSegments]objects.Segment, error) {
	var segments [numSegments]objects.Segment

	for idx := range numSegments {
		added := false

		for range maxAttempts {
			start := objects.Point{rand.Float64(), rand.Float64()}
			end, err := generateValidEndPoint(start)
			if err != nil {
				continue
			}

			newSegment, err := objects.NewSegment(start, end)
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
			return segments, errors.New("could not generated segment under max number of attempts")
		}
	}

	return segments, nil
}

func generateKingPiece(segments []objects.Segment) (objects.Point, error) {
	for range maxAttempts {
		viewPoint := objects.Point{rand.Float64(), 1}

		for range maxAttempts {
			kingPiece := objects.Point{rand.Float64(), rand.Float64()}

			fakeSegment, err := objects.NewSegment(viewPoint, kingPiece)
			for err != nil {
				fakeSegment, err = objects.NewSegment(viewPoint, kingPiece)
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
				return kingPiece, nil
			}
		}
	}
	return objects.Point{}, errors.New("could not generate valid king piece under max attempts")
}

func NewMap() objects.Map {
	player1Segments, err := generateRandomSegments()
	if err != nil {
		panic("could not create map")
	}

	player2Segments, err := generateRandomSegments()
	if err != nil {
		panic("could not create map")
	}

	player1king, err := generateKingPiece(player1Segments[:])
	if err != nil {
		panic("could not create map")
	}

	player2king, err := generateKingPiece(player2Segments[:])
	if err != nil {
		panic("could not create map")
	}

	return objects.Map{
		Player1Segments: player1Segments[:],
		Player2Segments: player2Segments[:],
		Player1King:     player1king,
		Player2King:     player2king,
	}
}
