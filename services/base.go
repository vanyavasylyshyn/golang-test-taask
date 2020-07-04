package services

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/vanyavasylyshyn/golang-test-task/helpers"
	"github.com/vanyavasylyshyn/golang-test-task/models"
)

// TokenDetails ...
type TokenDetails struct {
	AccessToken              []byte
	RefreshToken             []byte
	HashedAccessTokenObject  models.AccessToken
	HashedRefreshTokenObject models.RefreshToken
}

// CreateTokens ...
func CreateTokens(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}

	pairID := helpers.GenerateRandomUUID()

	accessToken := models.AccessToken{}
	err := accessToken.Generate(userID, pairID)
	if err != nil {
		return nil, err
	}
	refreshToken := models.RefreshToken{}
	err = refreshToken.Generate(userID, pairID)
	if err != nil {
		return nil, err
	}

	td.AccessToken = accessToken.Token
	td.RefreshToken = refreshToken.Token

	err = refreshToken.EncryptToken()
	if err != nil {
		return nil, err
	}

	td.HashedAccessTokenObject = accessToken
	td.HashedRefreshTokenObject = refreshToken

	return td, nil
}

// ExtractTokenMetadata ...
func ExtractTokenMetadata(tokenString []byte, secret string) (*models.TokenClaims, error) {
	claims := &models.TokenClaims{}

	tkn, err := jwt.ParseWithClaims(string(tokenString), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("Unexpected signing method: %v", err)
		}

		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}

	return claims, nil
}
