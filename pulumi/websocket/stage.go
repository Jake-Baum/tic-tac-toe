package websocket

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"time"
)

type ApiStage struct {
	pulumi.ResourceState

	Name         string
	WebSocketUrl pulumi.StringOutput
}

type ApiStageArgs struct {
	Api           *apigatewayv2.Api
	LambdaProxies []*LambdaProxy
}

func NewApiStage(ctx *pulumi.Context, name string, args ApiStageArgs, opts ...pulumi.ResourceOption) (*ApiStage, error) {
	apiStage := &ApiStage{
		Name: name,
	}
	err := ctx.RegisterComponentResource("jakebaum:websocket:ApiStage", name, apiStage, opts...)
	if err != nil {
		return nil, err
	}

	parentResourceOption := pulumi.ResourceOption(pulumi.Parent(apiStage))

	dependsOn := make([]pulumi.Resource, len(args.LambdaProxies))
	for index, lambdaProxy := range args.LambdaProxies {
		dependsOn[index] = lambdaProxy
	}

	websocketDeployment, err := apigatewayv2.NewDeployment(ctx, name, &apigatewayv2.DeploymentArgs{
		ApiId: args.Api.ID(),
		Triggers: pulumi.StringMap{
			"deployedAt": pulumi.String(time.Now().Format(time.RFC3339)), // This is somewhat of a hack to force the API to redeploy on changes.  Must be a better way
		},
	}, pulumi.DependsOn(dependsOn), parentResourceOption)
	if err != nil {
		return nil, err
	}

	stage, err := apigatewayv2.NewStage(ctx, name, &apigatewayv2.StageArgs{
		ApiId:        args.Api.ID(),
		DeploymentId: websocketDeployment.ID(),
		DefaultRouteSettings: apigatewayv2.StageDefaultRouteSettingsArgs{
			DataTraceEnabled:       pulumi.Bool(true),
			DetailedMetricsEnabled: pulumi.Bool(true),
			LoggingLevel:           pulumi.String("INFO"),
			ThrottlingBurstLimit:   pulumi.Int(100),
			ThrottlingRateLimit:    pulumi.Float64(100),
		},
	}, parentResourceOption)
	if err != nil {
		return nil, err
	}

	apiStage.WebSocketUrl = stage.InvokeUrl

	return apiStage, nil
}
