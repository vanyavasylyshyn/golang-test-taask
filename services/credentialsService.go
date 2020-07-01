package services

import (
	"context"
	"fmt"
	"os"

	b64 "encoding/base64"

	"github.com/vanyavasylyshyn/golang-test-task/helpers"
	"github.com/vanyavasylyshyn/golang-test-task/models"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GenerateCredentials ...
func GenerateCredentials(userID string) map[string]interface{} {
	client := models.Client
	db := client.Database(os.Getenv("DB_NAME"))
	accessTokenCollection := db.Collection("access-tokens")
	refreshTokenCollection := db.Collection("refresh-tokens")
	//If we could have user database,  check if  user exists

	pairID := helpers.GenerateRandomUUID()

	accessToken := models.AccessToken{}
	err := accessToken.Generate(userID, pairID)
	if err != nil {
		return u.Message(false, "Internal server eroor.")
	}
	refreshToken := models.RefreshToken{}
	err = refreshToken.Generate(userID, pairID)
	if err != nil {
		return u.Message(false, "Internal server eroor.")
	}

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessionContext mongo.SessionContext) (interface{}, error) {
		result, err := accessTokenCollection.InsertOne(
			sessionContext,
			accessToken,
		)
		if err != nil {
			return nil, err
		}

		result, err = refreshTokenCollection.InsertOne(
			sessionContext,
			refreshToken,
		)
		if err != nil {
			return nil, err
		}

		return result, err
	})
	if err != nil {
		u.LogError("[ERROR] Saving credentials: ", err)
		return u.Message(true, "Credentials has not been created.")
	}

	result := u.Message(true, "Credentials has been created.")
	result["accessToken"] = accessToken.Token
	result["refreshToken"] = b64.StdEncoding.EncodeToString(refreshToken.Token)
	return result
}

// DestroyRefreshCredentials ...
func DestroyRefreshCredentials(userID string) map[string]interface{} {
	client := models.Client
	db := client.Database(os.Getenv("DB_NAME"))
	refreshTokenCollection := db.Collection("refresh-tokens")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		u.LogError("[ERROR] Converting string userID to ObjectID for refresh token: ", err)
		return u.Message(false, "Error while deleting.")
	}

	result, err := refreshTokenCollection.DeleteMany(context.Background(), bson.M{
		"user_id": userObjectID,
	})
	if err != nil {
		return u.Message(false, "Error while deleting.")
	}

	return u.Message(true, fmt.Sprintf("Deleted %d refresh tokens.", result.DeletedCount))
}
