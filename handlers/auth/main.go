package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"siri-certified-api/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var cognitoClient *auth.CognitoClient

// init is run once each time the Lambda container is initialised. Allows the client to be reused
// between invocations of the handler from the same container before being deprovisioned.
func init() {
	newClient, err := auth.NewCognitoClient(os.Getenv("COGNITO_CLIENT_APP_ID"))
	if err != nil {
		log.Fatal(err.Error())
	}
	cognitoClient = newClient
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var srRequest SignUpRequest
	err := json.Unmarshal([]byte(request.Body), &srRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}, nil
	}

	service := cognitoClient
	res, err := auth.SignUpUser(context.TODO(), service, auth.SignUpUserParams{
		Email:    srRequest.Username,
		Password: srRequest.Password,
		ClientID: os.Getenv("COGNITO_CLIENT_APP_ID"),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(resJson),
	}, nil
}

func main() {
	lambda.Start(handler)
}
