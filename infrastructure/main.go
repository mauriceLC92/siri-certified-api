package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		account, err := aws.GetCallerIdentity(ctx, &aws.GetCallerIdentityArgs{})
		if err != nil {
			return err
		}

		region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
		if err != nil {
			return err
		}

		// IAM Role for the Lambda function
		lambdaRole, err := iam.NewRole(ctx, "lambdaRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				  "Version": "2012-10-17",
				  "Statement": [
					{
					  "Action": "sts:AssumeRole",
					  "Principal": {
						"Service": "lambda.amazonaws.com"
					  },
					  "Effect": "Allow",
					  "Sid": ""
					}
				  ]
				}`),
		})
		if err != nil {
			return err
		}

		// DynamoDB table
		table, err := dynamodb.NewTable(ctx, "ItemTable", &dynamodb.TableArgs{
			Attributes: dynamodb.TableAttributeArray{
				&dynamodb.TableAttributeArgs{
					Name: pulumi.String("ID"),
					Type: pulumi.String("S"),
				},
			},
			HashKey:       pulumi.String("ID"),
			ReadCapacity:  pulumi.Int(1),
			WriteCapacity: pulumi.Int(1),
		})
		if err != nil {
			return err
		}
		policyJSON := pulumi.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
			    {
					"Action": ["dynamodb:PutItem", "dynamodb:GetItem", "dynamodb:DeleteItem", "dynamodb:UpdateItem", "dynamodb:Query", "dynamodb:Scan"],
					"Effect": "Allow",
					"Resource": "%s"
				}
			]
		}`, table.Arn)

		// AWS Policy with DynamoDB read and write permissions
		policy, err := iam.NewPolicy(ctx, "lambdaDynamoPolicy", &iam.PolicyArgs{
			Description: pulumi.String("DynamoDB read and write access"),
			Policy:      policyJSON,
		})
		if err != nil {
			return err
		}

		// Connect policy with lambda function role
		_, err = iam.NewRolePolicyAttachment(ctx, "RolePolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      lambdaRole.Name,
			PolicyArn: policy.Arn,
		})
		if err != nil {
			return err
		}

		// Create a custom policy to allow writing logs to CloudWatch
		cwPpolicy, err := iam.NewPolicy(ctx, "cloudWatchLogWritePolicy", &iam.PolicyArgs{
			// TODO - update the Resource to limit it by region and account: `"arn:aws:logs:REGION:ACCOUNT_ID:*"`
			Description: pulumi.String("Allow writing logs to CloudWatch"),
			Policy: pulumi.String(`{
						"Version": "2012-10-17",
						"Statement": [{
							"Effect": "Allow",
							"Action": [
								"logs:CreateLogGroup",
								"logs:CreateLogStream",
								"logs:PutLogEvents"
							],
							"Resource": "*"
						}]
					}`),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, "cloudWatchLogWritePolicyRolePolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      lambdaRole.Name,
			PolicyArn: cwPpolicy.Arn,
		})
		if err != nil {
			return err
		}

		companiesLambda, err := createLambda(ctx, CreateLambda{functionName: "companiesFunction", archivePath: "../handlers/companies/companies.zip", role: lambdaRole})
		if err != nil {
			return err
		}

		// Exports
		ctx.Export("role", lambdaRole.Arn)
		ctx.Export("table", table.Arn)

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

		ctx.Export("userPoolArn", userPool.Arn)
		ctx.Export("userPoolClient ID", userPoolClient.ID())

		usersLambda, err := createLambda(ctx, CreateLambda{functionName: "usersFunction", archivePath: "../handlers/users/users.zip", role: lambdaRole})
		if err != nil {
			return err
		}

		authLambda, err := createLambda(ctx, CreateLambda{functionName: "dummyFunction", archivePath: "../handlers/auth/auth.zip", role: lambdaRole})
		if err != nil {
			return err
		}

		// All things API gateway related

		api, err := apigatewayv2.NewApi(ctx, "users-api", &apigatewayv2.ApiArgs{
			ProtocolType: pulumi.String("HTTP"),
		})
		if err != nil {
			return err
		}

		// ApplyT here would transform the Output to a string which is then inside an array
		clientIdInArray := userPoolClient.ID().ToStringOutput().ApplyT(func(clientId string) []string {
			return []string{clientId}
		}).(pulumi.StringArrayOutput)

		// Create Authorizer for the API
		authorizer, err := apigatewayv2.NewAuthorizer(ctx, "users-api-authorizer", &apigatewayv2.AuthorizerArgs{
			ApiId:           api.ID(),
			AuthorizerType:  pulumi.String("JWT"),
			IdentitySources: pulumi.ToStringArray([]string{"$request.header.Authorization"}),
			JwtConfiguration: apigatewayv2.AuthorizerJwtConfigurationArgs{
				// Create JwtConfiguration
				Audiences: clientIdInArray,
				Issuer:    pulumi.Sprintf("https://%s", userPool.Endpoint),
			},
			Name: pulumi.String("jwt-authorizer"),
		})
		if err != nil {
			return err
		}

		usersIntegration, err := apigatewayv2.NewIntegration(ctx, "usersIntegration", &apigatewayv2.IntegrationArgs{
			ApiId:             api.ID(),
			IntegrationType:   pulumi.String("AWS_PROXY"),
			IntegrationMethod: pulumi.String("POST"),
			IntegrationUri:    usersLambda.InvokeArn,
		})
		if err != nil {
			return err
		}

		authIntegration, err := apigatewayv2.NewIntegration(ctx, "authIntegration", &apigatewayv2.IntegrationArgs{
			ApiId:           api.ID(),
			IntegrationType: pulumi.String("AWS_PROXY"),
			IntegrationUri:  authLambda.InvokeArn,
		})
		if err != nil {
			return err
		}

		companiesIntegration, err := apigatewayv2.NewIntegration(ctx, "companiesIntegration", &apigatewayv2.IntegrationArgs{
			ApiId:           api.ID(),
			IntegrationType: pulumi.String("AWS_PROXY"),
			IntegrationUri:  companiesLambda.InvokeArn,
		})
		if err != nil {
			return err
		}

		_, err = apigatewayv2.NewRoute(ctx, "userRoute", &apigatewayv2.RouteArgs{
			ApiId:             api.ID(),
			RouteKey:          pulumi.String("POST /users"),
			Target:            pulumi.Sprintf("integrations/%s", usersIntegration.ID()),
			AuthorizationType: pulumi.String("JWT"),
			AuthorizerId:      authorizer.ID(),
		})
		if err != nil {
			return err
		}

		_, err = apigatewayv2.NewRoute(ctx, "authRoute", &apigatewayv2.RouteArgs{
			ApiId:    api.ID(),
			RouteKey: pulumi.String("POST /auth"),
			Target:   pulumi.Sprintf("integrations/%s", authIntegration.ID()),
		})
		if err != nil {
			return err
		}

		_, err = apigatewayv2.NewRoute(ctx, "companiesRoute", &apigatewayv2.RouteArgs{
			ApiId:    api.ID(),
			RouteKey: pulumi.String("GET /companies"),
			Target:   pulumi.Sprintf("integrations/%s", companiesIntegration.ID()),
		})
		if err != nil {
			return err
		}

		_, err = apigatewayv2.NewStage(ctx, "defaultStage", &apigatewayv2.StageArgs{
			ApiId:      api.ID(),
			AutoDeploy: pulumi.Bool(true),         // Automatically deploy changes to this stage
			Name:       pulumi.String("$default"), // Required parameter,
		})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "apiGatewayUsersInvoke", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  usersLambda.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			// "arn:aws:execute-api:region:account-id:api-id/stage/METHOD_HTTP_VERB/Resource-path"
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/$default/POST/users", region.Name, account.AccountId, api.ID()),
		})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "apiGatewayAuthInvoke", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  authLambda.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/$default/POST/auth", region.Name, account.AccountId, api.ID()),
		})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "apiGatewayCompaniesInvoke", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  companiesLambda.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/$default/GET/companies", region.Name, account.AccountId, api.ID()),
		})
		if err != nil {
			return err
		}

		ctx.Export("api.Arn", api.Arn)
		ctx.Export("api.ApiEndpoint", api.ApiEndpoint)

		return nil
	})
}

// Next steps
// Sign up a customer and get the verification code
// Write the function used to verify that customer
// Create a log in function which logs a customer in and returns the token
// Use that token to access the /users route and see if it lets me in, if yes, I now have a functioning Auth system
