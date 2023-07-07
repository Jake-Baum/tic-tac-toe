package game

import (
	"fmt"
)

type Piece string

type FinishedError struct {
	winner *Piece
}

func (e *FinishedError) Error() string {
	if e.winner != nil {
		return fmt.Sprintf("game has finished. The winner is %s!", *e.winner)
	}
	return "game has ended in a draw"
}

type InvalidMoveError struct {
	move int
}

func (e *InvalidMoveError) Error() string {
	return fmt.Sprintf("%d is not a valid move", e.move)
}

type NotPlayersTurnError struct {
	player string
}

func (e *NotPlayersTurnError) Error() string {
	return fmt.Sprintf("it is not %s's turn", e.player)
}

type PlayerDoesNotExistError struct {
	player string
	gameId string
}

func (e *PlayerDoesNotExistError) Error() string {
	return fmt.Sprintf("user %s is not permitted to make moves in game %s", e.player, e.gameId)
}

const (
	X     Piece = "X"
	O     Piece = "O"
	EMPTY Piece = "_"
)

type Game struct {
	Id          string
	Board       [][]Piece
	CurrentTurn Piece
	PlayerX     string
	PlayerO     string
}

func NewGame(playerX string, playerO string) *Game {
	return &Game{
		Board: [][]Piece{
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
		CurrentTurn: X,
		PlayerX:     playerX,
		PlayerO:     playerO,
	}
}

func (game *Game) MakeMove(player string, move int) error {
	if !game.IsPlayer(player) {
		return &PlayerDoesNotExistError{
			player: player,
			gameId: game.Id,
		}
	}

	if game.isDraw() {
		return &FinishedError{}
	}
	if isWinner, winner := game.isWinner(); isWinner {
		return &FinishedError{winner: &winner}
	}

	if !game.isValidMove(move) {
		return &InvalidMoveError{move: move}
	}

	if (game.CurrentTurn == X && player != game.PlayerX) || (game.CurrentTurn == O && player != game.PlayerO) {
		return &NotPlayersTurnError{player: player}
	}

	game.Board[move/3][move%3] = game.CurrentTurn

	if game.CurrentTurn == X {
		game.CurrentTurn = O
	} else {
		game.CurrentTurn = X
	}

	return nil
}

func (game *Game) IsPlayer(player string) bool {
	if game.PlayerX == player || game.PlayerO == player {
		return true
	}
	return false
}

func (game *Game) GetOtherPlayer(player string) string {
	if game.PlayerX == player {
		return game.PlayerO
	}
	return game.PlayerX
}

func (game *Game) isWinner() (bool, Piece) {
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

func (game *Game) isBoardFull() bool {
	for _, row := range game.Board {
		for _, cell := range row {
			if cell == EMPTY {
				return false
			}
		}
	}

	return true
}

func (game *Game) isDraw() bool {
	if !game.isBoardFull() {
		return false
	}

	isWinner, _ := game.isWinner()

	return !isWinner
}

func (game *Game) isValidMove(square int) bool {
	if square < 0 || square > 8 {
		return false
	}

	row, column := square/3, square%3

	if game.Board[row][column] != EMPTY {
		return false
	}

	return true
}

func areItemsInArrayEqual(arr []Piece) bool {
	for _, item := range arr {
		if item != arr[0] {
			return false
		}
	}
	return true
}
