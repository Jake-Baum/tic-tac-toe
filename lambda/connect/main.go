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

func connect(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID
	message := fmt.Sprintf("Connection ID: %s", connectionId)

	if _, err := db.CreateConnection(connectionId); err != nil {
		log.Errorf("An error occurred while creating connection with ID %s - %s", connectionId, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(utils.MessageResponseJson(message)), nil
}

func main() {
	utils.Initialize()
	lambda.Start(connect)
}
