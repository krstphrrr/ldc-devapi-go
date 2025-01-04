package auth

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// Setup Cognito client
func SetupCognito(userPoolID, clientID string) (*cognitoidentityprovider.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return cognitoidentityprovider.NewFromConfig(cfg), nil
}

// VerifyJWTWithCognito verifies the JWT token 
func VerifyJWTWithCognito(token string, userPoolID, clientID string) (map[string]interface{}, error) {

	return nil, errors.New("AWS Cognito verification not yet implemented")
}
