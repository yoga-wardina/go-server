package main

import (
	"fmt"
	"log"
	"net/http"

	"go-server/Routes"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize a new mux router
	r := mux.NewRouter()

	// Register routes
	Routes.RegisterRoutes(r)

	// Start the server
	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
