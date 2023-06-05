package main

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-aws-apigateway/sdk/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const BASE_PATH = "/api"

func createLambda(ctx *pulumi.Context, name string, role *iam.Role) (*lambda.Function, error) {
	return lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Runtime: pulumi.String("go1.x"),
		Handler: pulumi.String(name),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive(fmt.Sprintf("../bin/handlers/%s.zip", name)),
	})
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// An execution role to use for the Lambda function
		policy, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "lambda.amazonaws.com",
					},
				},
			},
		})
		if err != nil {
			return err
		}

		database, err := dynamodb.NewTable(ctx, "game", &dynamodb.TableArgs{
			Attributes: dynamodb.TableAttributeArray{
				&dynamodb.TableAttributeArgs{
					Name: pulumi.String("Id"),
					Type: pulumi.String("S"),
				},
			},
			BillingMode:   pulumi.String("PROVISIONED"),
			HashKey:       pulumi.String("Id"),
			ReadCapacity:  pulumi.Int(20),
			WriteCapacity: pulumi.Int(20),
		})
		if err != nil {
			return err
		}
		ctx.Export("Database name", database.Name)

		role, err := iam.NewRole(ctx, "role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(policy),
			ManagedPolicyArns: pulumi.StringArray{
				iam.ManagedPolicyAWSLambdaBasicExecutionRole,
				iam.ManagedPolicyAmazonDynamoDBFullAccess,
			},
		})
		if err != nil {
			return err
		}

		gamesHandlerFunction, err := createLambda(ctx, "games_handler", role)
		if err != nil {
			return err
		}

		gameHandlerFunction, err := createLambda(ctx, "game_handler", role)
		if err != nil {
			return err
		}

		// A REST API to route requests to HTML content and the Lambda function
		method := apigateway.MethodANY
		api, err := apigateway.NewRestAPI(ctx, "api", &apigateway.RestAPIArgs{
			Routes: []apigateway.RouteArgs{
				{Path: BASE_PATH + "/game", Method: &method, EventHandler: gamesHandlerFunction},
				{Path: BASE_PATH + "/game/{gameId}", Method: &method, EventHandler: gameHandlerFunction},
			},
		})
		if err != nil {
			return err
		}

		// The URL at which the REST API will be served
		ctx.Export("url", api.Url)
		return nil

	})
}
