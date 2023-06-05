package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Jake-Baum/tic-tac-toe/db"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

func messageResponseJson(message string) string {
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
		return InternalServerErrorResponse()
	}
}

func BadRequestResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    defaultHeaders,
		Body:       messageResponseJson(message),
	}
}

func NotFoundResponse(err db.EntityDoesNotExistError) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers:    defaultHeaders,
		Body:       messageResponseJson(err.Error()),
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
		Body:       messageResponseJson(message),
	}
}

func InternalServerErrorResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       messageResponseJson("An unexpected error occurred"),
		Headers:    defaultHeaders,
	}
}
