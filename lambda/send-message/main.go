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

func sendMessage(_ context.Context, websocketEvent events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionId := websocketEvent.RequestContext.ConnectionID

	var requestBody utils.Request
	err := json.Unmarshal([]byte(websocketEvent.Body), &requestBody)
	if err != nil {
		log.Errorf("An error occurred while deserialising request body %s - %s", websocketEvent.Body, err)
		return utils.InternalServerErrorResponse(), nil
	}
	messageTo := requestBody.MessageTo

	_, err = db.GetConnection(messageTo)
	if err != nil {
		switch err.(type) {
		case *db.EntityDoesNotExistError:
			log.Info(err)
			return utils.NotFoundResponse(err), nil
		default:
			log.Errorf("An error occurred while retrieving connection with ID %s - %s", messageTo, err)
			return utils.InternalServerErrorResponse(), nil
		}
	}

	message := utils.MessageResponseJson(fmt.Sprintf("Hi %s!  From %s", messageTo, connectionId))
	if err = utils.SendMessage(websocketEvent, messageTo, message); err != nil {
		log.Errorf("An error occurred when sending message from %s to %s - %s", connectionId, messageTo, err)
		return utils.InternalServerErrorResponse(), nil
	}

	return utils.OkResponse(utils.MessageResponseJson(fmt.Sprintf("Message sent to %s", messageTo))), nil
}

func main() {
	utils.Initialize()
	lambda.Start(sendMessage)
}
