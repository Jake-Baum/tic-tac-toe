package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	. "jakebaum.uk/tictactoe/game"
	"net/http"
	"unicode"
)

var games = make(map[int]Game)

func StartGame(c *gin.Context) {
	id := len(games)
	games[id] = *NewGame(id)
	c.IndentedJSON(http.StatusOK, games[id])
}

func executeWithGame(c *gin.Context, f func(game Game)) {
	if gameId, err := IntPathParam(c, "gameId"); err != nil {
		c.IndentedJSON(http.StatusBadRequest, CreateMessageResponse("Game ID is not valid"))
	} else {
		game, exists := games[gameId]
		if !exists {
			c.IndentedJSON(http.StatusNotFound, CreateMessageResponse("Game with ID %d does not exist", gameId))
			return
		}

		f(game)
	}
}

func GetGame(c *gin.Context) {
	executeWithGame(c, func(game Game) {
		c.IndentedJSON(http.StatusOK, game)
	})
}

func MakeMove(c *gin.Context) {
	executeWithGame(c, func(game Game) {

		if game.IsDraw() {
			message := "The game has ended in a draw"
			log.Info(message)
			c.IndentedJSON(http.StatusConflict, CreateMessageResponse(message))
			return
		}

		if isWinner, winner := game.IsWinner(); isWinner {
			message := fmt.Sprintf("The game has ended.  The winner is %s!", winner)
			log.Info(message)
			c.IndentedJSON(http.StatusOK, CreateMessageResponse(message))
			return
		}

		input := c.Query("move")

		square, err := getSquare(input)
		if err != nil {
			log.Infof(err.Error())
			c.IndentedJSON(http.StatusBadRequest, CreateMessageResponse(err.Error()))
			return
		}

		if !game.IsValidMove(square) {
			message := fmt.Sprintf("%d is not a valid move", square)
			log.Info(message)
			c.IndentedJSON(http.StatusBadRequest, CreateMessageResponse(message))
			return
		}

		game.Board[square/3][square%3] = game.CurrentTurn

		if game.CurrentTurn == X {
			game.CurrentTurn = O
		} else {
			game.CurrentTurn = X
		}

		c.IndentedJSON(http.StatusOK, game.Board)

	})
}

func getSquare(s string) (int, error) {
	if len(s) <= 0 || !unicode.IsDigit(rune(s[0])) {
		return -1, parseError{fmt.Sprintf("%s is not a valid move selection", s)}
	}

	return int(s[0]) - 48, nil
}
