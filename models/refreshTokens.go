package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RefreshToken ...
type RefreshToken struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	User     primitive.ObjectID `bson:"user_id,omitempty"`
	Token    []byte             `bson:"refresh_token,omitempty"`
	IsActive bool               `bson:"is_active,omitempty"`
}

// Generate ...
func (refreshToken *RefreshToken) Generate(userID string, pairID string) error {
	claims := &TokenClaims{
		UserID: userID,
		PairID: pairID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	t, err := token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		u.LogError("[ERROR] Signing refresh token: ", err)
		return err
	}

	refreshToken.Token = []byte(t)
	refreshToken.IsActive = true
	refreshToken.User, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		u.LogError("[ERROR] Converting string userID to ObjectID for refresh token: ", err)
		return err
	}

	return nil
}

// EncryptToken ...
func (refreshToken *RefreshToken) EncryptToken() error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken.Token), bcrypt.DefaultCost)
	if err != nil {
		u.LogError("[ERROR] Encryption refresh token: ", err)
		return err
	}

	refreshToken.Token = hashedToken
	return nil
}
