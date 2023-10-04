package main

import (
	"context"
	"encoding/json"

	"siri-certified-api/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	var srRequest SignUpRequest
	err = json.Unmarshal([]byte(request.Body), &srRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid input"}, nil
	}

	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)
	service := &auth.CognitoClient{Client: cognitoClient}

	res, err := auth.SignUpUser(context.TODO(), service, srRequest.Username, srRequest.Password, "")
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
