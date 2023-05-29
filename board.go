package main

import "fmt"

type Piece string

const (
	X     Piece = "X"
	O     Piece = "O"
	EMPTY Piece = "_"
)

func IsWinner(board [][]Piece) (bool, Piece) {
	for _, row := range board {
		if areItemsInArrayEqual(row) && row[0] != EMPTY {
			return true, row[0]
		}
	}

	diag1Count, diag2Count := 0, 0
	for i := 0; i < 3; i++ {

		count := 0
		for j := 0; j < 3; j++ {
			if board[j][i] == board[0][i] {
				count++
			}
		}

		if count >= 3 && board[0][i] != EMPTY {
			return true, board[0][i]
		}

		if board[i][i] == board[0][0] {
			diag1Count++
		}
		if board[i][2-i] == board[0][2] {
			diag2Count++
		}
	}

	if diag1Count >= 3 && board[0][0] != EMPTY {
		return true, board[0][0]
	}
	if diag2Count >= 3 && board[0][2] != EMPTY {
		return true, board[0][2]
	}

	return false, EMPTY
}

func areItemsInArrayEqual(arr []Piece) bool {
	for _, item := range arr {
		if item != arr[0] {
			return false
		}
	}
	return true
}

func IsValidMove(board [][]Piece, square int) bool {
	if square < 0 || square > 8 {
		return false
	}

	row, column := square/3, square%3

	if board[row][column] != EMPTY {
		return false
	}

	return true
}

func Print(board [][]Piece) {

	for _, row := range board {
		for _, cell := range row {
			fmt.Print(cell)
		}

		fmt.Println()
	}

}
