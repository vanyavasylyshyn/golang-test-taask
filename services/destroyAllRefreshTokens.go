package services

import (
	"context"
	"fmt"
	"os"

	"github.com/vanyavasylyshyn/golang-test-task/models"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DestroyAllRefreshTokens ...
func DestroyAllRefreshTokens(userID string) map[string]interface{} {
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
		u.LogError("[ERROR] Deleting refresh tokens: ", err)
		return u.Message(false, "Error while deleting.")
	}

	return u.Message(true, fmt.Sprintf("Deleted %d refresh tokens.", result.DeletedCount))
}
