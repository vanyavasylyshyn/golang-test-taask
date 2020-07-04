package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanyavasylyshyn/golang-test-task/helpers"
	"github.com/vanyavasylyshyn/golang-test-task/services"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
)

// CreateCredentials ...
var CreateCredentials = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	resp := services.GenerateTokens(userID)

	u.Respond(w, resp)
}

// DestroyAllRefreshCredentials ...
var DestroyAllRefreshCredentials = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	resp := services.DestroyAllRefreshTokens(userID)
	u.Respond(w, resp)
}

// RefreshCredentials ...
var RefreshCredentials = func(w http.ResponseWriter, r *http.Request) {
	token := helpers.ExtractToken(r)

	resp := services.RefreshTokens(token)

	u.Respond(w, resp)
}

// DestroyRefreshCredential ...
var DestroyRefreshCredential = func(w http.ResponseWriter, r *http.Request) {
	token := helpers.ExtractToken(r)

	resp := services.DestroyRefreshToken(token)

	u.Respond(w, resp)
}
