package objects_tests

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

type segmentShortestDistanceToPointTestCase struct {
	Segment  segment `json:"s"`
	Point    point   `json:"p"`
	Distance float64 `json:"distance"`
}

type segmentShortestDistanceToSegmentTestCase struct {
	S1       segment `json:"s1"`
	S2       segment `json:"s2"`
	Distance float64 `json:"distance"`
}

type segmentIntersectionSegmentTestCase struct {
	S1           segment `json:"s1"`
	S2           segment `json:"s2"`
	Intersection *point  `json:"intersection"`
}

type segmentTestJSON struct {
	ShortestDistanceToPointTest   []segmentShortestDistanceToPointTestCase   `json:"pointDistance"`
	ShortestDistanceToSegmentTest []segmentShortestDistanceToSegmentTestCase `json:"segmentDistance"`
	IntersectionSegmentTest       []segmentIntersectionSegmentTestCase       `json:"segmentIntersection"`
}

// Distance to point
func testSegmentShortestDistanceToPoint(t *testing.T, testCases []segmentShortestDistanceToPointTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Distance to Point Test Case %d", i), func(t *testing.T) {
			p1 := objects.Point{testCase.Segment.P1.X, testCase.Segment.P1.Y}
			p2 := objects.Point{testCase.Segment.P2.X, testCase.Segment.P2.Y}
			segment, error := objects.NewSegment(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line.", i)
				return
			}

			point := objects.Point{testCase.Point.X, testCase.Point.Y}

			distance := segment.ShortestDistanceToPoint(point)
			if CloseEnough(distance, testCase.Distance) {
				t.Logf("Test Case %d Passed.", i)
			} else {
				t.Errorf("Test Case %d Failed:\n\tSegment: ((%.3f, %.3f), (%.3f, %.3f))\n\tPoint: (%.3f, %.3f)\n\tExpected: %.3f, Got: %.3f\n\tDifference: %f",
					i, segment.P1.X, segment.P1.Y, segment.P2.X, segment.P2.Y, point.X, point.Y, testCase.Distance, distance, math.Abs(testCase.Distance-distance))
			}
		})
	}
}

// Distance to segment
func testSegmentShortestDistanceToSegment(t *testing.T, testCases []segmentShortestDistanceToSegmentTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Distance to Segment Test Case %d", i), func(t *testing.T) {
			p1 := objects.Point{testCase.S1.P1.X, testCase.S1.P1.Y}
			p2 := objects.Point{testCase.S1.P2.X, testCase.S1.P2.Y}
			s1, error := objects.NewSegment(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line.", i)
				return
			}

			p1 = objects.Point{testCase.S2.P1.X, testCase.S2.P1.Y}
			p2 = objects.Point{testCase.S2.P2.X, testCase.S2.P2.Y}
			s2, error := objects.NewSegment(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line.", i)
				return
			}

			distance := s1.ShortestDistanceToSegment(s2)
			if CloseEnough(distance, testCase.Distance) {
				t.Logf("Test Case %d Passed.", i)
			} else {
				t.Errorf("Test Case %d Failed:\n\tS1: ((%.3f, %.3f), (%.3f, %.3f))\n\tS2: ((%.3f, %.3f),(%.3f, %.3f))\n\tExpected: %.3f, Got: %.3f\n\tDifference: %f",
					i, s1.P1.X, s1.P1.Y, s1.P2.X, s1.P2.Y, s2.P1.X, s2.P1.Y, s2.P2.X, s2.P2.Y, testCase.Distance, distance, math.Abs(testCase.Distance-distance))
			}
		})
	}
}

// Intersection with segment
func testSegmentIntersectionSegment(t *testing.T, testCases []segmentIntersectionSegmentTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Distance to Point Test Case %d", i), func(t *testing.T) {
			p11 := objects.Point{testCase.S1.P1.X, testCase.S1.P1.Y}
			p12 := objects.Point{testCase.S1.P2.X, testCase.S1.P2.Y}
			s1, error := objects.NewSegment(p11, p12)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line.", i)
				return
			}

			p21 := objects.Point{testCase.S2.P1.X, testCase.S2.P1.Y}
			p22 := objects.Point{testCase.S2.P2.X, testCase.S2.P2.Y}
			s2, error := objects.NewSegment(p21, p22)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line.", i)
				return
			}

			doIntersect, intersection := s1.IntersectionSegment(s2)
			passed := false
			if doIntersect && testCase.Intersection != nil {
				if CloseEnough(intersection.X, testCase.Intersection.X) && CloseEnough(intersection.Y, testCase.Intersection.Y) {
					t.Logf("Test Case %d Passed.", i)
					passed = true
				}
			} else if !doIntersect && testCase.Intersection == nil {
				t.Logf("Test Case %d Passed.", i)
				passed = true
			}

			if !passed {
				t.Errorf("Test Case %d Failed:\n\tSegment 1: ((%.3f, %.3f), (%.3f, %.3f))\n\tSegment 2: ((%.3f, %.3f), (%.3f, %.3f))\n\tExpected: (%t), Got: (%t)",
					i, s1.P1.X, s1.P1.Y, s1.P2.X, s1.P2.Y, s2.P1.X, s2.P1.Y, s2.P2.X, s2.P2.Y, testCase.Intersection != nil, doIntersect)
				if doIntersect && testCase.Intersection != nil {
					t.Errorf("Expected: (%.3f, %.3f), Got: (%.3f, %.3f), Difference: (%f, %f)", testCase.Intersection.X, testCase.Intersection.Y, intersection.X, intersection.Y, math.Abs(testCase.Intersection.X-intersection.X), math.Abs(testCase.Intersection.Y-intersection.Y))
				} else if doIntersect {
					t.Errorf("Got: (%.3f, %.3f)", intersection.X, intersection.Y)
				}
			}
		})
	}
}

// Whole test
func TestSegment(t *testing.T) {
	jsonFile := "segment.json"

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}

	var testData segmentTestJSON
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	testSegmentShortestDistanceToPoint(t, testData.ShortestDistanceToPointTest)
	testSegmentShortestDistanceToSegment(t, testData.ShortestDistanceToSegmentTest)
	testSegmentIntersectionSegment(t, testData.IntersectionSegmentTest)
}
