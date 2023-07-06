package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

func wsFallback(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody utils.Request
	err := json.Unmarshal([]byte(websocketEvent.Body), &requestBody)
	if err != nil {
		log.Errorf("An error occurred while deserialising request body %s - %s", websocketEvent.Body, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.BadRequestResponse(fmt.Sprintf("Action (%s) does not correspond to any routes", requestBody.Action)), nil
}

func main() {
	utils.Initialize()
	lambda.Start(wsFallback)
}
