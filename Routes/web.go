package Routes

import (
	"net/http"

	"go-server/middleware"

	"go-server/handler"

	"github.com/gorilla/mux"
)


func RegisterRoutes(r *mux.Router) {
	// Route with middleware
	r.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected route"))
	}))).Methods("GET")

	// Other routes
	r.HandleFunc("/", handler.RootHandler).Methods("GET")
	r.HandleFunc("/users", handler.UsersHandler).Methods("GET")
	r.HandleFunc("/auth/signup", handler.CreateUserHandler).Methods("POST")
	r.HandleFunc("/auth/login", handler.LoginHandler).Methods("POST")
	
}