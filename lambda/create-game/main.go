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

type CreateGameRequest struct {
	PlayerO string `json:"playerO"`
}

func createGame(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID

	var requestBody CreateGameRequest
	err := json.Unmarshal([]byte(websocketEvent.Body), &requestBody)
	if err != nil {
		log.Errorf("An error occurred while deserialising request body %s - %s", websocketEvent.Body, err)
		return utils.InternalServerErrorResponse(), nil
	}

	playerOConnection, err := db.GetConnection(requestBody.PlayerO)
	if err != nil {
		switch err.(type) {
		case *db.EntityDoesNotExistError:
			log.Info(err)
			return utils.NotFoundResponse(err), nil
		default:
			log.Errorf("An error occurred while retrieving connection with ID %s - %s", requestBody.PlayerO, err)
			return utils.InternalServerErrorResponse(), nil
		}
	}

	g := *game.NewGame(connectionId, playerOConnection.Id)

	g, err = db.CreateGame(g)
	if err != nil {
		log.Errorf("An error occurred when creating game - %s", err)
		return utils.InternalServerErrorResponse(), nil
	}

	if err = utils.SendMessage(websocketEvent, playerOConnection.Id, g); err != nil {
		log.Errorf("An error occurred when sending message from %s to %s - %s", connectionId, playerOConnection.Id, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(g), nil
}

func main() {
	utils.Initialize()
	lambda.Start(createGame)
}
