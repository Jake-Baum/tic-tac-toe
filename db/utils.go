package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var svc = createClient()

type EntityDoesNotExistError struct {
	entityType string
	id         string
}

func (e *EntityDoesNotExistError) Error() string {
	return fmt.Sprintf("%s with ID %s does not exist", e.entityType, e.id)
}

func createClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}
