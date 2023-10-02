package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CreateLambda struct {
	functionName string
	archivePath  string
	role         *iam.Role
}

func createLambda(ctx *pulumi.Context, cr CreateLambda) (*lambda.Function, error) {
	// The following works as is if you build the bootstrap binary first and then zip it
	// AWS Lambda function
	lambdaFunc, err := lambda.NewFunction(ctx, cr.functionName, &lambda.FunctionArgs{
		// Code: pulumi.NewAssetArchive(map[string]interface{}{
		// 	"folder": pulumi.NewFileArchive("./handler"),
		// }),
		Code:          pulumi.NewFileArchive(cr.archivePath),
		Handler:       pulumi.String("bootstrap"),
		Runtime:       lambda.RuntimeCustomAL2,
		Role:          cr.role.Arn,
		Architectures: pulumi.ToStringArray([]string{"arm64"}),
		Timeout:       pulumi.Int(300),
		MemorySize:    pulumi.Int(128),
	})
	if err != nil {
		return &lambda.Function{}, err
	}

	return lambdaFunc, nil
}
