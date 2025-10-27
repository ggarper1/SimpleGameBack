package services

import (
	"errors"
	"math"
	"math/rand/v2"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

const (
	numSegments = 4

	maxAttempts = 20

	maxSegmentLength float64 = 0.2
	minSegmentLength float64 = 0.1
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

func GenerateRandomSegments() ([numSegments]objects.Segment, error) {
	var segments [numSegments]objects.Segment

	idx := 0
	for idx < numSegments {
		added := false

		for attempts := 0; attempts < maxAttempts; attempts++ {
			start := objects.Point{rand.Float64(), rand.Float64()}
			end, err := generateValidEndPoint(start)
			if err != nil {
				continue
			}

			newSegment, err := objects.NewSegment(start, end)
			if err != nil {
				panic("this error should always be nil, there must be abug in generateValidEndPoint")
			}

			// Check if segment has conflict
			segments[idx] = newSegment
			added = true
			break
		}
		if !added {
			return segments, errors.New("could not generated segment under max number of attempts")
		}

		idx++
	}

	return segments, nil
}
