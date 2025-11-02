package objects

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
)

const (
	NumPieces   = 3
	numSegments = 4

	minSegmentSeperation = 0.01
	minPieceSeparation   = 0.005

	MidSectionHalfWidth = 0.1

	maxAttempts = 10

	maxSegmentLength float64 = 0.4
	minSegmentLength float64 = 0.15
)

type Map struct {
	Player1Segments []Segment `json:"player1Segments"`
	Player2Segments []Segment `json:"player2Segments"`

	Player1King Point `json:"player1King"`
	Player2King Point `json:"player2King"`

	Player1Pieces []Piece
	Player2Pieces []Piece
}

// Map Initializer Code
func generateValidEndPoint(start *Point, player Player) (*Point, error) {
	for range maxAttempts {
		angle := rand.Float64() * 2 * math.Pi

		var intersection1, intersection2 Point
		var doesIntersect1, doesIntersect2 bool
		line := NewLineFromAngle(*start, angle)

		if angle <= math.Pi/2 {
			rightEdge, _ := NewLine(Point{1, 1}, Point{1, 0})
			var otherEdge Line
			switch player {
			case Player1:
				otherEdge, _ = NewLine(Point{0, 1 + MidSectionHalfWidth}, Point{1, 1 + MidSectionHalfWidth})
			case Player2:
				otherEdge, _ = NewLine(Point{0, -MidSectionHalfWidth}, Point{1, -MidSectionHalfWidth})
			default:
				panic("player here should either be player 1 or player 2")
			}

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(otherEdge)
		} else if angle <= math.Pi {
			leftEdge, _ := NewLine(Point{0, 0}, Point{0, 1})
			var otherEdge Line
			switch player {
			case Player1:
				otherEdge, _ = NewLine(Point{0, 1 + MidSectionHalfWidth}, Point{1, 1 + MidSectionHalfWidth})
			case Player2:
				otherEdge, _ = NewLine(Point{0, -MidSectionHalfWidth}, Point{1, -MidSectionHalfWidth})
			default:
				panic("player here should either be player 1 or player 2")
			}

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(otherEdge)
		} else if angle <= 3*math.Pi/2 {
			leftEdge, _ := NewLine(Point{0, 0}, Point{0, 1})
			var otherEdge Line
			switch player {
			case Player1:
				otherEdge, _ = NewLine(Point{0, MidSectionHalfWidth}, Point{1, MidSectionHalfWidth})
			case Player2:
				otherEdge, _ = NewLine(Point{0, -1 - MidSectionHalfWidth}, Point{1, -1 - MidSectionHalfWidth})
			default:
				panic("player here should either be player 1 or player 2")
			}

			doesIntersect1, intersection1 = line.IntersectionLine(leftEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(otherEdge)
		} else {
			rightEdge, _ := NewLine(Point{1, 1}, Point{1, 0})
			var otherEdge Line
			switch player {
			case Player1:
				otherEdge, _ = NewLine(Point{0, MidSectionHalfWidth}, Point{1, MidSectionHalfWidth})
			case Player2:
				otherEdge, _ = NewLine(Point{0, -1 - MidSectionHalfWidth}, Point{1, -1 - MidSectionHalfWidth})
			default:
				panic("player here should either be player 1 or player 2")
			}

			doesIntersect1, intersection1 = line.IntersectionLine(rightEdge)
			doesIntersect2, intersection2 = line.IntersectionLine(otherEdge)
		}

		maxRho := maxSegmentLength
		if doesIntersect1 {
			distance := start.DistanceTo(intersection1)
			if distance < maxRho {
				maxRho = distance
			}
		}
		if doesIntersect2 {
			distance := start.DistanceTo(intersection2)
			if distance < maxRho {
				maxRho = distance
			}
		}

		if maxRho > minSegmentLength {
			rho := rand.Float64()*(maxRho-minSegmentLength) + minSegmentLength
			return &Point{start.X + rho*math.Cos(angle), start.Y + rho*math.Sin(angle)}, nil
		}
	}
	return &Point{}, errors.New("could not generate endpoint under max number of attempts")
}

func generateRandomSegments(player Player) ([]Segment, error) {
	var segments [numSegments]Segment

	for idx := range numSegments {
		added := false

		for range maxAttempts {
			var start Point
			switch player {
			case Player1:
				start = Point{rand.Float64(), rand.Float64() + MidSectionHalfWidth}
			case Player2:
				start = Point{rand.Float64(), -rand.Float64() - MidSectionHalfWidth}
			default:
				panic("player here should be either player 1 or 2")
			}
			end, err := generateValidEndPoint(&start, player)
			if err != nil {
				continue
			}

			newSegment, err := NewSegment(start, *end)
			if err != nil {
				panic("this error should always be nil, there must be a bug in generateValidEndPoint")
			}

			isValid := true
			for _, segment := range segments {
				if newSegment.ShortestDistanceToSegment(segment) < minSegmentSeperation {
					isValid = false
					break
				}
			}

			if isValid {
				segments[idx] = newSegment
				added = true
				break
			}
		}
		if !added {
			return segments[:], errors.New("could not generated segment under max number of attempts")
		}
	}

	return segments[:], nil
}

func generateKingPiece(segments []Segment, player Player) (*Point, error) {
	for range maxAttempts {
		viewPoint := Point{rand.Float64(), 0}

		for range maxAttempts {
			var kingPiece Point
			switch player {
			case Player1:
				kingPiece = Point{rand.Float64(), rand.Float64() + MidSectionHalfWidth}
			case Player2:
				kingPiece = Point{rand.Float64(), -rand.Float64() - MidSectionHalfWidth}
			default:
				panic("player here should be either player 1 or 2")
			}

			fakeSegment, err := NewSegment(viewPoint, kingPiece)
			for err != nil {
				fakeSegment, err = NewSegment(viewPoint, kingPiece)
			}

			isValid := true
			for _, segment := range segments {
				doIntersect, _ := fakeSegment.IntersectionSegment(segment)
				if doIntersect {
					isValid = false
					break
				}
			}
			if isValid {
				return &kingPiece, nil
			}
		}
	}
	return &Point{}, errors.New("could not generate valid king piece under max attempts")
}

func NewMap() Map {
	player1Segments, err := generateRandomSegments(Player1)
	if err != nil {
		panic(fmt.Sprintf("Could not create map: %v", err))
	}

	player2Segments, err := generateRandomSegments(Player2)
	if err != nil {
		panic(fmt.Sprintf("could not create map: %v", err))
	}

	player1king, err := generateKingPiece(player1Segments, Player1)
	if err != nil {
		panic(fmt.Sprintf("could not create map: %v", err))
	}

	player2king, err := generateKingPiece(player2Segments, Player2)
	if err != nil {
		panic(fmt.Sprintf("could not create map: %v", err))
	}

	return Map{
		Player1Segments: player1Segments,
		Player2Segments: player2Segments,
		Player1King:     *player1king,
		Player2King:     *player2king,
		Player1Pieces:   make([]Piece, 0, NumPieces),
		Player2Pieces:   make([]Piece, 0, NumPieces),
	}
}

// Other Map functions
func (m Map) ToDTO() []byte {
	// Convert to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		panic("unable to jsonify this map")
	}

	return jsonData
}

