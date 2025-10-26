package objects_tests

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

type pointDistanceToTestCase struct {
	P1       point   `json:"p1"`
	P2       point   `json:"p2"`
	Distance float64 `json:"distance"`
}

type pointTestJSON struct {
	DistanceTest []pointDistanceToTestCase `json:"distance"`
}

func testPointDistanceTo(t *testing.T, testCases []pointDistanceToTestCase) {
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Point Distance Test Case %d", i), func(t *testing.T) {
			p1 := objects.Point{testCase.P1.X, testCase.P1.Y}
			p2 := objects.Point{testCase.P2.X, testCase.P2.Y}

			distance := p1.DistanceTo(p2)

			if CloseEnough(distance, testCase.Distance) {
				t.Logf("Test Case %d Passed.", i)
			} else {
				t.Errorf("Test Case %d Failed:\n\tPoint 1: (%.3f, %.3f)\n\tPoint 2: (%.3f, %.3f)\n\tExpected: %.3f, Got: %.3f\n\tDifference: %f",
					i, p1.X, p1.Y, p2.X, p2.Y, testCase.Distance, distance, math.Abs(distance-testCase.Distance))
			}
		})
	}
}

func TestPoint(t *testing.T) {
	jsonFile := "point.json"

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read test data file: %v", err)
	}

	var testData pointTestJSON
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	testPointDistanceTo(t, testData.DistanceTest)
}
