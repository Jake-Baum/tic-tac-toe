package game

import (
	"fmt"
	"unicode"
)

type Piece string

const (
	X     Piece = "X"
	O     Piece = "O"
	EMPTY Piece = "_"
)

type Game struct {
	Id          string
	Board       [][]Piece
	CurrentTurn Piece
}

func NewGame() *Game {
	return &Game{
		Board: [][]Piece{
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
		CurrentTurn: X,
	}
}

func (game Game) IsWinner() (bool, Piece) {
	for _, row := range game.Board {
		if areItemsInArrayEqual(row) && row[0] != EMPTY {
			return true, row[0]
		}
	}

	diag1Count, diag2Count := 0, 0
	for i := 0; i < 3; i++ {

		count := 0
		for j := 0; j < 3; j++ {
			if game.Board[j][i] == game.Board[0][i] {
				count++
			}
		}

		if count >= 3 && game.Board[0][i] != EMPTY {
			return true, game.Board[0][i]
		}

		if game.Board[i][i] == game.Board[0][0] {
			diag1Count++
		}
		if game.Board[i][2-i] == game.Board[0][2] {
			diag2Count++
		}
	}

	if diag1Count >= 3 && game.Board[0][0] != EMPTY {
		return true, game.Board[0][0]
	}
	if diag2Count >= 3 && game.Board[0][2] != EMPTY {
		return true, game.Board[0][2]
	}

	return false, EMPTY
}

func (game Game) IsBoardFull() bool {
	for _, row := range game.Board {
		for _, cell := range row {
			if cell == EMPTY {
				return false
			}
		}
	}

	return true
}

func (game Game) IsDraw() bool {
	if !game.IsBoardFull() {
		return false
	}

	isWinner, _ := game.IsWinner()

	return !isWinner
}

func (game Game) IsValidMove(square int) bool {
	if square < 0 || square > 8 {
		return false
	}

	row, column := square/3, square%3

	if game.Board[row][column] != EMPTY {
		return false
	}

	return true
}

func (game Game) Print() {

	for _, row := range game.Board {
		for _, cell := range row {
			fmt.Print(cell)
		}

		fmt.Println()
	}
}

func areItemsInArrayEqual(arr []Piece) bool {
	for _, item := range arr {
		if item != arr[0] {
			return false
		}
	}
	return true
}

func (game Game) Serialize() string {
	var serialized string
	for _, row := range game.Board {
		for _, cell := range row {
			serialized += string(cell)
		}
	}
	return serialized
}

func Deserialize(id string, state string, currentTurn string) Game {
	return Game{
		Id:          id,
		Board:       deserializeState(state),
		CurrentTurn: Piece(currentTurn),
	}
}

func deserializeState(state string) [][]Piece {
	return [][]Piece{
		{Piece(state[0]), Piece(state[1]), Piece(state[2])},
		{Piece(state[3]), Piece(state[4]), Piece(state[5])},
		{Piece(state[6]), Piece(state[7]), Piece(state[8])},
	}
}

func GetSquare(s string) (int, error) {
	if len(s) <= 0 || !unicode.IsDigit(rune(s[0])) {
		return -1, fmt.Errorf("%s is not a valid move selection", s)
	}

	return int(s[0]) - 48, nil
}
