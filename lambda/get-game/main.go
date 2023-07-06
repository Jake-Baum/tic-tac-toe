package main

import (
	"context"
	"encoding/json"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type GetGameRequest struct {
	Id string `json:"id"`
}

func getGame(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID

	var requestBody GetGameRequest
	err := json.Unmarshal([]byte(websocketEvent.Body), &requestBody)
	if err != nil {
		log.Errorf("An error occurred while deserialising request body %s - %s", websocketEvent.Body, err)
		return utils.InternalServerErrorResponse(), nil
	}
	gameId := requestBody.Id

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

	if !g.IsPlayer(connectionId) {
		log.Infof("User with ID %s tried to retrieve game %s which does not belong to them", connectionId, gameId)
		return utils.ForbiddenResponse(), nil
	}

	return utils.OkResponse(g), nil
}

func main() {
	utils.Initialize()
	lambda.Start(getGame)
}
