package main

import (
	"encoding/json"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"tic-tac-toe/websocket"
)

func createDynamoTable(ctx *pulumi.Context, name string, attributes dynamodb.TableAttributeArray, isTtlEnabled bool) (*dynamodb.Table, error) {
	var ttl dynamodb.TableTtlPtrInput

	if isTtlEnabled {
		ttl = dynamodb.TableTtlArgs{
			AttributeName: pulumi.String("Ttl"),
			Enabled:       pulumi.Bool(isTtlEnabled),
		}
	}

	return dynamodb.NewTable(ctx, name, &dynamodb.TableArgs{
		Attributes:    attributes,
		BillingMode:   pulumi.String("PROVISIONED"),
		HashKey:       pulumi.String("Id"),
		ReadCapacity:  pulumi.Int(20),
		WriteCapacity: pulumi.Int(20),
		Ttl:           ttl,
	})
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		region, err := aws.GetRegion(ctx, nil, nil)
		if err != nil {
			return err
		}

		// An execution lambdaRole to use for the Lambda function
		lambdaPolicy, err := json.Marshal(map[string]interface{}{
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

		lambdaRole, err := iam.NewRole(ctx, "lambdaRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(lambdaPolicy),
			ManagedPolicyArns: pulumi.StringArray{
				iam.ManagedPolicyAWSLambdaBasicExecutionRole,
				iam.ManagedPolicyAmazonDynamoDBFullAccess,
				iam.ManagedPolicyAmazonAPIGatewayInvokeFullAccess,
			},
		})
		if err != nil {
			return err
		}

		connectionTable, err := createDynamoTable(ctx, "connection", dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Id"),
				Type: pulumi.String("S"),
			},
		}, true)
		if err != nil {
			return err
		}

		api, err := apigatewayv2.NewApi(ctx, "websocket-api", &apigatewayv2.ApiArgs{
			ProtocolType:             pulumi.String("WEBSOCKET"),
			RouteSelectionExpression: pulumi.String("$request.body.action"),
		})
		if err != nil {
			return err
		}

		connectLambdaProxy, err := websocket.NewLambdaProxy(ctx, "connect", websocket.LambdaProxyArgs{
			LambdaRole: lambdaRole,
			Api:        api,
			LambdaEnvironment: pulumi.StringMap{
				"CONNECTION_TABLE_NAME": connectionTable.Name,
			},
			RouteKey: "$connect",
		})
		if err != nil {
			return err
		}

		defaultLambdaProxy, err := websocket.NewLambdaProxy(ctx, "ws-fallback", websocket.LambdaProxyArgs{
			LambdaRole: lambdaRole,
			Api:        api,
			RouteKey:   "$default",
		})
		if err != nil {
			return err
		}

		disconnectLambdaProxy, err := websocket.NewLambdaProxy(ctx, "disconnect", websocket.LambdaProxyArgs{
			LambdaRole: lambdaRole,
			Api:        api,
			LambdaEnvironment: pulumi.StringMap{
				"CONNECTION_TABLE_NAME": connectionTable.Name,
			},
			RouteKey: "$disconnect",
		})
		if err != nil {
			return err
		}

		gameTable, err := createDynamoTable(ctx, "game", dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Id"),
				Type: pulumi.String("S"),
			},
		}, false)
		if err != nil {
			return err
		}

		gameEnvironment := pulumi.StringMap{
			"CONNECTION_TABLE_NAME": connectionTable.Name,
			"GAME_TABLE_NAME":       gameTable.Name,
			"REGION":                pulumi.String(region.Name),
		}

		createGameLambdaProxy, err := websocket.NewLambdaProxy(ctx, "create-game", websocket.LambdaProxyArgs{
			LambdaRole:        lambdaRole,
			Api:               api,
			LambdaEnvironment: gameEnvironment,
			RouteKey:          "create-game",
		})
		if err != nil {
			return err
		}

		getGameLambdaProxy, err := websocket.NewLambdaProxy(ctx, "get-game", websocket.LambdaProxyArgs{
			LambdaRole:        lambdaRole,
			Api:               api,
			LambdaEnvironment: gameEnvironment,
			RouteKey:          "get-game",
		})
		if err != nil {
			return err
		}

		makeMoveLambdaProxy, err := websocket.NewLambdaProxy(ctx, "make-move", websocket.LambdaProxyArgs{
			LambdaRole:        lambdaRole,
			Api:               api,
			LambdaEnvironment: gameEnvironment,
			RouteKey:          "make-move",
		})
		if err != nil {
			return err
		}

		apiStage, err := websocket.NewApiStage(ctx, "dev", websocket.ApiStageArgs{
			Api: api,
			LambdaProxies: []*websocket.LambdaProxy{
				connectLambdaProxy,
				defaultLambdaProxy,
				disconnectLambdaProxy,
				createGameLambdaProxy,
				getGameLambdaProxy,
				makeMoveLambdaProxy,
			},
		})

		ctx.Export("websocketUrl", apiStage.WebSocketUrl)

		return err
	})
}
