package main

import (
	"log"
	"net/http"
	"news-app/config"
	"news-app/routes"

	"github.com/gorilla/mux"
)

// func main() {
// 	// Initialize the database connection
// 	db, err := config.InitDatabase()
// 	if err != nil {
// 		log.Fatalf("Could not connect to the database: %v", err)
// 	}
// 	defer db.Close()

// 	// Create a new router
// 	router := mux.NewRouter()

// 	// Register routes
// 	routes.AuthRoutes(router)

// 	// Start the server
// 	log.Println("Server is running on http://localhost:8080")
// 	err = http.ListenAndServe(":8080", router)
// 	if err != nil {
// 		log.Fatalf("Could not start the server: %v", err)
// 	}
// }


func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not set
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, Railway!"))
    })

    http.ListenAndServe(":"+port, nil)
}
