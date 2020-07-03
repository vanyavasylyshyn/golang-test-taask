package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vanyavasylyshyn/golang-test-task/controllers"
	"github.com/vanyavasylyshyn/golang-test-task/models"
)

func main() {
	models.Connect()

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.RootPath).Methods("GET")
	router.HandleFunc("/api/credentials/{userID}/new", controllers.CreateCredentials).Methods("GET")
	router.HandleFunc("/api/credentials/{userID}/destroy/refresh",
		controllers.DestroyRefreshTokensForUser).Methods("DELETE")
	router.HandleFunc("/api/credentials/refresh", controllers.RefreshCredentials).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println("Server starten on port: " + port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}

}
