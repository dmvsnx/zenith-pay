package mocks

import "github.com/savanyv/zenith-pay/internal/utils/helpers"

type JWTService struct {
	GenerateAccessTokenFn func(userID, username, role string, tokenVersion int) (string, error)
	ValidateAccessTokenFn func(tokenString string) (*helpers.JWTClaims, error)
}

func (m *JWTService) GenerateAccessToken(userID, username, role string, tokenVersion int) (string, error) {
	return m.GenerateAccessTokenFn(userID, username, role, tokenVersion)
}

func (m *JWTService) ValidateAccessToken(tokenString string) (*helpers.JWTClaims, error) {
	return m.ValidateAccessTokenFn(tokenString)
}
