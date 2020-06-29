package controllers

import (
	"net/http"

	u "github.com/vanyavasylyshyn/golang-test-task/utils"
)

// RootPath ...
var RootPath = func(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{"message": "Hi! See API description here: https://github.com/vanyavasylyshyn/golang-test-task"}

	u.Respond(w, resp)
}
