package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path"
	"siri-certified-api/companies"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Hello from companies lambda")
	lambdaDir := os.Getenv("LAMBDA_TASK_ROOT")
	jsonFilePath := path.Join(lambdaDir, "company-data.json")

	c, err := companies.GetCompanies(jsonFilePath)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, nil
}

func main() {
	lambda.Start(handler)
}
