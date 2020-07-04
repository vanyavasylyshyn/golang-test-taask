package services

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"

	"github.com/vanyavasylyshyn/golang-test-task/models"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DestroyRefreshToken ...
func DestroyRefreshToken(refreshToken string) map[string]interface{} {
	client := models.Client
	db := client.Database(os.Getenv("DB_NAME"))
	refreshTokenCollection := db.Collection("refresh-tokens")

	decodedRefreshToken, err := b64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		u.LogError("[ERROR] Decoding from string: ", err)
		return u.Message(false, "Token is not valid.")
	}

	refreshClaims, err := ExtractTokenMetadata(decodedRefreshToken, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		u.LogError("[ERROR] Extracting refresh token metadata: ", err)
		return u.Message(false, "Token is expired.")
	}

	refreshObjectID, err := primitive.ObjectIDFromHex(refreshClaims.Id)
	if err != nil {
		u.LogError("[ERROR] Converting string refreshID to ObjectID: ", err)
		return u.Message(false, "Error while deleting.")
	}

	fmt.Print("\n")
	fmt.Print(refreshClaims.Id)
	fmt.Print("\n")

	result := refreshTokenCollection.FindOneAndDelete(context.Background(), bson.M{
		"_id": refreshObjectID,
	})
	if result.Err() != nil {
		u.LogError("[ERROR] Deleting token from database: ", result.Err())
		return u.Message(false, "Error while deleting.")
	}

	return u.Message(true, "Token was deleted.")
}
