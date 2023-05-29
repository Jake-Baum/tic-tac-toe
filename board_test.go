package main

import "testing"

func TestIsWinnerEmptyBoard(t *testing.T) {
	board := [][]Piece{
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
	}

	var isWinner, winner = IsWinner(board)

	assertBool(t, isWinner, false)
	assertPiece(t, winner, EMPTY)
}

func TestIsWinnerRow(t *testing.T) {
	board := [][]Piece{
		{X, X, X},
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
	}

	var isWinner, winner = IsWinner(board)

	assertBool(t, isWinner, true)
	assertPiece(t, winner, X)
}

func TestIsWinnerColumn(t *testing.T) {
	board := [][]Piece{
		{X, X, O},
		{EMPTY, EMPTY, O},
		{EMPTY, X, O},
	}

	var isWinner, winner = IsWinner(board)

	assertBool(t, isWinner, true)
	assertPiece(t, winner, O)
}

func TestIsWinnerDiag1(t *testing.T) {
	board := [][]Piece{
		{X, X, O},
		{EMPTY, X, O},
		{EMPTY, O, X},
	}

	var isWinner, winner = IsWinner(board)

	assertBool(t, isWinner, true)
	assertPiece(t, winner, X)
}

func TestIsWinnerDiag2(t *testing.T) {
	board := [][]Piece{
		{X, X, O},
		{EMPTY, O, O},
		{O, X, X},
	}

	var isWinner, winner = IsWinner(board)

	assertBool(t, isWinner, true)
	assertPiece(t, winner, O)
}

func TestIsValidMoveSuccess(t *testing.T) {
	board := [][]Piece{
		{X, X, O},
		{EMPTY, O, O},
		{O, X, X},
	}

	isValidMove := IsValidMove(board, 3)

	assertBool(t, isValidMove, true)
}

func TestIsValidMoveTooBig(t *testing.T) {
	board := [][]Piece{
		{X, X, O},
		{EMPTY, O, O},
		{O, X, X},
	}

	isValidMove := IsValidMove(board, 9)

	assertBool(t, isValidMove, false)

}

func TestIsValidMoveTooSmall(t *testing.T) {

	board := [][]Piece{
		{X, X, O},
		{EMPTY, O, O},
		{O, X, X},
	}

	isValidMove := IsValidMove(board, -1)

	assertBool(t, isValidMove, false)
}

func TestIsValidMoveSquareAlreadyOccupied(t *testing.T) {
	board := [][]Piece{
		{X, X, O},
		{EMPTY, O, O},
		{O, X, X},
	}

	isValidMove := IsValidMove(board, 0)

	assertBool(t, isValidMove, false)
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
