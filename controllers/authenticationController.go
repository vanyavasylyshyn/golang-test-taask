package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanyavasylyshyn/golang-test-task/services"
	u "github.com/vanyavasylyshyn/golang-test-task/utils"
)

// CreateCredentials ...
var CreateCredentials = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	resp := services.GenerateCredentials(userID)
	u.Respond(w, resp)
}

// DestroyRefreshTokensForUser ...
var DestroyRefreshTokensForUser = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	resp := services.DestroyRefreshCredentials(userID)
	u.Respond(w, resp)
}
