package utils

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func MessageResponseJson(message string) string {
	return fmt.Sprintf("{\"message\": \"%s\"}", message)
}

var defaultHeaders = map[string]string{"Content-Type": "application/json"}

func OkResponse(responseBody interface{}) events.APIGatewayProxyResponse {
	if responseBodySerialized, err := json.MarshalIndent(responseBody, "", ""); err == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    defaultHeaders,
			Body:       string(responseBodySerialized),
		}
	} else {
		log.Errorf("An error occurred when serialising response (%s) to JSON - %s", responseBody, err)
		return InternalServerErrorResponse()
	}
}

func BadRequestResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    defaultHeaders,
		Body:       MessageResponseJson(message),
	}
}

func ForbiddenResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusForbidden,
		Headers:    defaultHeaders,
		Body:       MessageResponseJson("You are not allowed to access this resource"),
	}
}

func NotFoundResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers:    defaultHeaders,
		Body:       MessageResponseJson(err.Error()),
	}
}

func MethodNotAllowedResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Headers:    defaultHeaders,
	}
}

func ConflictResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusConflict,
		Headers:    defaultHeaders,
		Body:       MessageResponseJson(message),
	}
}

func InternalServerErrorResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       MessageResponseJson("An unexpected error occurred"),
		Headers:    defaultHeaders,
	}
}
