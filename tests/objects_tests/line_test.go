package objects_tests

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

type lineShortestDistanceToTestCase struct {
	Line     line    `json:"l"`
	Point    point   `json:"p"`
	Distance float64 `json:"distance"`
}

type lineIntersectionLineTestCase struct {
	L1           line   `json:"l1"`
	L2           line   `json:"l2"`
	Intersection *point `json:"intersection"`
}

type lineIntersectionSegmentTestCase struct {
	Line         line    `json:"l"`
	Segment      segment `json:"s"`
	Intersection *point  `json:"intersection"`
}

type lineTestJSON struct {
	DistanceToPointTest     []lineShortestDistanceToTestCase  `json:"distance"`
	IntersectionLineTest    []lineIntersectionLineTestCase    `json:"lineIntersection"`
	IntersectionSegmentTest []lineIntersectionSegmentTestCase `json:"segmentIntersection"`
}

// Distance to point
func testLineShortestDistanceTo(t *testing.T, testCases []lineShortestDistanceToTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Distance to Point Test Case %d", i), func(t *testing.T) {
			p1 := objects.Point{testCase.Line.P1.X, testCase.Line.P1.Y}
			p2 := objects.Point{testCase.Line.P2.X, testCase.Line.P2.Y}
			line, error := objects.NewLine(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line.", i)
				return
			}
			point := objects.Point{testCase.Point.X, testCase.Point.Y}

			distance := line.ShortestDistanceTo(point)
			if CloseEnough(distance, testCase.Distance) {
				t.Errorf("Test Case %d Failed:\n\tLine: ((%.3f, %.3f), (%.3f, %.3f))\n\tPoint: (%.3f, %.3f)\n\tExpected: %.3f, Got: %.3f",
					i, line.P1.X, line.P1.Y, line.P2.X, line.P2.Y, point.X, point.Y, testCase.Distance, distance)
			} else {
				t.Logf("Test Case %d Passed.", i)
			}
		})
	}
}

// Line intersection test
func testLineIntersectionLine(t *testing.T, testCases []lineIntersectionLineTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Distance to Line Test Case %d", i), func(t *testing.T) {
			p1 := objects.Point{testCase.L1.P1.X, testCase.L1.P1.Y}
			p2 := objects.Point{testCase.L1.P2.X, testCase.L1.P2.Y}
			l1, error := objects.NewLine(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line 1: (%.3f, %3.f), (%.3f, %3.f)", i, testCase.L1.P1.X, testCase.L1.P1.Y, testCase.L1.P2.Y, testCase.L1.P2.Y)
				return
			}

			p1 = objects.Point{testCase.L2.P1.X, testCase.L2.P1.Y}
			p2 = objects.Point{testCase.L2.P2.X, testCase.L2.P2.Y}
			l2, error := objects.NewLine(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line 1: (%.3f, %3.f), (%.3f, %3.f)", i, testCase.L2.P1.X, testCase.L2.P1.Y, testCase.L2.P2.Y, testCase.L2.P2.Y)
				return
			}

			doIntersect, intersection := l1.IntersectionLine(l2)

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
				t.Errorf("Test Case %d Failed:\n\tLine 1: ((%.3f, %.3f), (%.3f, %.3f))\n\tLine 2: ((%.3f, %.3f), (%.3f, %.3f))\n\tExpected: (%t), Got: (%t)",
					i, l1.P1.X, l1.P1.Y, l1.P2.X, l1.P2.Y, l2.P1.X, l2.P1.Y, l2.P2.X, l2.P2.Y, testCase.Intersection != nil, doIntersect)
				if doIntersect && testCase.Intersection != nil {
					t.Errorf("Expected: (%.3f, %.3f), Got: (%.3f, %.3f), Difference: (%f, %f)", testCase.Intersection.X, testCase.Intersection.Y, intersection.X, intersection.Y, math.Abs(testCase.Intersection.X-intersection.X), math.Abs(testCase.Intersection.Y-intersection.Y))
				} else if doIntersect {
					t.Errorf("Got: (%.3f, %.3f)", intersection.X, intersection.Y)
				}
			}
		})
	}
}

// Segment intersection test
func testLineIntersectionSegment(t *testing.T, testCases []lineIntersectionSegmentTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Distance to Segment Test Case %d", i), func(t *testing.T) {
			p1 := objects.Point{testCase.Line.P1.X, testCase.Line.P1.Y}
			p2 := objects.Point{testCase.Line.P2.X, testCase.Line.P2.Y}
			line, error := objects.NewLine(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a line: (%.3f, %3.f), (%.3f, %3.f)", i, testCase.Line.P1.X, testCase.Line.P1.Y, testCase.Line.P2.Y, testCase.Line.P2.Y)
				return
			}

			p1 = objects.Point{testCase.Segment.P1.X, testCase.Segment.P1.Y}
			p2 = objects.Point{testCase.Segment.P2.X, testCase.Segment.P2.Y}
			segment, error := objects.NewSegment(p1, p2)
			if error != nil {
				t.Logf("Test Case %d skipped: was not possible to generate a segment: (%.3f, %3.f), (%.3f, %3.f)", i, testCase.Segment.P1.X, testCase.Segment.P1.Y, testCase.Segment.P2.Y, testCase.Segment.P2.Y)
				return
			}

			doIntersect, intersection := line.IntersectionSegment(segment)

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
				t.Errorf("Test Case %d Failed:\n\tLine: ((%.3f, %.3f), (%.3f, %.3f))\n\tSegment: ((%.3f, %.3f), (%.3f, %.3f))\n\tExpected: (%t), Got: (%t)",
					i, line.P1.X, line.P1.Y, line.P2.X, line.P2.Y, segment.P1.X, segment.P1.Y, segment.P2.X, segment.P2.Y, testCase.Intersection != nil, doIntersect)
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
func TestLine(t *testing.T) {
	jsonFile := "line.json"

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}

	var testData lineTestJSON
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	testLineShortestDistanceTo(t, testData.DistanceToPointTest)
	testLineIntersectionLine(t, testData.IntersectionLineTest)
	testLineIntersectionSegment(t, testData.IntersectionSegmentTest)
}
