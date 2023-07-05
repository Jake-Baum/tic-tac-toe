package main

import (
	"context"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/db"
	. "github.com/Jake-Baum/tic-tac-toe/lambda"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

func sendMessage(ctx context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID

	connection, err := db.GetConnection(connectionId)
	if err != nil {
		log.Errorf("An error occurred while retrieving connection with ID %s: %s", connectionId, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(utils.MessageResponseJson(fmt.Sprintf("Hi %s", connection.Id))), nil
}

func main() {
	Initialize()
	lambda.Start(sendMessage)
}
