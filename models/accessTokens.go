package models

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// AccessToken ...
type AccessToken struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	User     primitive.ObjectID `bson:"user_id,omitempty"`
	Token    []byte             `bson:"access_token,omitempty"`
	IsActive bool               `bson:"is_active,omitempty"`
}

// Generate ...
func (accessToken *AccessToken) Generate(userID string, pairID string) error {
	tokenID := primitive.NewObjectID()

	claims := &TokenClaims{
		UserID: userID,
		PairID: pairID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Id:        tokenID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	t, err := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		u.LogError("[ERROR] Signing access token: ", err)
		return err
	}

	accessToken.ID = tokenID
	accessToken.Token = []byte(t)
	accessToken.IsActive = true
	accessToken.User, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Print(userID)
		u.LogError("[ERROR] Converting string userID to ObjectID for access token: ", err)
		return err
	}

	return nil
}

// EncryptToken ...
func (accessToken *AccessToken) EncryptToken() error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(accessToken.Token), bcrypt.DefaultCost)
	if err != nil {
		u.LogError("[ERROR] Encryption refresh token: ", err)
		return err
	}

	accessToken.Token = hashedToken
	return nil
}
