package main

import (
	"log"
	"os"
	"path"
	"siri-certified-api/companies"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() ([]companies.Company, error) {
	log.Println("Hello from companies lambda")
	lambdaDir := os.Getenv("LAMBDA_TASK_ROOT")
	jsonFilePath := path.Join(lambdaDir, "company-data.json")

	c, err := companies.GetCompanies(jsonFilePath)
	if err != nil {
		return []companies.Company{}, err
	}

	return c, nil
}

func main() {
	lambda.Start(handler)
}
