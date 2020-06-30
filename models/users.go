package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	u "github.com/vanyavasylyshyn/golang-test-task/utils"
)

// User ...
type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          string               `bson:"name,omitempty"`
	Surname       string               `bson:"surname,omitempty"`
	AccessToken   primitive.ObjectID   `bson:"access_token,omitempty"`
	RefreshTokens []primitive.ObjectID `bson:"refresh_tokens,omitempty"`
}

// Validate ...
func (user *User) Validate() (map[string]interface{}, bool) {
	if len(user.Name) < 2 {
		return u.Message(false, "Name is too short."), false
	}

	if len(user.Surname) < 2 {
		return u.Message(false, "Surname is too short."), false
	}

	return u.Message(true, "Requirement passed."), true
}

// Create ...
func (user *User) Create(ctx context.Context, usersCollection *mongo.Collection) map[string]interface{} {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	insertResult, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		u.LogError("[ERROR] Insertion error: ", err)
		return u.Message(true, "Internal server error.")
	}

	result := u.Message(true, "Account has been created")
	result["user"] = insertResult
	return result
}
