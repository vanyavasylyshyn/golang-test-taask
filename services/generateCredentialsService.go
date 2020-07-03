package services

import (
	"context"
	b64 "encoding/base64"
	"os"

	"github.com/vanyavasylyshyn/golang-test-task/helpers"
	"github.com/vanyavasylyshyn/golang-test-task/models"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/mongo"
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

// GenerateCredentials ...
func GenerateCredentials(userID string) map[string]interface{} {
	client := models.Client
	db := client.Database(os.Getenv("DB_NAME"))
	accessTokenCollection := db.Collection("access-tokens")
	refreshTokenCollection := db.Collection("refresh-tokens")
	//If we could have user database,  check if  user exists

	tokenDetails, err := CreateTokens(userID)
	if err != nil {
		u.LogError("[ERROR] Creating tokens: ", err)
		return u.Message(false, "Internal server error.")
	}

	session, err := client.StartSession()
	if err != nil {
		u.LogError("[ERROR] Startion transaction session: ", err)
		return u.Message(false, "Internal server error.")
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessionContext mongo.SessionContext) (interface{}, error) {
		result, err := accessTokenCollection.InsertOne(
			sessionContext,
			tokenDetails.HashedAccessTokenObject,
		)
		if err != nil {
			return nil, err
		}

		result, err = refreshTokenCollection.InsertOne(
			sessionContext,
			tokenDetails.HashedRefreshTokenObject,
		)
		if err != nil {
			return nil, err
		}

		return result, err
	})
	if err != nil {
		u.LogError("[ERROR] Saving credentials: ", err)
		return u.Message(false, "Credentials has not been created.")
	}

	b64RefreshToken := b64.StdEncoding.EncodeToString(tokenDetails.RefreshToken)

	result := u.Message(true, "Credentials has been created.")
	result["accessToken"] = string(tokenDetails.AccessToken)
	result["refreshToken"] = b64RefreshToken
	return result
}
