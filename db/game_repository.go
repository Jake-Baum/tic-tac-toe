package db

import (
	"github.com/Jake-Baum/tic-tac-toe/game"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
)

var tableName = os.Getenv("TABLE_NAME")

var svc = createClient()

func CreateGame(g game.Game) (game.Game, error) {
	id := uuid.New().String()
	g.Id = id

	result, err := dynamodbattribute.MarshalMap(g)
	if err != nil {
		return game.Game{}, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      result,
	}
	if _, err := svc.PutItem(input); err != nil {
		return game.Game{}, err
	}

	return g, nil
}

func GetGame(id string) (game.Game, error) {
	g := game.Game{}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}
	if result, err := svc.GetItem(input); err != nil {
		return game.Game{}, err

	} else if result.Item == nil {
		return game.Game{}, &EntityDoesNotExistError{
			entityType: "game",
			id:         id,
		}

	} else if err := dynamodbattribute.UnmarshalMap(result.Item, &g); err != nil {
		return game.Game{}, err

	} else {
		return g, nil
	}
}

func UpdateGame(g game.Game) (game.Game, error) {
	result, err := dynamodbattribute.MarshalMap(g)
	if err != nil {
		return game.Game{}, err
	}

	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":Board":       result["Board"],
		":CurrentTurn": result["CurrentTurn"],
	}
	log.Info(expressionAttributeValues)

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeValues: expressionAttributeValues,
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(g.Id),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET Board = :Board, CurrentTurn = :CurrentTurn"),
	}

	updatedGame := game.Game{}
	if result, err := svc.UpdateItem(input); err != nil {
		return game.Game{}, err
	} else if err := dynamodbattribute.UnmarshalMap(result.Attributes, &updatedGame); err != nil {
		return game.Game{}, err
	} else {
		return updatedGame, nil
	}
}
