package main

import (
	"fmt"
	"log"
	"net/http"

	"go-server/Routes"
	"go-server/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	Routes.RegisterRoutes(r)

	corsRouter := middleware.CORSMiddleware(r)
	port := 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), corsRouter))
}
