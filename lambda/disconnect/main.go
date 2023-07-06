package main

import (
	"context"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

func disconnect(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID
	message := fmt.Sprintf("Goodbye %s", connectionId)

	if _, err := db.DeleteConnection(connectionId); err != nil {
		switch err.(type) {
		case *db.EntityDoesNotExistError:
			log.Info(err)
			return utils.NotFoundResponse(err), nil
		default:
			log.Errorf("An error occurred while deleting connection with ID %s - %s", connectionId, err)
			return utils.InternalServerErrorResponse(), nil
		}
	}

	return utils.OkResponse(utils.MessageResponseJson(message)), nil
}

func main() {
	utils.Initialize()
	lambda.Start(disconnect)
}
