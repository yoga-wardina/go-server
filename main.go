package main

import (
	"fmt"
	"net/http"

	"go-server/Routes"
)

func main() {
	http.HandleFunc("/", Routes.RootHandler)
	http.HandleFunc("/users", Routes.UsersHandler)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
