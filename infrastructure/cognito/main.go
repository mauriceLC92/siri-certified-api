package main

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Cognito User Pool
		userPool, err := createCognitoPool(ctx)
		if err != nil {
			return err
		}

		// Cognito User Pool Client
		userPoolClient, err := createCognitoPoolClient(ctx, userPool)
		if err != nil {
			return err
		}

		ctx.Export("userPool Arn", userPool.Arn)
		ctx.Export("userPool Endpoint", userPool.Endpoint)
		ctx.Export("userPoolClient ID", userPoolClient.ID())

		return nil
	})
}
