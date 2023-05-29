package game

import "testing"

func TestGame_IsWinner_EmptyBoard(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
	}

	var isWinner, winner = game.IsWinner()

	assertBool(t, isWinner, false)
	assertPiece(t, winner, EMPTY)
}

func TestGame_IsWinner_Row(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, X},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
	}

	var isWinner, winner = game.IsWinner()

	assertBool(t, isWinner, true)
	assertPiece(t, winner, X)
}

func TestGame_IsWinner_Column(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, EMPTY, O},
			{EMPTY, X, O},
		},
	}

	var isWinner, winner = game.IsWinner()

	assertBool(t, isWinner, true)
	assertPiece(t, winner, O)
}

func TestGame_IsWinner_Diag1(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, X, O},
			{EMPTY, O, X},
		},
	}

	var isWinner, winner = game.IsWinner()

	assertBool(t, isWinner, true)
	assertPiece(t, winner, X)
}

func TestGame_IsWinner_Diag2(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	var isWinner, winner = game.IsWinner()

	assertBool(t, isWinner, true)
	assertPiece(t, winner, O)
}

func TestGame_IsValidMove_Success(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.IsValidMove(3)

	assertBool(t, isValidMove, true)
}

func TestGame_IsValidMove_TooBig(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.IsValidMove(9)

	assertBool(t, isValidMove, false)

}

func TestGame_IsValidMove_TooSmall(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.IsValidMove(-1)

	assertBool(t, isValidMove, false)
}

func TestGame_IsValidMove_SquareAlreadyOccupied(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.IsValidMove(0)

	assertBool(t, isValidMove, false)
}

func TestGame_IsBoardFull_BoardFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{O, X, X},
			{X, O, O},
		},
	}

	isDraw := game.IsDraw()

	assertBool(t, isDraw, true)
}

func TestGame_IsBoardFull_BoardNotFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{O, X, X},
			{X, O, EMPTY},
		},
	}

	isDraw := game.IsDraw()

	assertBool(t, isDraw, false)
}

func TestGame_IsDraw_BoardFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{O, X, X},
			{X, O, O},
		},
	}

	isDraw := game.IsDraw()

	assertBool(t, isDraw, true)
}

func TestGame_IsDraw_BoardNotFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{X, EMPTY, O},
			{O, X, X},
		},
	}

	isDraw := game.IsDraw()

	assertBool(t, isDraw, false)
}

func TestGame_IsDraw_IsWin(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, X},
			{O, EMPTY, O},
			{O, X, EMPTY},
		},
	}

	isDraw := game.IsDraw()

	assertBool(t, isDraw, false)
}

func assertBool(t *testing.T, actual bool, expected bool) {
	if actual != expected {
		t.Errorf("Result was incorrect, got %t, should be: %t.", actual, expected)
	}
}

func assertPiece(t *testing.T, actual Piece, expected Piece) {
	if actual != expected {
		t.Errorf("Result was incorrect, got %s, should be: %s.", actual, expected)
	}
}
