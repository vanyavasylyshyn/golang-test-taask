package controllers

import (
	"encoding/json"
	"net/http"

	u "github.com/vanyavasylyshyn/golang-test-task/utils"

	"github.com/vanyavasylyshyn/golang-test-task/models"
)

// CreateUser ...
var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	db := models.Client.Database("golang-test-task-db")
	usersCollection := db.Collection("users")

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create(models.Context, usersCollection)
	u.Respond(w, resp)
}
