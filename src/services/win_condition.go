package services

import (
	"math/rand"

	"ggarper1/SimpleGameBack/src/storage/objects"
)

func CheckWhoWon(m *objects.Map) objects.Player {
	player1HitMatrix := [objects.NumPieces][objects.NumPieces]bool{}
	for j := range objects.NumPieces {
		for i := range objects.NumPieces {
			player1HitMatrix[j][i] = m.Player1Pieces[j].DoesHit(m, &m.Player2Pieces[i].Position)
		}
	}
	player1Assignations := [objects.NumPieces]int{}
	for j := range objects.NumPieces {
		validIndices := make([]int, 0, objects.NumPieces)
		for i := range objects.NumPieces {
			if player1HitMatrix[j][i] {
				validIndices = append(validIndices, i)
			}
		}
		if len(validIndices) > 0 {
			player1Assignations[j] = validIndices[rand.Intn(len(validIndices))]
		}
	}

	player2HitMatrix := [objects.NumPieces][objects.NumPieces]bool{}
	for j := range objects.NumPieces {
		for i := range objects.NumPieces {
			player2HitMatrix[j][i] = m.Player2Pieces[j].DoesHit(m, &m.Player1Pieces[i].Position)
		}
	}
	player2Assignations := [objects.NumPieces]int{}
	for j := range objects.NumPieces {
		validIndices := make([]int, 0, objects.NumPieces)
		for i := range objects.NumPieces {
			if player2HitMatrix[j][i] {
				validIndices = append(validIndices, i)
			}
		}
		if len(validIndices) > 0 {
			player2Assignations[j] = validIndices[rand.Intn(len(validIndices))]
		}
	}

	numPlayer1PiecesHitKing := 0
	for j := range objects.NumPieces {
		isDead := false
		for _, pieceIdx := range player2Assignations {
			if j == pieceIdx {
				isDead = true
				break
			}
		}

		if !isDead && m.Player1Pieces[j].DoesHit(m, &m.Player2King) {
			numPlayer1PiecesHitKing++
		}
	}

	numPlayer2PiecesHitKing := 0
	for j := range objects.NumPieces {
		isDead := false
		for _, pieceIdx := range player1Assignations {
			if j == pieceIdx {
				isDead = true
				break
			}
		}

		if !isDead && m.Player2Pieces[j].DoesHit(m, &m.Player1King) {
			numPlayer2PiecesHitKing++
		}
	}

	if numPlayer1PiecesHitKing > numPlayer2PiecesHitKing {
		return objects.Player1
	} else if numPlayer2PiecesHitKing > numPlayer1PiecesHitKing {
		return objects.Player2
	}
	return objects.NoPlayer
}
