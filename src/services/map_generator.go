package services

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

const (
	numSegments = 4

	minSegmentSeperation = 0.01

	maxAttempts = 10

	maxSegmentLength float64 = 0.4
	minSegmentLength float64 = 0.15
)

func generateValidEndPoint(start objects.Point) (objects.Point, error) {
	for range maxAttempts {
		angle := rand.Float64() * 2 * math.Pi

		fmt.Printf("------------------------------\n")
		fmt.Printf("\tAngle: %f\n", angle)

		var intersection1, intersection2 objects.Point
		var doesIntersect1, doesIntersect2 bool
		line := objects.NewLineFromAngle(start, angle)

		if angle <= math.Pi/2 {
			rightEdge, _ := objects.NewLine(objects.Point{1, 1}, objects.Point{1, 0})
			topEdge, _ := objects.NewLine(objects.Point{0, 1}, objects.Point{1, 1})

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(topEdge)
		} else if angle <= math.Pi {
			leftEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{0, 1})
			topEdge, _ := objects.NewLine(objects.Point{0, 1}, objects.Point{1, 1})

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(topEdge)
		} else if angle <= 3*math.Pi/2 {
			leftEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{0, 1})
			bottomEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{1, 0})

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(bottomEdge)
		} else {
			rightEdge, _ := objects.NewLine(objects.Point{1, 1}, objects.Point{1, 0})
			bottomEdge, _ := objects.NewLine(objects.Point{0, 0}, objects.Point{1, 0})

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(bottomEdge)
		}

		maxRho := maxSegmentLength
		if doesIntersect1 {
			distance := start.DistanceTo(intersection1)
			fmt.Printf("\tIntersetcion 1: (%f, %f), Distance: %f\n", intersection1.X, intersection1.Y, distance)
			if distance < maxRho {
				maxRho = distance
			}
		}
		if doesIntersect2 {
			distance := start.DistanceTo(intersection2)
			fmt.Printf("\tIntersetcion 2: (%f, %f), Distance: %f\n", intersection2.X, intersection2.Y, distance)
			if distance < maxRho {
				maxRho = distance
			}
		}
		fmt.Printf("\tMax rho: %f\n", maxRho)

		if maxRho > minSegmentLength {
			rho := rand.Float64()*(maxRho-minSegmentLength) + minSegmentLength
			fmt.Printf("\trho: %f\n", rho)
			fmt.Printf("\tx: %f\n", start.X+rho*math.Cos(angle))
			fmt.Printf("\ty: %f\n", start.Y+rho*math.Sin(angle))
			return objects.Point{start.X + rho*math.Cos(angle), start.Y + rho*math.Sin(angle)}, nil
		}
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

	return objects.Map{
		Player1Segments: player1Segments[:],
		Player2Segments: player2Segments[:],
		Player1King:     player1king,
		Player2King:     player2king,
	}
}
