package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/Jake-Baum/tic-tac-toe/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	log "github.com/sirupsen/logrus"
)

func sendMessage(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID

	var requestBody utils.Request
	err := json.Unmarshal([]byte(websocketEvent.Body), &requestBody)
	if err != nil {
		log.Errorf("An error occurred while deserialising request body %s - %s", websocketEvent.Body, err)
		return utils.InternalServerErrorResponse(), nil
	}

	_, err = db.GetConnection(requestBody.MessageTo)
	if err != nil {
		switch err.(type) {
		case *db.EntityDoesNotExistError:
			log.Info(err)
			return utils.NotFoundResponse(err), nil
		default:
			log.Errorf("An error occurred while retrieving connection with ID %s - %s", requestBody.MessageTo, err)
			return utils.InternalServerErrorResponse(), nil
		}
	}

	apiGatewaySession, err := utils.NewApiGatewaySession(websocketEvent)
	if err != nil {
		log.Errorf("An error occurred when creating API Gateway session - %s", err)
		return utils.InternalServerErrorResponse(), nil
	}

	connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(requestBody.MessageTo),
		Data:         []byte(utils.MessageResponseJson(fmt.Sprintf("To %s.  From %s", requestBody.MessageTo, connectionId))),
	}
	_, err = apiGatewaySession.PostToConnection(connectionInput)
	if err != nil {
		log.Errorf("An error occurred when sending message from %s to %s - %s", connectionId, requestBody.MessageTo, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(utils.MessageResponseJson(fmt.Sprintf("Message sent to %s", requestBody.MessageTo))), nil
}

func main() {
	utils.Initialize()
	lambda.Start(sendMessage)
}
