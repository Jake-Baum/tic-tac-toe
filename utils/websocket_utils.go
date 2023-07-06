package utils

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	log "github.com/sirupsen/logrus"
	"os"
)

type Request struct {
	Action    string `json:"action"`
	MessageTo string `json:"messageTo"`
}

var region = os.Getenv("REGION")

func newApiGatewaySession(websocketEvent events.APIGatewayWebsocketProxyRequest) (*apigatewaymanagementapi.ApiGatewayManagementApi, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(fmt.Sprintf("https://%s/%s", websocketEvent.RequestContext.DomainName, websocketEvent.RequestContext.Stage)),
	})

	if err != nil {
		return nil, err
	}

	return apigatewaymanagementapi.New(sess), nil
}

func SendMessage(websocketEvent events.APIGatewayWebsocketProxyRequest, messageTo string, message interface{}) error {
	responseBodySerialized, err := json.MarshalIndent(message, "", "")
	if err != nil {
		return err
	}

	apiGatewaySession, err := newApiGatewaySession(websocketEvent)
	if err != nil {
		log.Errorf("An error occurred when creating API Gateway session - %s", err)
		return err
	}

	connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(messageTo),
		Data:         responseBodySerialized,
	}
	_, err = apiGatewaySession.PostToConnection(connectionInput)
	return err
}
