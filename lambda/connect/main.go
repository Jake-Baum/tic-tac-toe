package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type WebSocketRequestBody struct {
	Action string
	Data   json.RawMessage
}

func connect(ctx context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID
	callbackUrl := fmt.Sprintf("wss://%s/%s", websocketEvent.RequestContext.DomainName, websocketEvent.RequestContext.Stage)

	if _, err := db.CreateConnection(connectionId); err != nil {
		log.Errorf("An error occurred while creating connection with ID %s - %s", connectionId, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(utils.MessageResponseJson(callbackUrl)), nil
}

func main() {
	utils.Initialize()
	lambda.Start(connect)
}
