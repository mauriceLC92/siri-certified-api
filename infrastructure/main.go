package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// IAM Role for the Lambda function
		role, err := iam.NewRole(ctx, "lambdaRole", &iam.RoleArgs{
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
			Role:      role.Name,
			PolicyArn: policy.Arn,
		})
		if err != nil {
			return err
		}

		// Create a custom policy to allow writing logs to CloudWatch
		cwPpolicy, err := iam.NewPolicy(ctx, "cloudWatchLogWritePolicy", &iam.PolicyArgs{
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
			Role:      role.Name,
			PolicyArn: cwPpolicy.Arn,
		})
		if err != nil {
			return err
		}

		// The following works as is if you build the bootstrap binary first and then zip it
		// AWS Lambda function
		itemFunction, err := lambda.NewFunction(ctx, "itemFunction", &lambda.FunctionArgs{
			// Code: pulumi.NewAssetArchive(map[string]interface{}{
			// 	"folder": pulumi.NewFileArchive("./handler"),
			// }),
			Code:          pulumi.NewFileArchive("../handler/myFunction.zip"),
			Handler:       pulumi.String("bootstrap"),
			Runtime:       pulumi.String("provided.al2"),
			Role:          role.Arn,
			Architectures: pulumi.ToStringArray([]string{"arm64"}),
			Timeout:       pulumi.Int(300),
			MemorySize:    pulumi.Int(128),
		})
		if err != nil {
			return err
		}

		// CloudWatch event rule for monthly execution
		rule, err := cloudwatch.NewEventRule(ctx, "monthlyRule", &cloudwatch.EventRuleArgs{
			ScheduleExpression: pulumi.String("cron(0 12 1 * ? *)"),
		})
		if err != nil {
			return err
		}

		// CloudWatch event target for the lambda function
		_, err = cloudwatch.NewEventTarget(ctx, "monthlyTarget", &cloudwatch.EventTargetArgs{
			Arn:  itemFunction.Arn,
			Rule: rule.Name,
		})
		if err != nil {
			return err
		}

		// Exports
		ctx.Export("role", role.Arn)
		ctx.Export("table", table.Arn)

		return nil
	})
}
