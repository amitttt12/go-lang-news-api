package routes

import (
	"news-app/handler"

	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router) {
	// Define the signup route and link it to the Signup handler
	router.HandleFunc("/api/signup", handler.Signup).Methods("POST")
	router.HandleFunc("/api/signin", handler.Signin).Methods("POST")
}
