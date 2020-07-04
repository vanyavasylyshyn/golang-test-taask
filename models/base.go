package models

import (
	"context"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	u "github.com/vanyavasylyshyn/golang-test-task/utils"
)

// Client ...
var Client *mongo.Client

// Context ...
var Context context.Context

// TokenClaims ...
type TokenClaims struct {
	UserID string
	PairID string
	jwt.StandardClaims
}

// Connect ...
func Connect() {
	//Load .env file
	e := godotenv.Load()
	if e != nil {
		u.LogError("[ERROR] Load env variables: ", e)
	}

	// Create client
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_ATLAS_URI")))
	if err != nil {
		u.LogError("[ERROR] Create mongo client: ", err)
	}

	// Create connect
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		u.LogError("[ERROR] Create mongo connection: ", err)
	}

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		u.LogError("[ERROR] Check mongo connection: ", err)
	}

	Client = client
	Context = ctx
}
