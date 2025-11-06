package services_tests

import (
	"encoding/json"
	"errors"
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

type mapsToJSON struct {
	Data []json.RawMessage `json:"data"`
}

type mapsInJSON struct {
	Data []objects.Map `json:"data"`
}

func generateMaps() error {
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
					return errors.New(fmt.Sprintf("Error occured while adding piece:\n%v", err))
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
					return errors.New(fmt.Sprintf("Error occured while adding piece:\n%v", err))
				}
			}
		}

		jsonMaps[i] = m.ToDTO()
	}

	maps := mapsToJSON{jsonMaps[:]}
	jsonData, err := json.Marshal(maps)
	if err != nil {
		return errors.New("Could not encode to json")
	}

	err = os.WriteFile("maps.json", jsonData, 0o644)
	if err != nil {
		return errors.New("was not capable of writting to file")
	}
	return nil
}

func TestWinCondition(t *testing.T) {
	if _, err := os.Stat("maps.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Print("Generating maps...\n")
		generateMaps()
	}

	var testMaps mapsInJSON
	data, err := os.ReadFile("maps.json")
	if err != nil {
		panic("something went wrong while reading file")
	}
	json.Unmarshal(data, &testMaps)

	for i, m := range testMaps.Data {
		player := services.CheckWhoWon(&m)
		fmt.Printf("Player %d won in map %d\n", player, i)
	}
}
