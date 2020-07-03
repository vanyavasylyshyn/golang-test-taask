package services

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/vanyavasylyshyn/golang-test-task/models"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

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

// RefreshCredentials ...
func RefreshCredentials(refreshToken string) map[string]interface{} {
	client := models.Client
	db := client.Database(os.Getenv("DB_NAME"))
	refreshTokenCollection := db.Collection("refresh-tokens")
	accessTokenCollection := db.Collection("access-tokens")

	decodedRefreshToken, err := b64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		u.LogError("[ERROR] Decoding from string: ", err)
		return u.Message(false, "Internal server error.")
	}

	refreshClaims, err := ExtractTokenMetadata(decodedRefreshToken, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		u.LogError("[ERROR] Extracting refresh token metadata: ", err)
		return u.Message(false, "Token is expired.")
	}

	userObjectID, err := primitive.ObjectIDFromHex(refreshClaims.UserID)
	if err != nil {
		u.LogError("[ERROR] Converting string userID to ObjectID for refresh token: ", err)
		return u.Message(false, "Error while refreshing.")
	}

	accessTokenHash := &models.AccessToken{}
	refreshTokenHash := &models.RefreshToken{}

	// винести в окрему функцію в base.go
	err = accessTokenCollection.FindOne(context.Background(), bson.M{
		"user_id":   userObjectID,
		"is_active": true,
	}).Decode(accessTokenHash)
	if err != nil {
		u.LogError("[ERROR] Decoding token from database: ", err)
		return u.Message(false, "Error while refreshing.")
	}

	err = refreshTokenCollection.FindOne(context.Background(), bson.M{
		"user_id":   userObjectID,
		"is_active": true,
	}).Decode(refreshTokenHash)
	if err != nil {
		u.LogError("[ERROR] Decoding token from database: ", err)
		return u.Message(false, "Error while refreshing.")
	}

	// Винести в модель
	err = bcrypt.CompareHashAndPassword(refreshTokenHash.Token, []byte(decodedRefreshToken))
	if err != nil {
		u.LogError("[ERROR] Comparing hash: ", err)
		return u.Message(false, "You have no rights to refresh.")
	}
	// ...

	accessClaims, err := ExtractTokenMetadata(accessTokenHash.Token, os.Getenv("ACCESS_SECRET"))
	if err != nil {
		u.LogError("[ERROR] Extracting access token metadata: ", err)
		return u.Message(false, "")
	}

	if accessClaims.PairID != refreshClaims.PairID {
		return u.Message(false, "You have no rights to refresh.")
	}

	tokenDetails, err := CreateTokens(refreshClaims.UserID)
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
		// винести у base.go
		accessTokenCollection.FindOneAndUpdate(
			sessionContext,
			bson.M{
				"user_id":   userObjectID,
				"is_active": true,
			},
			bson.D{
				{"$set", bson.D{{"is_active", false}}},
			},
		)

		refreshTokenCollection.FindOneAndUpdate(
			sessionContext,
			bson.M{
				"user_id":   userObjectID,
				"is_active": true,
			},
			bson.D{
				{"$set", bson.D{{"is_active", false}}},
			},
		)
		// ...

		// винести у base.go
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
		// ...

		return result, err
	})
	if err != nil {
		u.LogError("[ERROR] Refreshing credentials: ", err)
		return u.Message(false, "Credentials has not been created.")
	}
	b64RefreshToken := b64.StdEncoding.EncodeToString(tokenDetails.RefreshToken)

	result := u.Message(true, "Credentials has been created.")
	result["accessToken"] = tokenDetails.AccessToken
	result["refreshToken"] = b64RefreshToken
	return result
}
