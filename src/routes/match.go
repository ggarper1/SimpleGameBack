package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ggarper1/SimpleGameBack/src/services"
	"ggarper1/SimpleGameBack/src/storage/objects"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Player struct {
	conn   *websocket.Conn
	player objects.Player
}

type MatchCreator struct {
	connections chan *websocket.Conn
	matches     chan *Match
}

type Match struct {
	player1  Player
	player2  Player
	matchMap objects.Map
}

type MatchDTO struct {
	Map    json.RawMessage `json:"map"`
	Player int             `json:"player"`
}

type MatchResult struct {
	Valid  bool           `json:"valid"`
	Result objects.Player `json:"result"`
}

func NewMatchCreator() *MatchCreator {
	manager := &MatchCreator{
		connections: make(chan *websocket.Conn),
		matches:     make(chan *Match),
	}
	go manager.handleMatches()
	return manager
}

func (matchCreator MatchCreator) RecieveConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	matchCreator.connections <- conn
}

func (matchCreator MatchCreator) handleMatches() {
	for {
		player1 := Player{
			conn:   <-matchCreator.connections,
			player: objects.Player1,
		}

		player2 := Player{
			conn:   <-matchCreator.connections,
			player: objects.Player2,
		}

		match := &Match{
			player1:  player1,
			player2:  player2,
			matchMap: objects.NewMap(),
		}

		go matchCreator.handleSingleMatch(match)
	}
}

func getPieces(conn *websocket.Conn, pieceChan chan objects.Piece) {
	var pieces struct {
		Pieces []objects.Piece `json:"pieces"`
	}
	err := conn.ReadJSON(&pieces)
	if err != nil {
		panic(fmt.Sprintf("Error recieving pieces:\n %v", err))
	}
	for _, p := range pieces.Pieces {
		pieceChan <- p
	}
	close(pieceChan)
}

func (matchCreator MatchCreator) handleSingleMatch(m *Match) {
	defer m.player1.conn.Close()
	defer m.player2.conn.Close()

	mapDTO := m.matchMap.ToDTO()

	// Send match data to player 1
	player1MatchData := MatchDTO{
		Map:    mapDTO,
		Player: 1,
	}
	player1Data, err := json.Marshal(player1MatchData)
	if err != nil {
		panic("Error marshaling player 1 data")
	}
	err = m.player1.conn.WriteMessage(websocket.TextMessage, player1Data)
	if err != nil {
		panic("Could not send data to player 1")
	}

	// Send match data to player 2
	player2MatchData := MatchDTO{
		Map:    mapDTO,
		Player: 2,
	}
	player2Data, err := json.Marshal(player2MatchData)
	if err != nil {
		panic("Error marshaling player 2 data")
	}
	err = m.player2.conn.WriteMessage(websocket.TextMessage, player2Data)
	if err != nil {
		panic("Could not send data to player 1")
	}

	// Get pieces
	player1Pieces := make(chan objects.Piece, objects.NumPieces)
	go getPieces(m.player1.conn, player1Pieces)
	player2Pieces := make(chan objects.Piece, objects.NumPieces)
	go getPieces(m.player2.conn, player2Pieces)

	player1Failed := false
	ctr := 0
	for piece := range player1Pieces {
		valid, err := m.matchMap.AddPlayer1Piece(piece)
		if err != nil {
			panic("Error when placing piece for player 2")
		}
		if !valid {
			player1Failed = true
			break
		}
		ctr += 1
	}
	if ctr < objects.NumPieces && !player1Failed {
		panic("Expected more pieces")
	}

	player2Failed := false
	ctr = 0
	for piece := range player2Pieces {
		valid, err := m.matchMap.AddPlayer2Piece(piece)
		if err != nil {
			panic("Error when placing piece for player 2")
		}
		if !valid {
			player2Failed = true
			break
		}
		ctr += 1
	}
	if ctr < objects.NumPieces && !player2Failed {
		panic("Expected more pieces")
	}

	result := objects.NoPlayer
	if !player1Failed && !player2Failed {
		result = services.CheckWhoWon(&m.matchMap)
	}

	player1MatchResult := MatchResult{
		Valid:  !player1Failed,
		Result: result,
	}
	player2MatchResult := MatchResult{
		Valid:  !player2Failed,
		Result: result,
	}

	m.player1.conn.WriteJSON(player1MatchResult)
	m.player2.conn.WriteJSON(player2MatchResult)
}
