package router

import (
	"net/http"

	"github.com/Anand-S23/complete-auth/src/controller"
	"github.com/gorilla/handlers"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func Fn(fn apiFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := fn(w, r)
        if err != nil {
            controller.WriteJSON(w, http.StatusInternalServerError, controller.ErrMsg(err.Error()))
        }
    }
}

func NewRouter(c *controller.Controller) http.Handler {
    router := http.NewServeMux()
    router.HandleFunc("GET /ping", Fn(c.Ping)) // Health Check

    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

    return corsHandler(router)
}

