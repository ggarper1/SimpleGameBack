package objects_tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

const (
	numTests = 10
)

type mapsInJSON struct {
	Data []json.RawMessage `json:"data"`
}

func TestMap(t *testing.T) {
	t.Log("This is a reminder that this test simply generates jsons and does some simple checks to verify\n that segments and king pieces coordinates are between 0 and 1.\n For more in depth testing visualized the results with `plot_maps.py`.")
	var jsonMaps [numTests]json.RawMessage
	for i := range numTests {
		t.Run(fmt.Sprintf("Point Distance Test Case %d", i), func(t *testing.T) {
			m := objects.NewMap()

			for _, segment := range m.Player1Segments {
				if segment.P1.X < 0 || segment.P1.X > 1 ||
					segment.P1.Y < 0 || segment.P1.Y > 1 ||
					segment.P2.X < 0 || segment.P2.X > 1 ||
					segment.P2.Y < 0 || segment.P2.Y > 1 {
					t.Errorf("Segment is not valid do to negative coordinates:\n\tStart: (%f, %f)\n\tEnd: (%f, %f)", segment.P1.X, segment.P1.Y, segment.P2.X, segment.P2.Y)
				}
			}
			for _, segment := range m.Player2Segments {
				if segment.P1.X < 0 || segment.P1.X > 1 ||
					segment.P1.Y < 0 || segment.P1.Y > 1 ||
					segment.P2.X < 0 || segment.P2.X > 1 ||
					segment.P2.Y < 0 || segment.P2.Y > 1 {
					t.Errorf("Segment is not valid do to negative coordinates for segments:\n\tStart: (%f, %f)\n\tEnd: (%f, %f)", segment.P1.X, segment.P1.Y, segment.P2.X, segment.P2.Y)
				}
			}
			if m.Player1King.X < 0 || m.Player1King.Y > 1 ||
				m.Player2King.X < 0 || m.Player2King.Y > 1 {
				t.Errorf("Segment is not valid do to negative coordinates for king pieces:\n\tKing piece 1: (%f, %f)\n\tKing piece 2: (%f, %f)", m.Player1King.X, m.Player1King.Y, m.Player2King.X, m.Player2King.Y)
			}

			jsonData := m.ToDTO()

			jsonMaps[i] = jsonData
		})

		maps := mapsInJSON{jsonMaps[:]}
		jsonData, err := json.Marshal(maps)
		if err != nil {
			panic("Could not encode to json")
		}

		err = os.WriteFile("maps.json", jsonData, 0o644)
		if err != nil {
			t.Error("Was not capable for writing json to file")
		}

	}
}
