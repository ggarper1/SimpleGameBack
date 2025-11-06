package routes

import (
	"encoding/json"
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

type player struct {
	conn   *websocket.Conn
	player objects.Player
}

type MatchCreator struct {
	connections chan *websocket.Conn
	matches     chan *match
}

type match struct {
	player1  player
	player2  player
	matchMap objects.Map
}

type matchDTO struct {
	Map    json.RawMessage `json:"map"`
	Player int             `json:"player"`
}

type piecePlacement struct {
	Pieces []objects.Piece `json:"pieces"`
}

type matchResult struct {
	ValidPieces bool           `json:"validPieces"`
	Result      objects.Player `json:"result"`
}

func NewMatchCreator() *MatchCreator {
	manager := &MatchCreator{
		connections: make(chan *websocket.Conn),
		matches:     make(chan *match),
	}
	go manager.createMatches()
	go manager.handleMatches()
	return manager
}

func (matchCreator MatchCreator) RecieveConnection(w http.ResponseWriter, r *http.Request) {
	// Step 1: Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	matchCreator.connections <- conn
}

func (matchCreator MatchCreator) createMatches() {
	player1 := player{
		conn:   <-matchCreator.connections,
		player: objects.Player1,
	}

	player2 := player{
		conn:   <-matchCreator.connections,
		player: objects.Player2,
	}

	newMatch := &match{
		player1:  player1,
		player2:  player2,
		matchMap: objects.NewMap(),
	}

	matchCreator.matches <- newMatch
}

func getPieces(conn *websocket.Conn, pieceChan chan objects.Piece) {
	var pieces piecePlacement
	err := conn.ReadJSON(&pieces)
	if err != nil {
		panic("Error recieving pieces")
	}
	for _, p := range pieces.Pieces {
		pieceChan <- p
	}
}

func (matchCreator MatchCreator) handleMatches() {
	match := <-matchCreator.matches

	defer match.player1.conn.Close()
	defer match.player2.conn.Close()

	mapDTO := match.matchMap.ToDTO()

	// Send match data to player 1
	player1MatchData := matchDTO{
		Map:    mapDTO,
		Player: 1,
	}
	player1Data, err := json.Marshal(player1MatchData)
	if err != nil {
		panic("Error marshaling player 1 data")
	}
	match.player1.conn.WriteMessage(websocket.TextMessage, player1Data)

	// Send match data to player 2
	player2MatchData := matchDTO{
		Map:    mapDTO,
		Player: 2,
	}
	player2Data, err := json.Marshal(player2MatchData)
	if err != nil {
		panic("Error marshaling player 2 data")
	}
	match.player2.conn.WriteMessage(websocket.TextMessage, player2Data)

	// Get pieces
	player1Pieces := make(chan objects.Piece)
	go getPieces(match.player1.conn, player1Pieces)
	player2Pieces := make(chan objects.Piece)
	go getPieces(match.player2.conn, player2Pieces)

	player1Failed := false
	for piece := range player1Pieces {
		valid, err := match.matchMap.AddPlayer1Piece(piece)
		if err != nil {
			panic("Error when placing piece for player 2")
		}
		if !valid {
			player1Failed = true
			break
		}
	}
	player2Failed := false
	for piece := range player2Pieces {
		valid, err := match.matchMap.AddPlayer2Piece(piece)
		if err != nil {
			panic("Error when placing piece for player 2")
		}
		if !valid {
			player2Failed = true
			break
		}
	}

	result := objects.NoPlayer
	if !player1Failed && !player2Failed {
		result = services.CheckWhoWon(&match.matchMap)
	}

	player1MatchResult := matchResult{
		ValidPieces: player1Failed,
		Result:      result,
	}
	player2MatchResult := matchResult{
		ValidPieces: player2Failed,
		Result:      result,
	}

	match.player1.conn.WriteJSON(player1MatchResult)
	match.player2.conn.WriteJSON(player2MatchResult)
}
