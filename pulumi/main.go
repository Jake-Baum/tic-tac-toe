package main

import (
	"encoding/json"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"tic-tac-toe/websocket"
	"time"
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

func createApiGatewayWebsocket(ctx *pulumi.Context, name string) (*apigatewayv2.Api, error) {
	websocket, err := apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		ProtocolType:             pulumi.String("WEBSOCKET"),
		RouteSelectionExpression: pulumi.String("$request.body.action"),
	})
	if err != nil {
		return nil, err
	}

	return websocket, nil
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

		apiGatewayPolicy, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "apigateway.amazonaws.com",
					},
				},
			},
		})
		apiGatewayRole, err := iam.NewRole(ctx, "apiGatewayRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(apiGatewayPolicy),
			ManagedPolicyArns: pulumi.StringArray{
				iam.ManagedPolicyAmazonAPIGatewayPushToCloudWatchLogs,
			},
		})

		_, err = apigateway.NewAccount(ctx, "apiGatewayAccount", &apigateway.AccountArgs{
			CloudwatchRoleArn: apiGatewayRole.Arn,
		})

		apiGatewayLogGroup, err := cloudwatch.NewLogGroup(ctx, "apiGatewayLogGroup", &cloudwatch.LogGroupArgs{})
		if err != nil {
			return nil
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

		api, err := createApiGatewayWebsocket(ctx, "websocket-api")

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

		websocketDeployment, err := apigatewayv2.NewDeployment(ctx, "websocketDeployment", &apigatewayv2.DeploymentArgs{
			ApiId: api.ID(),
			Triggers: pulumi.StringMap{
				"deployedAt": pulumi.String(time.Now().Format(time.RFC3339)), // This is somewhat of a hack to force the API to redeploy on changes.  Must be a better way
			},
		}, pulumi.DependsOn([]pulumi.Resource{
			connectLambdaProxy,
			defaultLambdaProxy,
			disconnectLambdaProxy,
			createGameLambdaProxy,
			getGameLambdaProxy,
			makeMoveLambdaProxy}))
		if err != nil {
			return err
		}

		stage, err := apigatewayv2.NewStage(ctx, "dev", &apigatewayv2.StageArgs{
			ApiId:        api.ID(),
			DeploymentId: websocketDeployment.ID(),
			AccessLogSettings: apigatewayv2.StageAccessLogSettingsArgs{
				DestinationArn: apiGatewayLogGroup.Arn,
				Format:         pulumi.String("$context.apiId, $context.requestId, $context.authorize.error, $context.authorize.status, $context.authorizer.error, $context.authorizer.integrationStatus, $context.authorizer.requestId, $context.authorizer.status"),
			},
			DefaultRouteSettings: apigatewayv2.StageDefaultRouteSettingsArgs{
				DataTraceEnabled:       pulumi.Bool(true),
				DetailedMetricsEnabled: pulumi.Bool(true),
				LoggingLevel:           pulumi.String("INFO"),
				ThrottlingBurstLimit:   pulumi.Int(100),
				ThrottlingRateLimit:    pulumi.Float64(100),
			},
		})
		if err != nil {
			return err
		}

		// The URL at which the REST API will be served
		//ctx.Export("rest_url", api.Url)
		ctx.Export("websocket_url", stage.InvokeUrl)
		return nil

	})
}
