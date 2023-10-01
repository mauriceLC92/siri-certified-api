package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createCognitoPool(ctx *pulumi.Context) (*cognito.UserPool, error) {
	userPool, err := cognito.NewUserPool(ctx, "siri-certified-pulumi", &cognito.UserPoolArgs{
		AutoVerifiedAttributes: pulumi.StringArray{
			pulumi.String("email"),
		},
		AdminCreateUserConfig: cognito.UserPoolAdminCreateUserConfigArgs{
			AllowAdminCreateUserOnly: pulumi.Bool(false),
			// Set InviteMessageTemplate to default to enable Cognito to automatically send messages
			InviteMessageTemplate: &cognito.UserPoolAdminCreateUserConfigInviteMessageTemplateArgs{}},
		UsernameAttributes: pulumi.ToStringArray([]string{"email"}),
		EmailConfiguration: cognito.UserPoolEmailConfigurationArgs{
			EmailSendingAccount: pulumi.String("COGNITO_DEFAULT"),
		},
	})
	if err != nil {
		return &cognito.UserPool{}, err
	}

	return userPool, nil
}

func createCognitoPoolClient(ctx *pulumi.Context, userPool *cognito.UserPool) (*cognito.UserPoolClient, error) {
	userPoolClient, err := cognito.NewUserPoolClient(ctx, "client-maurice", &cognito.UserPoolClientArgs{
		UserPoolId:                 userPool.ID(),
		EnableTokenRevocation:      pulumi.Bool(true),
		ExplicitAuthFlows:          pulumi.ToStringArray([]string{"ALLOW_USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_ADMIN_USER_PASSWORD_AUTH"}),
		PreventUserExistenceErrors: pulumi.String("ENABLED"),
	})
	if err != nil {
		return &cognito.UserPoolClient{}, err
	}

	return userPoolClient, nil
}
