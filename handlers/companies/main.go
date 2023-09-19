package main

import (
	"log"
	"siri-certified-api/companies"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler() (*companies.Company, error) {
	log.Println("Hello from companies lambda")

	company := &companies.Company{
		Title:     "Hello World Corp",
		Link:      "https://www.helloworld.com",
		CVRNumber: "1814569",
	}

	return company, nil
}

func main() {
	lambda.Start(handler)
}
