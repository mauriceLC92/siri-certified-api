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
var cognitoClientAppId = os.Getenv("COGNITO_CLIENT_APP_ID")

// init is run once each time the Lambda container is initialised. Allows the client to be reused
// between invocations of the handler from the same container before being deprovisioned.
func init() {
	newClient, err := auth.NewCognitoClient(cognitoClientAppId)
	if err != nil {
		log.Fatal(err.Error())
	}
	cognitoClient = newClient
}

type ConfirmationCodeRequest struct {
	ConfirmationCode string `json:"confirmationCode"`
	Email            string `json:"email"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var confirmCodeRequest ConfirmationCodeRequest
	err := json.Unmarshal([]byte(request.Body), &confirmCodeRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}, nil
	}

	service := cognitoClient
	res, err := auth.ConfirmVerificationCode(ctx, service, auth.ConfirmVerificationCodeParams{
		ConfirmationCode: confirmCodeRequest.ConfirmationCode,
		ClientID:         cognitoClientAppId,
		Email:            confirmCodeRequest.Email,
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
