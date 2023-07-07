package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func defaultGame() Game {
	return Game{
		Board: [][]Piece{},
	}
}

func TestGame_IsWinner_EmptyBoard(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
	}

	var isWinner, winner = game.isWinner()

	assert.Equal(t, isWinner, false)
	assert.Equal(t, winner, EMPTY)
}

func TestGame_IsWinner_Row(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, X},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
	}

	var isWinner, winner = game.isWinner()

	assert.Equal(t, isWinner, true)
	assert.Equal(t, winner, X)
}

func TestGame_IsWinner_Column(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, EMPTY, O},
			{EMPTY, X, O},
		},
	}

	var isWinner, winner = game.isWinner()

	assert.Equal(t, isWinner, true)
	assert.Equal(t, winner, O)
}

func TestGame_IsWinner_Diag1(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, X, O},
			{EMPTY, O, X},
		},
	}

	var isWinner, winner = game.isWinner()

	assert.Equal(t, isWinner, true)
	assert.Equal(t, winner, X)
}

func TestGame_IsWinner_Diag2(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	var isWinner, winner = game.isWinner()

	assert.Equal(t, isWinner, true)
	assert.Equal(t, winner, O)
}

func TestGame_IsValidMove_Success(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.isValidMove(3)

	assert.Equal(t, isValidMove, true)
}

func TestGame_IsValidMove_TooBig(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.isValidMove(9)

	assert.Equal(t, isValidMove, false)

}

func TestGame_IsValidMove_TooSmall(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.isValidMove(-1)

	assert.Equal(t, isValidMove, false)
}

func TestGame_IsValidMove_SquareAlreadyOccupied(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{EMPTY, O, O},
			{O, X, X},
		},
	}

	isValidMove := game.isValidMove(0)

	assert.Equal(t, isValidMove, false)
}

func TestGame_IsBoardFull_BoardFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{O, X, X},
			{X, O, O},
		},
	}

	isDraw := game.isDraw()

	assert.Equal(t, isDraw, true)
}

func TestGame_IsBoardFull_BoardNotFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{O, X, X},
			{X, O, EMPTY},
		},
	}

	isDraw := game.isDraw()

	assert.Equal(t, isDraw, false)
}

func TestGame_IsDraw_BoardFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{O, X, X},
			{X, O, O},
		},
	}

	isDraw := game.isDraw()

	assert.Equal(t, isDraw, true)
}

func TestGame_IsDraw_BoardNotFull(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, O},
			{X, EMPTY, O},
			{O, X, X},
		},
	}

	isDraw := game.isDraw()

	assert.Equal(t, isDraw, false)
}

func TestGame_IsDraw_IsWin(t *testing.T) {
	var game = Game{
		Board: [][]Piece{
			{X, X, X},
			{O, EMPTY, O},
			{O, X, EMPTY},
		},
	}

	isDraw := game.isDraw()

	assert.Equal(t, isDraw, false)
}

func TestGame_IsPlayer_PlayerX(t *testing.T) {
	game := NewGame("playerX", "playerO")

	assert.Equal(t, true, game.IsPlayer("playerX"))
}

func TestGame_IsPlayer_PlayerO(t *testing.T) {
	game := NewGame("playerX", "playerO")

	assert.Equal(t, true, game.IsPlayer("playerO"))
}

func TestGame_IsPlayer_NotAPlayer(t *testing.T) {
	game := NewGame("playerX", "playerO")

	assert.Equal(t, false, game.IsPlayer("someOtherPlayer"))
}

func TestGame_GetOtherPlayer_PlayerX(t *testing.T) {
	game := NewGame("playerX", "playerO")

	assert.Equal(t, "playerO", game.GetOtherPlayer("playerX"))
}

func TestGame_GetOtherPlayer_PlayerO(t *testing.T) {
	game := NewGame("playerX", "playerO")

	assert.Equal(t, "playerX", game.GetOtherPlayer("playerO"))
}

func TestGame_GetOtherPlayer_NotAPlayer(t *testing.T) {
	game := NewGame("playerX", "playerO")

	assert.Equal(t, "playerX", game.GetOtherPlayer("someOtherPlayer"))
}

func TestGame_MakeMove_NotAPlayer(t *testing.T) {
	game := NewGame("playerX", "playerO")

	err := game.MakeMove("someOtherPlayer", 0)
	assert.Equal(t, &PlayerDoesNotExistError{
		player: "someOtherPlayer",
		gameId: "",
	}, err)

	expectedGame := NewGame("playerX", "playerO")
	assert.Equal(t, expectedGame, game)
}

func TestGame_MakeMove_IsDraw(t *testing.T) {
	game := NewGame("playerX", "playerO")
	game.Board = [][]Piece{
		{X, O, X},
		{X, O, O},
		{O, X, X},
	}

	err := game.MakeMove("playerX", 0)
	assert.Equal(t, &FinishedError{}, err)
}

func TestGame_MakeMove_XIsWinner(t *testing.T) {
	game := NewGame("playerX", "playerO")
	game.Board = [][]Piece{
		{X, O, X},
		{X, O, O},
		{X, EMPTY, EMPTY},
	}
	game.CurrentTurn = O

	err := game.MakeMove("playerO", 8)

	expectedWinner := X
	assert.Equal(t, &FinishedError{winner: &expectedWinner}, err)
}

func TestGame_MakeMove_InvalidMove(t *testing.T) {
	game := NewGame("playerX", "playerO")
	game.Board = [][]Piece{
		{X, O, X},
		{X, O, O},
		{EMPTY, EMPTY, EMPTY},
	}

	err := game.MakeMove("playerX", 1)

	assert.Equal(t, &InvalidMoveError{move: 1}, err)
}

func TestGame_MakeMove_NotPlayersTurn(t *testing.T) {
	game := NewGame("playerX", "playerO")
	game.Board = [][]Piece{
		{X, O, X},
		{X, O, O},
		{EMPTY, EMPTY, EMPTY},
	}

	err := game.MakeMove("playerO", 7)

	assert.Equal(t, &NotPlayersTurnError{player: "playerO"}, err)
}

func TestGame_MakeMove_ValidMove(t *testing.T) {
	game := NewGame("playerX", "playerO")
	game.Board = [][]Piece{
		{X, O, X},
		{X, O, O},
		{EMPTY, EMPTY, EMPTY},
	}

	err := game.MakeMove("playerX", 6)

	expectedGame := NewGame("playerX", "playerO")
	expectedGame.Board = [][]Piece{
		{X, O, X},
		{X, O, O},
		{X, EMPTY, EMPTY},
	}
	expectedGame.CurrentTurn = O
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedGame, game)
}
