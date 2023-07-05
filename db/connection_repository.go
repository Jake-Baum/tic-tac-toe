package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"time"
)

type Connection struct {
	Id  string
	Ttl int64
}

var connectionTableName = os.Getenv("TABLE_NAME")

func CreateConnection(id string) (Connection, error) {

	ttlExpiryTime := generateUnixTimestampIn20Minutes()
	c := Connection{
		Id:  id,
		Ttl: ttlExpiryTime,
	}

	result, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return Connection{}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(connectionTableName),
		Item:      result,
	}
	if _, err := svc.PutItem(input); err != nil {
		return Connection{}, err
	}

	return c, nil
}

func GetConnection(id string) (Connection, error) {
	c := Connection{}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(connectionTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	if result, err := svc.GetItem(input); err != nil {
		return Connection{}, err

	} else if result.Item == nil {
		return Connection{}, &EntityDoesNotExistError{
			entityType: connectionTableName,
			id:         id,
		}

	} else if err := dynamodbattribute.UnmarshalMap(result.Item, &c); err != nil {
		return Connection{}, err

	} else {
		return c, nil
	}
}

func RefreshTtl(id string) (Connection, error) {
	result, err := dynamodbattribute.MarshalMap(Connection{
		Id:  id,
		Ttl: generateUnixTimestampIn20Minutes(),
	})
	if err != nil {
		return Connection{}, err
	}

	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":Ttl": result["Ttl"],
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(connectionTableName),
		ExpressionAttributeValues: expressionAttributeValues,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET Ttl = :Ttl"),
	}

	updatedConnection := Connection{}
	if result, err := svc.UpdateItem(input); err != nil {
		return Connection{}, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Attributes, &updatedConnection); err != nil {
		return Connection{}, err
	} else {
		return updatedConnection, nil
	}
}

func generateUnixTimestampInXMinutes(x int) int64 {
	return time.Now().Add(time.Duration(x) * time.Minute).Unix()
}

func generateUnixTimestampIn20Minutes() int64 {
	return generateUnixTimestampInXMinutes(20)
}