func (m *Map) AddPlayer1Piece(piece Piece) (bool, error) {
	if len(m.Player1Pieces) == NumPieces {
		return false, errors.New("trying to add more peices than allowed")
	}

	if piece.Position.X < 0 || piece.Position.X > 1 {
		return false, errors.New("Invalid X coordinate")
	}
	if piece.Position.Y < MidSectionHalfWidth || piece.Position.Y > 1+MidSectionHalfWidth {
		return false, errors.New("Invalid Y coordinate")
	}

	for _, segment := range m.Player1Segments {
		if segment.ShortestDistanceToPoint(piece.Position) < minPieceSeparation {
			return false, nil
		}
	}
	if piece.Position.DistanceTo(m.Player1King) < minPieceSeparation {
		return false, nil
	}
	for _, otherPiece := range m.Player1Pieces {
		if piece.Position.DistanceTo(otherPiece.Position) < minPieceSeparation {
			return false, nil
		}
	}

	m.Player1Pieces = append(m.Player1Pieces, piece)
	return true, nil
}

func (m *Map) AddPlayer2Piece(piece Piece) (bool, error) {
	if len(m.Player2Pieces) == NumPieces {
		return false, errors.New("trying to add more peices than allowed")
	}

	if piece.Position.X < 0 || piece.Position.X > 1 {
		return false, errors.New("Invalid X coordinate")
	}
	if piece.Position.Y > -MidSectionHalfWidth || piece.Position.Y < -1-MidSectionHalfWidth {
		return false, errors.New("Invalid Y coordinate")
	}

	for _, segment := range m.Player2Segments {
		if segment.ShortestDistanceToPoint(piece.Position) < minPieceSeparation {
			return false, nil
		}
	}
	if piece.Position.DistanceTo(m.Player2King) < minPieceSeparation {
		return false, nil
	}
	for _, otherPiece := range m.Player2Pieces {
		if piece.Position.DistanceTo(otherPiece.Position) < minPieceSeparation {
			return false, nil
		}
	}

	m.Player2Pieces = append(m.Player2Pieces, piece)
	return true, nil
}
