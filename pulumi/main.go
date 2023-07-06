package main

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"

	//"github.com/pulumi/pulumi-aws-apigateway/sdk/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const basePath = "/api"

func createLambda(ctx *pulumi.Context, name string, role *iam.Role, environment pulumi.StringMap) (*lambda.Function, error) {
	return lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Runtime: pulumi.String("go1.x"),
		Handler: pulumi.String("main"),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive(fmt.Sprintf("../bin/lambda/%s/main.zip", name)),
		Environment: lambda.FunctionEnvironmentArgs{
			Variables: environment,
		},
	})
}

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

		//gameTable, err := createDynamoTable(ctx, "game", dynamodb.TableAttributeArray{
		//	&dynamodb.TableAttributeArgs{
		//		Name: pulumi.String("Id"),
		//		Type: pulumi.String("S"),
		//	},
		//}, false)
		//if err != nil {
		//	return err
		//}

		//gamesHandlerFunction, err := createLambda(ctx, "games", lambdaRole, pulumi.StringMap{"TABLE_NAME": gameTable.Name})
		//if err != nil {
		//	return err
		//}
		//
		//gameHandlerFunction, err := createLambda(ctx, "game", lambdaRole, pulumi.StringMap{"TABLE_NAME": gameTable.Name})
		//if err != nil {
		//	return err
		//}

		//// A REST API to route requests to HTML content and the Lambda function
		//apigateway.New
		//method := apigateway.MethodANY
		//api, err := apigateway.NewRestAPI(ctx, "api", &apigateway.RestAPIArgs{
		//	Routes: []apigateway.RouteArgs{
		//		{Path: basePath + "/game", Method: &method, EventHandler: gamesHandlerFunction},
		//		{Path: basePath + "/game/{gameId}", Method: &method, EventHandler: gameHandlerFunction},
		//	},
		//})
		//if err != nil {
		//	return err
		//}

		connectionTable, err := createDynamoTable(ctx, "connection", dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Id"),
				Type: pulumi.String("S"),
			},
		}, true)
		if err != nil {
			return err
		}

		websocket, err := apigatewayv2.NewApi(ctx, "connect", &apigatewayv2.ApiArgs{
			ProtocolType:             pulumi.String("WEBSOCKET"),
			RouteSelectionExpression: pulumi.String("$request.body.action"),
		})
		if err != nil {
			return err
		}

		websocketDeployment, err := apigatewayv2.NewDeployment(ctx, "websocketDeployment", &apigatewayv2.DeploymentArgs{
			ApiId: websocket.ID(),
		})
		if err != nil {
			return err
		}

		stage, err := apigatewayv2.NewStage(ctx, "dev", &apigatewayv2.StageArgs{
			ApiId:        websocket.ID(),
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

		connectFunction, err := createLambda(ctx, "connect", lambdaRole, pulumi.StringMap{"TABLE_NAME": connectionTable.Name})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "connectInvokePermission", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  connectFunction.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: websocket.ExecutionArn.ApplyT(func(executionArn string) (string, error) {
				return fmt.Sprintf("%v/*", executionArn), nil
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}

		connectIntegration, err := apigatewayv2.NewIntegration(ctx, "connectIntegration", &apigatewayv2.IntegrationArgs{
			ApiId:                   websocket.ID(),
			IntegrationType:         pulumi.String("AWS_PROXY"),
			ConnectionType:          pulumi.String("INTERNET"),
			ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
			IntegrationMethod:       pulumi.String("POST"),
			IntegrationUri:          connectFunction.InvokeArn,
			PassthroughBehavior:     pulumi.String("WHEN_NO_MATCH"),
		})
		if err != nil {
			return err
		}

		connectRoute, err := apigatewayv2.NewRoute(ctx, "connectRoute", &apigatewayv2.RouteArgs{
			ApiId:                            websocket.ID(),
			RouteKey:                         pulumi.String("$connect"),
			AuthorizationType:                pulumi.String("NONE"),
			RouteResponseSelectionExpression: pulumi.String("$default"),
			Target: connectIntegration.ID().ApplyT(func(id string) (string, error) {
				return fmt.Sprintf("integrations/%v", id), nil
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}

		_, err = apigatewayv2.NewRouteResponse(ctx, "connectRouteResponse", &apigatewayv2.RouteResponseArgs{
			RouteId:          connectRoute.ID(),
			ApiId:            websocket.ID(),
			RouteResponseKey: pulumi.String("$default"),
		})
		if err != nil {
			return err
		}

		sendMessageFunction, err := createLambda(ctx, "send-message", lambdaRole, pulumi.StringMap{
			"TABLE_NAME": connectionTable.Name,
			"API_GATEWAY_ENDPOINT": pulumi.All(stage.ApiId, region.Name, stage.Name).ApplyT(func(args []interface{}) (string, error) {
				return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s", args[0], args[1], args[2]), nil
			}).(pulumi.StringOutput),
			"REGION": pulumi.String(region.Name),
		})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "sendMessageInvokePermission", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  sendMessageFunction.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: websocket.ExecutionArn.ApplyT(func(executionArn string) (string, error) {
				return fmt.Sprintf("%v/*", executionArn), nil
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}

		sendMessageIntegration, err := apigatewayv2.NewIntegration(ctx, "sendMessageIntegration", &apigatewayv2.IntegrationArgs{
			ApiId:                   websocket.ID(),
			IntegrationType:         pulumi.String("AWS_PROXY"),
			ConnectionType:          pulumi.String("INTERNET"),
			ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
			IntegrationMethod:       pulumi.String("POST"),
			IntegrationUri:          sendMessageFunction.InvokeArn,
			PassthroughBehavior:     pulumi.String("WHEN_NO_MATCH"),
		})
		if err != nil {
			return err
		}

		sendMessageRoute, err := apigatewayv2.NewRoute(ctx, "sendMessageRoute", &apigatewayv2.RouteArgs{
			ApiId:                            websocket.ID(),
			RouteKey:                         pulumi.String("$default"),
			AuthorizationType:                pulumi.String("NONE"),
			RouteResponseSelectionExpression: pulumi.String("$default"),
			Target: sendMessageIntegration.ID().ApplyT(func(id string) (string, error) {
				return fmt.Sprintf("integrations/%v", id), nil
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}

		_, err = apigatewayv2.NewRouteResponse(ctx, "sendMessageRouteResponse", &apigatewayv2.RouteResponseArgs{
			RouteId:          sendMessageRoute.ID(),
			ApiId:            websocket.ID(),
			RouteResponseKey: pulumi.String("$default"),
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
