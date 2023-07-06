package utils

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"os"
)

type Request struct {
	Action    string `json:"action"`
	MessageTo string `json:"messageTo"`
}

var region = os.Getenv("REGION")

func NewApiGatewaySession(websocketEvent events.APIGatewayWebsocketProxyRequest) (*apigatewaymanagementapi.ApiGatewayManagementApi, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(fmt.Sprintf("https://%s/%s", websocketEvent.RequestContext.DomainName, websocketEvent.RequestContext.Stage)),
	})

	if err != nil {
		return nil, err
	}

	return apigatewaymanagementapi.New(sess), nil
}
