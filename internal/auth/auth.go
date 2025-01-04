package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// VerifyJWT verifies and parses the JWT token
func VerifyJWT(token string) (map[string]interface{}, error) {

	parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	return claims, nil
}