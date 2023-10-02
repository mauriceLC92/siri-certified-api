package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Print("hello...")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, from DUMMY, no auth",
	}, nil
}

func main() {
	lambda.Start(handler)
}
