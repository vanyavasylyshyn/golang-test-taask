package models

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccessToken ...
type AccessToken struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Token string             `bson:"access_token,omitempty"`
	User  primitive.ObjectID `bson:"user_id,omitempty"`
}

// Generate ...
func (accessToken *AccessToken) Generate(userID string, pairID string) error {

	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["pairID"] = pairID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		u.LogError("[ERROR] Signing access token: ", err)
		return err
	}

	accessToken.Token = t
	accessToken.User, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Print(userID)
		u.LogError("[ERROR] Converting string userID to ObjectID for access token: ", err)
		return err
	}

	return nil
}
