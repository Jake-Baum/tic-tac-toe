package websocket

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LambdaProxy struct {
	pulumi.ResourceState

	Name           string
	LambdaFunction *lambda.Function
	Route          *apigatewayv2.Route
}

type LambdaProxyArgs struct {
	LambdaRole        *iam.Role
	Api               *apigatewayv2.Api
	LambdaEnvironment pulumi.StringMap
	RouteKey          string
}

func NewLambdaProxy(ctx *pulumi.Context, name string, args LambdaProxyArgs, opts ...pulumi.ResourceOption) (*LambdaProxy, error) {
	lambdaProxyGroup := &LambdaProxy{
		Name: name,
	}
	err := ctx.RegisterComponentResource("jakebaum:websocket:LambdaProxy", name, lambdaProxyGroup, opts...)
	if err != nil {
		return nil, err
	}

	parentResourceOption := pulumi.ResourceOption(pulumi.Parent(lambdaProxyGroup))

	lambdaProxyGroup.LambdaFunction, err = lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Runtime: pulumi.String("go1.x"),
		Handler: pulumi.String("main"),
		Role:    args.LambdaRole.Arn,
		Code:    pulumi.NewFileArchive(fmt.Sprintf("../bin/lambda/%s/main.zip", name)),
		Environment: lambda.FunctionEnvironmentArgs{
			Variables: args.LambdaEnvironment,
		},
	}, parentResourceOption)
	if err != nil {
		return nil, err
	}

	_, err = lambda.NewPermission(ctx, name, &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  lambdaProxyGroup.LambdaFunction.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: args.Api.ExecutionArn.ApplyT(func(executionArn string) (string, error) {
			return fmt.Sprintf("%v/*", executionArn), nil
		}).(pulumi.StringOutput),
	}, parentResourceOption)
	if err != nil {
		return nil, err
	}

	integration, err := apigatewayv2.NewIntegration(ctx, name, &apigatewayv2.IntegrationArgs{
		ApiId:                   args.Api.ID(),
		IntegrationType:         pulumi.String("AWS_PROXY"),
		ConnectionType:          pulumi.String("INTERNET"),
		ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
		IntegrationMethod:       pulumi.String("POST"),
		IntegrationUri:          lambdaProxyGroup.LambdaFunction.InvokeArn,
		PassthroughBehavior:     pulumi.String("WHEN_NO_MATCH"),
	}, parentResourceOption)
	if err != nil {
		return nil, err
	}

	lambdaProxyGroup.Route, err = apigatewayv2.NewRoute(ctx, name, &apigatewayv2.RouteArgs{
		ApiId:                            args.Api.ID(),
		RouteKey:                         pulumi.String(args.RouteKey),
		AuthorizationType:                pulumi.String("NONE"),
		RouteResponseSelectionExpression: pulumi.String("$default"),
		Target: integration.ID().ApplyT(func(id string) (string, error) {
			return fmt.Sprintf("integrations/%v", id), nil
		}).(pulumi.StringOutput),
	}, parentResourceOption)
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRouteResponse(ctx, name, &apigatewayv2.RouteResponseArgs{
		RouteId:          lambdaProxyGroup.Route.ID(),
		ApiId:            args.Api.ID(),
		RouteResponseKey: pulumi.String("$default"),
	}, parentResourceOption)
	if err != nil {
		return nil, err
	}

	return lambdaProxyGroup, err
}
