package service_tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"ggarper1/SimpleGameBack/src/services"
)

const (
	numTests = 10
)

type mapsInJSON struct {
	Data []json.RawMessage `json:"data"`
}

func TestMap(t *testing.T) {
	var jsonMaps [numTests]json.RawMessage
	for i := range numTests {
		t.Run(fmt.Sprintf("Point Distance Test Case %d", i), func(t *testing.T) {
			m := services.NewMap()

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
