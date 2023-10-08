package auth

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type SignUpUserParams struct {
	Email, Password, ClientID string
}

// SignUpUser is the function which signs up a new user.
func SignUpUser(ctx context.Context, service CognitoService, params SignUpUserParams) (*cognitoidentityprovider.SignUpOutput, error) {
	if service == nil {
		return nil, errors.New("service cannot be nil")
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: &params.ClientID,
		Username: &params.Email,
		Password: &params.Password,
	}

	return service.SignUp(ctx, input)
}

type ConfirmVerificationCodeParams struct {
	Email, ConfirmationCode, ClientID string
}

// ConfirmConfirmationCode is the function which confirms the verification code sent to a user after a successful sign up.
func ConfirmVerificationCode(ctx context.Context, service CognitoService, params ConfirmVerificationCodeParams) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	if service == nil {
		return nil, errors.New("service cannot be nil")
	}

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         &params.ClientID,
		ConfirmationCode: &params.ConfirmationCode,
		Username:         &params.Email,
	}

	return service.ConfirmSignUp(ctx, input)
}

type LogInParams struct {
	Email, Password, ClientID string
}

// LogInUser is the function which logs in the user based on their email and password.
func LogInUser(ctx context.Context, service CognitoService, params LogInParams) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	if service == nil {
		return nil, errors.New("service cannot be nil")
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: &params.ClientID,
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": params.Email,
			"PASSWORD": params.Password,
		},
	}

	return service.InitiateAuth(ctx, input)
}

type CognitoService interface {
	SignUp(ctx context.Context, params *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error)
	ConfirmSignUp(ctx context.Context, params *cognitoidentityprovider.ConfirmSignUpInput) (*cognitoidentityprovider.ConfirmSignUpOutput, error)
	InitiateAuth(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
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

// ConfirmSignUp calls the ConfirmSignUp method on the actual Cognito client.
func (c *CognitoClient) ConfirmSignUp(ctx context.Context, params *cognitoidentityprovider.ConfirmSignUpInput) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	return c.Client.ConfirmSignUp(ctx, params)
}

// InitiateAuth calls the InitiateAuth method on the actual Cognito client.
func (c *CognitoClient) InitiateAuth(ctx context.Context, params *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return c.Client.InitiateAuth(ctx, params)
}
