package main

import (
	"context"
	"encoding/json"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/game"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type MakeMoveRequest struct {
	Id   string `json:"id"`
	Move int    `json:"move"`
}

func makeMove(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID

	var requestBody MakeMoveRequest
	err := json.Unmarshal([]byte(websocketEvent.Body), &requestBody)
	if err != nil {
		log.Errorf("An error occurred while deserialising request body %s - %s", websocketEvent.Body, err)
		return utils.InternalServerErrorResponse(), nil
	}
	gameId := requestBody.Id
	move := requestBody.Move

	g, err := db.GetGame(gameId)
	if err != nil {
		switch err.(type) {
		case *db.EntityDoesNotExistError:
			log.Info(err)
			return utils.NotFoundResponse(err), nil
		default:
			log.Errorf("An error occurred while retrieving game with ID %s - %s", gameId, err)
			return utils.InternalServerErrorResponse(), nil
		}
	}

	err = g.MakeMove(connectionId, move)
	if err != nil {
		switch err.(type) {
		case *game.FinishedError:
			return utils.ConflictResponse(err.Error()), nil
		case *game.InvalidMoveError:
			return utils.BadRequestResponse(err.Error()), nil
		case *game.NotPlayersTurnError:
			return utils.BadRequestResponse(err.Error()), nil
		case *game.PlayerDoesNotExistError:
			return utils.ForbiddenResponse(), nil
		}
	}

	updatedGame, err := db.UpdateGame(g)
	if err != nil {
		switch err.(type) {
		case *db.EntityDoesNotExistError:
			log.Info(err)
			return utils.NotFoundResponse(err), nil
		default:
			log.Errorf("An error occurred while retrieving game with ID %s - %s", gameId, err)
			return utils.InternalServerErrorResponse(), nil
		}
	}

	otherPlayerId := g.GetOtherPlayer(connectionId)
	if err = utils.SendMessage(websocketEvent, otherPlayerId, g); err != nil {
		log.Errorf("An error occurred when sending message from %s to %s - %s", connectionId, otherPlayerId, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(updatedGame), nil

}

func main() {
	utils.Initialize()
	lambda.Start(makeMove)
}
