package routes_tests

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"ggarper1/SimpleGameBack/src/storage/objects"
	"github.com/gorilla/websocket"
)

type matchDTO struct {
	Map    json.RawMessage `json:"map"`
	Player int             `json:"player"`
}

type piecePlacement struct {
	Pieces []objects.Piece `json:"pieces"`
}

type matchResult struct {
	Valid  bool `json:"valid"`
	Result int  `json:"result"`
}

func TestMatch(t *testing.T) {
	for range 10 {
		// Test Case with valid pieces:
		go startMatch(t, false)
		startMatch(t, false)
	}
	for range 10 {
		// Test Case with valid pieces:
		go startMatch(t, false)
		startMatch(t, true)
	}
}

func startMatch(t *testing.T, invalid bool) {
	serverURL := "ws://localhost:8080/ws"

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		t.Error("Failed to connect to WebSocket: ", err)
	}

	defer conn.Close()

	var matchData matchDTO
	err = conn.ReadJSON(&matchData)
	if err != nil {
		t.Error("Failed to receive match data:", err)
	}

	pieces := generateRandomPieces(objects.NumPieces, invalid, matchData.Player)

	// Send pieces to server
	placement := piecePlacement{
		Pieces: pieces,
	}

	err = conn.WriteJSON(placement)
	if err != nil {
		t.Errorf("Player %d: Failed to send pieces: %v", matchData.Player, err)
	}

	// Receive match result
	var result matchResult
	err = conn.ReadJSON(&result)
	if err != nil {
		t.Errorf("Player %d: Failed to receive result: %v", matchData.Player, err)
	}

	if invalid && result.Valid || invalid && result.Result != 0 {
		s := ""
		for _, p := range pieces {
			s += fmt.Sprintf("\t(%f, %f), %f\n", p.Position.X, p.Position.Y, p.Angle)
		}
		t.Errorf("Invalid pieces were accepted as valid for player %d\n%s", matchData.Player, s)
	}

	fmt.Printf("Player %d: Match Result - Valid Pieces: %v, Winner: %d\n", matchData.Player, result.Valid, result.Result)
}

func generateRandomPieces(count int, invalid bool, player int) []objects.Piece {
	pieces := make([]objects.Piece, count)

	playerMult := 1.0
	if player == 2 {
		playerMult *= -1.0
	}

	for i := range count {
		if invalid {
			var x, y float64
			prob := rand.Float32()
			if prob < 1.0/3 {
				x = (rand.Float64()) * -1
				y = (rand.Float64() + 0.01) * -playerMult
			} else if prob < 2/3.0 {
				x = (rand.Float64() + 1.0) * playerMult
				y = (rand.Float64() + 1.0) * playerMult

			} else {
				x = rand.Float64() * 0.01 * playerMult
				y = rand.Float64() * 0.01 * playerMult
			}
			pieces[i] = objects.Piece{
				Position: objects.Point{x, y},
				Angle:    rand.Float64(),
			}
		} else {
			pieces[i] = objects.Piece{
				Position: objects.Point{rand.Float64(), (rand.Float64() + objects.MidSectionHalfWidth) * playerMult},
				Angle:    rand.Float64() * math.Pi * 2,
			}
		}
	}
	return pieces
}
