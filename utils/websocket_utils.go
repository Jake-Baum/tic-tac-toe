package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"os"
)

type Request struct {
	Action    string `json:"action"`
	MessageTo string `json:"messageTo"`
}

var apiGatewayEndpoint = os.Getenv("API_GATEWAY_ENDPOINT")
var region = os.Getenv("REGION")

func NewApiGatewaySession() (*apigatewaymanagementapi.ApiGatewayManagementApi, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(apiGatewayEndpoint),
	})

	if err != nil {
		return nil, err
	}

	return apigatewaymanagementapi.New(sess), nil
}
