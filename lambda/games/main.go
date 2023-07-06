package main

import (
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/game"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func createGame() events.APIGatewayProxyResponse {
	g := *game.NewGame()

	if g, err := db.CreateGame(g); err == nil {
		return utils.OkResponse(g)
	} else {
		log.Error(err)
		return utils.InternalServerErrorResponse()
	}
}

func handleGamesHandlerRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodPost:
		return createGame(), nil
	default:
		return utils.MethodNotAllowedResponse(), nil
	}
}

func main() {
	utils.Initialize()
	lambda.Start(handleGamesHandlerRequest)
}
