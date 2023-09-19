package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler is a simple function that takes a string and does a ToUpper.
func handler(request events.CloudWatchEvent) error {
	log.Println("Hello Maurice!")

	log.Printf("request.AccountID: %v\n", request.AccountID)
	log.Printf("request.ID: %v\n", request.ID)
	log.Printf("detail-type: %s", request.DetailType)
	return nil
}

func main() {
	lambda.Start(handler)
}
