package main

import (
	"errors"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/game"
	. "github.com/Jake-Baum/tic-tac-toe/lambda"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func makeMove(id string, move string) events.APIGatewayProxyResponse {
	g, err := db.GetGame(id)
	if err != nil {
		log.Error(err)

		var entityDoesNotExistError *db.EntityDoesNotExistError
		switch {
		case errors.As(err, &entityDoesNotExistError):
			return utils.NotFoundResponse(*entityDoesNotExistError)
		default:
			return utils.InternalServerErrorResponse()
		}
	}

	if g.IsDraw() {
		return utils.ConflictResponse("The game has ended in a draw")
	}

	if isWinner, winner := g.IsWinner(); isWinner {
		return utils.ConflictResponse(fmt.Sprintf("The game has ended.  The winner is %s!", winner))
	}

	if square, err := game.GetSquare(move); err != nil {
		return utils.BadRequestResponse(err.Error())

	} else if !g.IsValidMove(square) {
		return utils.BadRequestResponse(fmt.Sprintf("%d is not a valid move", square))
	} else {
		g.Board[square/3][square%3] = g.CurrentTurn

		if g.CurrentTurn == game.X {
			g.CurrentTurn = game.O
		} else {
			g.CurrentTurn = game.X
		}

		if updatedGame, err := db.UpdateGame(g); err != nil {
			log.Error(err)
			return utils.InternalServerErrorResponse()
		} else {
			return utils.OkResponse(updatedGame)
		}
	}
}

func getGame(id string) events.APIGatewayProxyResponse {
	g, err := db.GetGame(id)
	if err != nil {
		log.Error(err)

		var entityDoesNotExistError *db.EntityDoesNotExistError
		switch {
		case errors.As(err, &entityDoesNotExistError):
			return utils.NotFoundResponse(*entityDoesNotExistError)
		default:
			return utils.InternalServerErrorResponse()
		}
	}

	return utils.OkResponse(g)
}

func handleGameHandlerRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	gameId := request.PathParameters["gameId"]

	switch request.HTTPMethod {
	case http.MethodGet:
		return getGame(gameId), nil
	case http.MethodPut:
		move := request.QueryStringParameters["move"]
		return makeMove(gameId, move), nil
	default:
		return utils.MethodNotAllowedResponse(), nil
	}

}

func main() {
	Initialize()
	lambda.Start(handleGameHandlerRequest)
}
