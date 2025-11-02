package services_tests

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"ggarper1/SimpleGameBack/src/services"
	"ggarper1/SimpleGameBack/src/storage/objects"
)

const (
	numTests = 10
)

type mapsInJSON struct {
	Data []json.RawMessage `json:"data"`
}

func TestWinCondition(t *testing.T) {
	var jsonMaps [numTests]json.RawMessage
	for i := range numTests {
		m := objects.NewMap()

		for range objects.NumPieces {
			added := false
			for !added {
				angle := rand.Float64()*(math.Pi/2) + math.Pi*5.0/4
				point := objects.Point{rand.Float64(), rand.Float64() + objects.MidSectionHalfWidth}
				piece := objects.Piece{point, angle}
				var err error
				added, err = m.AddPlayer1Piece(piece)
				if err != nil {
					panic(fmt.Sprintf("Error occured while adding piece:\n%v", err))
				}
			}

			added = false
			for !added {
				angle := rand.Float64()*math.Pi/2 + math.Pi/4
				point := objects.Point{rand.Float64(), -rand.Float64() - objects.MidSectionHalfWidth}
				piece := objects.Piece{point, angle}
				var err error
				added, err = m.AddPlayer2Piece(piece)
				if err != nil {
					panic(fmt.Sprintf("Error occured while adding piece:\n%v", err))
				}
			}
		}

		jsonMaps[i] = m.ToDTO()

		player := services.CheckWhoWon(&m)

		fmt.Printf("Map %d: Player %d Wone\n", i, player)

	}

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
