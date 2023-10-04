package auth

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// CognitoService is an interface that defines the SignUp method.
// This allows us to mock the Cognito client for testing.
type CognitoService interface {
	SignUp(ctx context.Context, params *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error)
}

// CognitoClient is a struct that wraps the actual Cognito client.
type CognitoClient struct {
	Client   *cognitoidentityprovider.Client
	ClientId string
}

// NewCognitoClient initializes a new CognitoClient with the given ClientId.
func NewCognitoClient(clientId string) (*CognitoClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	cognitoProviderClient := cognitoidentityprovider.NewFromConfig(cfg)
	return &CognitoClient{Client: cognitoProviderClient, ClientId: clientId}, nil
}

// SignUp calls the SignUp method on the actual Cognito client.
func (c *CognitoClient) SignUp(ctx context.Context, params *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error) {
	return c.Client.SignUp(ctx, params)
}

// SignUpUser is the function to sign up a user.
func SignUpUser(ctx context.Context, service CognitoService, email, password, clientID string) (*cognitoidentityprovider.SignUpOutput, error) {
	if service == nil {
		return nil, errors.New("service cannot be nil")
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: &clientID,
		Username: &email,
		Password: &password,
	}

	return service.SignUp(ctx, input)
}
