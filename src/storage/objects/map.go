package objects

import (
	"encoding/json"
)

type Map struct {
	Player1Segments []Segment `json:"player1Segments"`
	Player2Segments []Segment `json:"player2Segments"`

	Player1King Point `json:"player1King"`
	Player2King Point `json:"player2King"`
}

func (m Map) ToDTO() []byte {
	// Convert to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		panic("unable to jsonify this map")
	}

	return jsonData
}
