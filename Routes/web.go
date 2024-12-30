package Routes

import (
	"net/http"

	"go-server/Middleware"

	"go-server/Handlers"

	"github.com/gorilla/mux"
)


func RegisterRoutes(r *mux.Router) {
	// Route with middleware
	r.Handle("/protected", Middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected route"))
	})))

	// Other routes
	r.HandleFunc("/", Handlers.RootHandler)
	r.HandleFunc("/users", Handlers.UsersHandler)
}