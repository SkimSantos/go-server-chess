package chessbot

import (
	"github.com/notnil/chess"
)

func ApplyMove(state string, move string) (string, error) {
	fen, err := chess.FEN(state)
	if err != nil {
		return state, err
	}
	game := chess.NewGame(fen)

	// Parse the move (e.g., "e2e4") and apply it
	if err := game.MoveStr(move); err != nil {
		return "", err
	}

	// Return the new state in FEN notation
	return game.Position().String(), nil
}
